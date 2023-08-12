package controllers

import (
	"boilerplate-api/api/services"
	"boilerplate-api/api/validators"
	"boilerplate-api/constants"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/responses"
	"fmt"
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

func (cc MessageController) GetAllMessagesByRoomId(c *gin.Context) {
	roomId := c.Param("room_id")
	_roomId, _ := strconv.ParseInt(roomId, 10, 64)
	data, err := cc.messageService.GetAllMessagesByRoomId(_roomId)
	if err != nil {
		cc.logger.Zap.Error("Error [GetAllMessagesByRoomId]", err)
		err := errors.InternalError.New(fmt.Sprintf("Something went wrong [GetAllMessagesByRoomId] %v", err))
		responses.HandleError(c, err)
		return
	}
	responses.JSON(c, http.StatusOK, data)

}
