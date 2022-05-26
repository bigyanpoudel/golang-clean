package routes

import "go.uber.org/fx"

// Module exports dependency to container
var Module = fx.Options(
	fx.Provide(NewRoutes),
	fx.Provide(NewPostRoute),
	fx.Provide(NewUserRoute),
)

//Multiple routes
type Routes []Route

//routes interface
type Route interface {
	Setup()
}

// NewRoutes sets up routes
func NewRoutes(
	userRoutes UserRoutes,
	postRoutes PostRoute,
) Routes {
	return Routes{
		userRoutes,
		postRoutes,
	}
}

func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
