package controllers

import (
	"boilerplate-api/api/services"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/responses"
	"fmt"
	"net/http"
	"os/user"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// TodoController struct
type TodoController struct {
	logger      infrastructure.Logger
	TodoService services.TodoService
}

// NewTodoController constructor
func NewTodoController(
	logger infrastructure.Logger,
	TodoService services.TodoService,
) TodoController {
	return TodoController{
		logger:      logger,
		TodoService: TodoService,
	}
}

// CreateTodo Create Todo
func (cc TodoController) CreateTodo(c *gin.Context) {
	todo := models.Todo{}

	if err := c.ShouldBindJSON(&todo); err != nil {
		cc.logger.Zap.Error("Error [CreateTodo] (ShouldBindJson) : ", err)
		err := errors.BadRequest.Wrap(err, "Failed to bind Todo")
		responses.HandleError(c, err)
		return
	}

	if _, err := cc.TodoService.CreateTodo(todo); err != nil {
		cc.logger.Zap.Error("Error [CreateTodo] [db CreateTodo]: ", err.Error())
		err := errors.BadRequest.Wrap(err, "Failed To Create Todo")
		responses.HandleError(c, err)
		return
	}

	responses.SuccessJSON(c, http.StatusOK, "Todo Created Successfully")
}

// GetAllTodo Get All Todo
func (cc TodoController) GetAllTodo(c *gin.Context) {

	todos, count, err := cc.TodoService.GetAllTodo()
	if err != nil {
		cc.logger.Zap.Error("Error finding Todo records", err.Error())
		err := errors.InternalError.Wrap(err, "Failed To Find Todo")
		responses.HandleError(c, err)
		return
	}

	responses.JSONCount(c, http.StatusOK, todos, count)
}

func (cc TodoController) GetOneTodo(c *gin.Context) {
	ID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	todo, err := cc.TodoService.GetOneTodo(ID)
	if err != nil {
		cc.logger.Zap.Error("Error [GetOneTodo] [db GetOneTodo]: ", err.Error())
		err := errors.InternalError.Wrap(err, "Failed To Find Todo")
		responses.HandleError(c, err)
		return
	}

	responses.JSON(c, http.StatusOK, todo)
}

func (cc TodoController) ExportExcel(c *gin.Context) {
	// filePath := "exported_data.xlsx"
	// path := utils.GetPath()
	// fmt.Println(path, "UTILSSS-->")
	// if path != "" {
	// 	filePath = fmt.Sprintf("%v/exported_data.xlsx", path)
	// }
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Failed to retrieve user information:", err)
		return
	}

	// Construct the desktop path
	desktopPath := filepath.Join(usr.HomeDir, "Desktop")
	filePath := filepath.Join(desktopPath, "output.xlsx")
	f := excelize.NewFile()
	// defer func() {
	// 	if err := f.Close(); err != nil {
	// 		fmt.Println(err)
	// 		responses.ErrorJSON(c, http.StatusOK, err.Error())
	// 		return

	// 	}
	// }()
	getTodo, _, _ := cc.TodoService.GetAllTodo()
	todo := models.Todo{}
	// Create a new sheet.
	sheetName := todo.TableName()
	index, err := f.NewSheet(sheetName)
	if err != nil {
		fmt.Println("err1-->")
		responses.ErrorJSON(c, http.StatusOK, err.Error())
		return
	}
	// Set value of a cell.
	f.SetCellValue(sheetName, "A1", "ID")
	f.SetCellValue(sheetName, "B1", "TASK")
	f.SetCellValue(sheetName, "C1", "IS COMPLETED")
	for i, item := range getTodo {
		shell1 := fmt.Sprintf("A%v", i+2)
		shell2 := fmt.Sprintf("B%v", i+2)
		shell3 := fmt.Sprintf("C%v", i+2)
		f.SetCellValue(sheetName, shell1, item.ID)
		f.SetCellValue(sheetName, shell2, item.Task)
		f.SetCellValue(sheetName, shell3, item.Is_completed)
	}
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)
	fmt.Println(filePath, "FILEPATH-->")
	// Save spreadsheet by the given path.
	if err := f.SaveAs(filePath); err != nil {
		fmt.Println(err)
		fmt.Println("err2-->")
		responses.ErrorJSON(c, http.StatusOK, err.Error())
		return
	}
	if err := f.Close(); err != nil {
		fmt.Println(err)
		fmt.Println("err3-->")

		responses.ErrorJSON(c, http.StatusOK, err.Error())
		return

	}
	defer responses.JSON(c, http.StatusOK, "Exported")

}
