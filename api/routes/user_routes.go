package routes

import (
	"go-clean-api/api/controller"
	"go-clean-api/api/middleware"
	"go-clean-api/infrastructure"
)

type UserRoutes struct {
	handler       infrastructure.RequestHandler
	controller    controller.UserController
	authorization middleware.FirebaeAuth
}

func NewUserRoute(handler infrastructure.RequestHandler, c controller.UserController, authorization middleware.FirebaeAuth) UserRoutes {
	return UserRoutes{
		handler:       handler,
		controller:    c,
		authorization: authorization,
	}
}

func (u UserRoutes) Setup() {
	user := u.handler.Gin.Group("/api/user/")
	{
		user.POST("register", u.controller.Register())
		user.GET("read", u.controller.GetAllUser())
		user.POST("login", u.controller.LoginHandler())
		user.GET("profile", u.authorization.Handle(), u.controller.GetUserProfile())
		user.PUT("change-password", u.authorization.Handle(), u.controller.GetUserProfile())
		user.POST("verify", u.controller.VerifyUser())
		user.POST("verify-email", u.controller.UsersByEmail())
		user.POST("upload", u.controller.FileUpload())
	}
}
