package routes

import "go.uber.org/fx"

// Module exports dependency to container
var Module = fx.Options(
	fx.Provide(NewRoutes),
	fx.Provide(NewTodoRoutes),
	fx.Provide(NewUtilityRoutes),
	fx.Provide(NewDocsRoutes),
	fx.Provide(NewJwtAuthRoutes),
	fx.Provide(NewUserRoutes),
	fx.Provide(NewBackupRoutes),
)

// Routes contains multiple routes
type Routes []Route

// Route interface
type Route interface {
	Setup()
}

// NewRoutes sets up routes
func NewRoutes(
	todoRoutes TodoRoutes,
	utilityRoutes UtilityRoutes,
	userRoutes UserRoutes,
	jwtAuthRoutes JwtAuthRoutes,
	docsRoutes DocsRoutes,
	backupRoutes BackupRoutes,
) Routes {
	return Routes{
		todoRoutes,
		utilityRoutes,
		jwtAuthRoutes,
		userRoutes,
		docsRoutes,
		backupRoutes,
	}
}

// Setup all the route
func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
