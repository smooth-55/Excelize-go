package wbs

import "go.uber.org/fx"

// Module exports dependency

var Module = fx.Options(
	fx.Provide(NewWSRoute),
	fx.Provide(NewServerRoutes),
	fx.Provide(NewHandler),
	fx.Provide(NewHub),
	// fx.Provide(NewClient),
)

// Routes contains multiple routes
type WSRoutes []WSRoute

// Route interface
type WSRoute interface {
	Setup()
}

// NewRoutes sets up routes
func NewWSRoute(
	server ServerRoutes,
) WSRoutes {
	return WSRoutes{
		server,
	}
}

// Setup all the route
func (r WSRoutes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
