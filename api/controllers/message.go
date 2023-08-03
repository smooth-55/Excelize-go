package controllers

import (
	"boilerplate-api/api/services"
	"boilerplate-api/api/validators"
	"boilerplate-api/constants"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/responses"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MessageController struct {
	logger         infrastructure.Logger
	messageService services.MessageService
	env            infrastructure.Env
	validator      validators.UserValidator
}

// NewUserController Creates New user controller
func NewMessageController(
	logger infrastructure.Logger,
	messageService services.MessageService,
	env infrastructure.Env,
	validator validators.UserValidator,
) MessageController {
	return MessageController{
		logger:         logger,
		messageService: messageService,
		env:            env,
		validator:      validator,
	}
}

func (cc MessageController) GetMyConversations(c *gin.Context) {
	userId := c.GetInt64(constants.UserID)
	if userId == 0 {
		cc.logger.Zap.Error("User Id required")
		err := errors.BadRequest.New("User Id required")
		responses.HandleError(c, err)
		return
	}
	data, count, err := cc.messageService.GetMyConversations(userId)
	if err != nil {
		cc.logger.Zap.Error("Error [CreateUser] [db CreateUser]: User with this email already exists")
		err := errors.BadRequest.New("User with this email already exists")
		responses.HandleError(c, err)
		return
	}
	responses.JSONCount(c, http.StatusOK, data, count)

}
