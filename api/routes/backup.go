package routes

import (
	"boilerplate-api/api/controllers"
	"boilerplate-api/api/middlewares"
	"boilerplate-api/infrastructure"
)

// BackupRoutes struct
type BackupRoutes struct {
	logger           infrastructure.Logger
	router           infrastructure.Router
	backupController controllers.BackupController
	middleware       middlewares.FirebaseAuthMiddleware
	trxMiddleware    middlewares.DBTransactionMiddleware
}

// NewBackupRoutes creates new Todo controller
func NewBackupRoutes(
	logger infrastructure.Logger,
	router infrastructure.Router,
	backupController controllers.BackupController,
	middleware middlewares.FirebaseAuthMiddleware,
	trxMiddleware middlewares.DBTransactionMiddleware,
) BackupRoutes {
	return BackupRoutes{
		router:           router,
		logger:           logger,
		backupController: backupController,
		middleware:       middleware,
		trxMiddleware:    trxMiddleware,
	}
}

// Setup todo routes
func (c BackupRoutes) Setup() {
	c.logger.Zap.Info(" Setting up Todo routes")
	todo := c.router.Gin.Group("/backup")
	{
		todo.GET("/export", c.backupController.BackupData)
		todo.POST("/import", c.trxMiddleware.DBTransactionHandle(), c.backupController.RestoreData)
	}
}
