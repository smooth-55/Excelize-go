package controllers

import (
	"boilerplate-api/api/services"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/responses"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
