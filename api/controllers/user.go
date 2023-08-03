package controllers

import (
	"boilerplate-api/api/services"
	"boilerplate-api/api/validators"
	"boilerplate-api/constants"
	"boilerplate-api/dtos"
	"boilerplate-api/errors"
	"boilerplate-api/infrastructure"
	"boilerplate-api/models"
	"boilerplate-api/paginations"
	"boilerplate-api/responses"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	logger      infrastructure.Logger
	userService services.UserService
	env         infrastructure.Env
	validator   validators.UserValidator
}

// NewUserController Creates New user controller
func NewUserController(
	logger infrastructure.Logger,
	userService services.UserService,
	env infrastructure.Env,
	validator validators.UserValidator,
) UserController {
	return UserController{
		logger:      logger,
		userService: userService,
		env:         env,
		validator:   validator,
	}
}

// CreateUser Create User
// @Summary				Create User
// @Description			Create User
// @Param				data body dtos.CreateUserRequestData true "Enter JSON"
// @Param 				Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Produce				application/json
// @Tags				User
// @Success				200 {object} responses.Success "OK"
// @Failure      		400 {object} responses.Error
// @Failure      		500 {object} responses.Error
// @Router				/users [post]
func (cc UserController) CreateUser(c *gin.Context) {
	reqData := dtos.CreateUserRequestData{}
	trx := c.MustGet(constants.DBTransaction).(*gorm.DB)

	if err := c.ShouldBindJSON(&reqData); err != nil {
		cc.logger.Zap.Error("Error [CreateUser] (ShouldBindJson) : ", err)
		err := errors.BadRequest.Wrap(err, "Failed to bind user data")
		responses.HandleError(c, err)
		return
	}
	if validationErr := cc.validator.Validate.Struct(reqData); validationErr != nil {
		err := errors.BadRequest.Wrap(validationErr, "Validation error")
		err = errors.SetCustomMessage(err, "Invalid input information")
		err = errors.AddErrorContextBlock(err, cc.validator.GenerateValidationResponse(validationErr))
		responses.HandleError(c, err)
		return
	}

	if _, err := cc.userService.GetOneUserWithEmail(reqData.Email); err == nil {
		cc.logger.Zap.Error("Error [CreateUser] [db CreateUser]: User with this email already exists")
		err := errors.BadRequest.New("User with this email already exists")
		responses.HandleError(c, err)
		return
	}

	if _, err := cc.userService.GetOneUserWithUsername(reqData.Username); err == nil {
		cc.logger.Zap.Error("Error [CreateUser] [db CreateUser]: User with this username already exists")
		err := errors.BadRequest.New("User with this username already exists")
		responses.HandleError(c, err)
		return
	}
	user := reqData.GetUserObj()
	if err := cc.userService.WithTrx(trx).CreateUser(user); err != nil {
		cc.logger.Zap.Error("Error [CreateUser] [db CreateUser]: ", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to create user")
		responses.HandleError(c, err)
		return
	}

	responses.SuccessJSON(c, http.StatusOK, "User Created Successfully")
}

// GetAllUsers Get All User
// @Summary				Get all User.
// @Param				page_size query string false "10"
// @Param				page query string false "Page no" "1"
// @Param				keyword query string false "search by name"
// @Param				Keyword2 query string false "search by type"
// @Description			Return all the User
// @Produce				application/json
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Tags				User
// @Success 			200 {array} responses.DataCount{data=[]dtos.GetUserResponse}
// @Failure      		500 {object} responses.Error
// @Router				/users [get]
func (cc UserController) GetAllUsers(c *gin.Context) {
	pagination := paginations.BuildPagination[*paginations.UserPagination](c)

	users, count, err := cc.userService.GetAllUsers(*pagination)
	if err != nil {
		cc.logger.Zap.Error("Error finding user records", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to get users data")
		responses.HandleError(c, err)
		return
	}

	responses.JSONCount(c, http.StatusOK, users, count)
}

// GetUserProfile Returns logged-in user profile
// @Summary				Get one user by id
// @Description			Get one user by id
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Produce				application/json
// @Tags				User
// @Success 			200 {array} responses.Data{data=dtos.GetUserResponse}
// @Failure      		500 {object} responses.Error
// @Router				/profile [get]
func (cc UserController) GetUserProfile(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.MustGet(constants.UserID).(string), 10, 64)

	user, followFollowing, err := cc.userService.GetOneUser(userId)
	if err != nil {
		cc.logger.Zap.Error("Error finding user profile", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to get users profile data")
		responses.HandleError(c, err)
		return
	}
	resp := make(map[string]interface{})
	resp["user"] = user
	resp["data"] = followFollowing
	responses.JSON(c, http.StatusOK, resp)
}

func (cc UserController) FollowUser(c *gin.Context) {
	fmt.Println("------->>>>iam here")
	trx := c.MustGet(constants.DBTransaction).(*gorm.DB)

	userId, _ := strconv.ParseInt(c.MustGet(constants.UserID).(string), 10, 64)
	reqData := dtos.FollowUser{}
	if err := c.ShouldBindJSON(&reqData); err != nil {
		cc.logger.Zap.Error("Error [CreateUser] (ShouldBindJson) : ", err)
		err := errors.BadRequest.Wrap(err, "Failed to bind user data")
		responses.HandleError(c, err)
		return
	}
	if userId == reqData.FollowUserId {
		cc.logger.Zap.Error("You can't follow yourself")
		responses.ErrorJSON(c, http.StatusBadRequest, "You can't follow yourself")
		return
	}
	user, _, err := cc.userService.GetOneUser(userId)
	if err != nil {
		cc.logger.Zap.Error("Error finding user profile", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to get users profile data")
		responses.HandleError(c, err)
		return
	}
	if _, _, err := cc.userService.GetOneUser(reqData.FollowUserId); err != nil {
		cc.logger.Zap.Error("Error finding user profile", err.Error())
		erMessage := fmt.Sprintf("User with id: %v, does'nt exists", reqData.FollowUserId)
		err := errors.InternalError.Wrap(err, erMessage)
		responses.HandleError(c, err)
		return
	}
	follow := models.FollowUser{}
	follow.FollowedById = user.ID
	follow.FollowedToId = reqData.FollowUserId
	if user.IsPrivate {
		follow.IsApproved = true
	}
	if err := cc.userService.WithTrx(trx).FollowUser(follow); err != nil {
		cc.logger.Zap.Error("Error finding user profile", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to get users profile data")
		responses.HandleError(c, err)
		return
	}
	responses.JSON(c, http.StatusOK, user)
}

func (cc UserController) FollowSuggestions(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.MustGet(constants.UserID).(string), 10, 64)
	if _, _, err := cc.userService.GetOneUser(userId); err != nil {
		cc.logger.Zap.Error("Error finding user profile", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to get users profile data")
		responses.HandleError(c, err)
		return
	}
	suggestedUsers, err := cc.userService.FollowSuggestions(userId)
	if err != nil {
		cc.logger.Zap.Error("Error finding follow suggestion", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to get follwer suggestions")
		responses.HandleError(c, err)
		return
	}
	responses.JSON(c, http.StatusOK, suggestedUsers)
}

func (cc UserController) GetTwoWayFollowers(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.MustGet(constants.UserID).(string), 10, 64)
	if _, _, err := cc.userService.GetOneUser(userId); err != nil {
		cc.logger.Zap.Error("Error finding user profile", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to get users profile data")
		responses.HandleError(c, err)
		return
	}
	suggestedUsers, err := cc.userService.GetTwoWayFollowers(userId)
	if err != nil {
		cc.logger.Zap.Error("Error finding follow suggestion", err.Error())
		err := errors.InternalError.Wrap(err, "Failed to get follwer suggestions")
		responses.HandleError(c, err)
		return
	}
	responses.JSON(c, http.StatusOK, suggestedUsers)
}
