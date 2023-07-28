package controllers

import (
	"boilerplate-api/api/services"
	"boilerplate-api/constants"
	"boilerplate-api/dtos"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/paginations"
	"boilerplate-api/responses"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BackupController struct
type BackupController struct {
	logger      infrastructure.Logger
	TodoService services.TodoService
	UserService services.UserService
}

// NewBackupController constructor
func NewBackupController(
	logger infrastructure.Logger,
	TodoService services.TodoService,
	UserService services.UserService,
) BackupController {
	return BackupController{
		logger:      logger,
		TodoService: TodoService,
		UserService: UserService,
	}
}

func (cc BackupController) BackupData(ctx *gin.Context) {
	// Sample data to export
	data := make(map[string]interface{})
	cur := time.Now().Unix()
	filename := fmt.Sprintf("media/backup_%v.json", cur)

	getTodo, _, _ := cc.TodoService.GetAllTodo()
	pagination := paginations.BuildPagination[*paginations.UserPagination](ctx)
	users, _, _ := cc.UserService.GetAllUsers(*pagination)
	data["todo"] = getTodo
	data["users"] = users

	file, _ := json.MarshalIndent(data, "", " ")
	if err := ioutil.WriteFile(filename, file, 0644); err != nil {
		responses.JSON(ctx, http.StatusInternalServerError, err)
		return
	}
	message := fmt.Sprintf("Data exported to : %v", filename)
	responses.JSON(ctx, http.StatusOK, message)
}

func (cc BackupController) RestoreData(ctx *gin.Context) {
	dto := dtos.FilePath{}
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		cc.logger.Zap.Error("Error [CreateTodo] (ShouldBindJson) : ", err)
		err := errors.BadRequest.Wrap(err, "Failed to bind Todo")
		responses.HandleError(ctx, err)
		return
	}

	trx := ctx.MustGet(constants.DBTransaction).(*gorm.DB)
	data := dtos.BackupDTO{}
	jsonFile, err := os.Open(fmt.Sprintf("%v/%v", dto.Path, dto.FileName))

	if err != nil {
		responses.ErrorJSON(ctx, http.StatusBadRequest, err)
		return
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &data)

	for _, item := range data.Todo {
		item.ID = 0
		item.CreatedAt = time.Now()
		item.UpdatedAt = time.Now()
	}
	for _, item := range data.Users {
		item.ID = 0
		item.CreatedAt = time.Now()
		item.UpdatedAt = time.Now()
	}
	if err := cc.TodoService.WithTrx(trx).BulkCreateTodo(data.Todo); err != nil {
		responses.ErrorJSON(ctx, http.StatusInternalServerError, err)
		return
	}

	if err := cc.UserService.WithTrx(trx).BulkCreateUser(data.Users); err != nil {
		responses.ErrorJSON(ctx, http.StatusInternalServerError, err)
		return
	}
	defer jsonFile.Close()
	responses.JSON(ctx, http.StatusOK, "Data imported successfully")
}
