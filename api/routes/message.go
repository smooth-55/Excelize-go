package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/api/middlewares"
	"boilerplate-api/constants"
	"boilerplate-api/infrastructure"
)

// MessageRoutes struct
type MessageRoutes struct {
	logger              infrastructure.Logger
	router              infrastructure.Router
	messageController   controllers.MessageController
	middleware          middlewares.FirebaseAuthMiddleware
	jwtMiddleware       middlewares.JWTAuthMiddleWare
	trxMiddleware       middlewares.DBTransactionMiddleware
	rateLimitMiddleware middlewares.RateLimitMiddleware
}

func NewMessageRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	messageController controllers.MessageController,
	jwtMiddleware middlewares.JWTAuthMiddleWare,
	trxMiddleware middlewares.DBTransactionMiddleware,
	rateLimitMiddleware middlewares.RateLimitMiddleware,

) MessageRoutes {
	return MessageRoutes{
		router:              router,
		logger:              logger,
		messageController:   messageController,
		jwtMiddleware:       jwtMiddleware,
		trxMiddleware:       trxMiddleware,
		rateLimitMiddleware: rateLimitMiddleware,
	}
}

// Setup user routes
func (i MessageRoutes) Setup() {
	i.logger.Zap.Info(" Setting up user routes")
	message := i.router.Gin.Group("/messages").
		Use(i.rateLimitMiddleware.
			HandleRateLimit(constants.BasicRateLimit, constants.BasicPeriod),
		).
		Use(i.jwtMiddleware.Handle())
	{
		message.GET("get-all-conversations", i.messageController.GetMyConversations)
	}
}
