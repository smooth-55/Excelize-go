package controllers

import (
	"boilerplate-api/api/services"
	"boilerplate-api/api/validators"
	"boilerplate-api/constants"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/responses"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
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
	_userId := c.MustGet(constants.UserID).(string)
	userId, _ := strconv.ParseInt(_userId, 10, 64)

	if userId == 0 {
		cc.logger.Zap.Error("User Id required")
		err := errors.BadRequest.New("User Id required")
		responses.HandleError(c, err)
		return
	}
	data, count, _ := cc.messageService.GetMyConversations(userId)

	responses.JSONCount(c, http.StatusOK, data, count)

}
