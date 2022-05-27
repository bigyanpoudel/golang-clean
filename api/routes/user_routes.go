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
	pagination    middleware.Pagination
}

func NewUserRoute(handler infrastructure.RequestHandler, c controller.UserController, authorization middleware.FirebaeAuth, p middleware.Pagination) UserRoutes {
	return UserRoutes{
		handler:       handler,
		controller:    c,
		authorization: authorization,
		pagination:    p,
	}
}

func (u UserRoutes) Setup() {
	user := u.handler.Gin.Group("/api/user/")
	{
		user.POST("register", u.controller.Register())
		user.GET("read", u.pagination.IncludePagination(), u.controller.GetAllUser())
		user.POST("login", u.controller.LoginHandler())
		user.GET("profile", u.authorization.Handle(), u.controller.GetUserProfile())
		user.PUT("change-password", u.authorization.Handle(), u.controller.GetUserProfile())
		user.POST("verify", u.controller.VerifyUser())
		user.POST("verify-email", u.controller.UsersByEmail())
		user.POST("upload", u.controller.FileUpload())
		user.POST("search", u.pagination.IncludePagination(), u.controller.SearchUser())
		user.GET("/read/users", u.pagination.IncludePagination(), u.controller.GetUsers())
	}
}
