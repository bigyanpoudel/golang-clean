package middleware

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewMiddlewares),
	fx.Provide(NewAuthMiddlerware),
	fx.Provide(NewFirebaseAuth),
)

type IMiddleware interface {
	SetUp()
}

type Middleware []IMiddleware

func NewMiddlewares(
	authMiddleware AuthMiddleWare,
) Middleware {
	return Middleware{
		authMiddleware,
	}
}

func (m Middleware) SetUp() {
	for _, middleware := range m {
		middleware.SetUp()
	}
}
