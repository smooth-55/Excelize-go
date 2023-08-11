package wbs

import (
	"boilerplate-api/api/middlewares"
	"boilerplate-api/constants"
	"boilerplate-api/infrastructure"
)

// MessageRoutes struct
type ServerRoutes struct {
	logger infrastructure.Logger
	router infrastructure.Router
	// jwtMiddleware       middlewares.JWTAuthMiddleWare
	trxMiddleware       middlewares.DBTransactionMiddleware
	rateLimitMiddleware middlewares.RateLimitMiddleware
	handler             *Handler
}

func NewServerRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	jwtMiddleware middlewares.JWTAuthMiddleWare,
	trxMiddleware middlewares.DBTransactionMiddleware,
	rateLimitMiddleware middlewares.RateLimitMiddleware,
	handler *Handler,

) ServerRoutes {
	return ServerRoutes{
		router: router,
		logger: logger,
		// jwtMiddleware:       jwtMiddleware,
		trxMiddleware:       trxMiddleware,
		rateLimitMiddleware: rateLimitMiddleware,
		handler:             handler,
	}
}

// Setup user routes
func (i ServerRoutes) Setup() {
	i.logger.Zap.Info(" Setting up user routes")
	ws := i.router.Gin.Group("/ws").Use(i.rateLimitMiddleware.
		HandleRateLimit(constants.BasicRateLimit, constants.BasicPeriod),
	) //.Use(i.jwtMiddleware.Handle())
	{
		// ws.POST("/create-room", i.handler.CreateRoom)
		ws.GET("/join-room/:roomId", i.handler.JoinRoom)
		ws.GET("/join-room", i.handler.JoinRoomTest)

		// ws.GET("/get-rooms", i.handler.GetRooms)
		// ws.GET("/get-clients/:roomId", i.handler.GetClients)
	}
}
