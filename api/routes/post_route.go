package routes

import (
	"go-clean-api/api/controller"
	"go-clean-api/api/middleware"
	"go-clean-api/infrastructure"
)

type PostRoute struct {
	handler       infrastructure.RequestHandler
	controller    controller.PostController
	Logger        infrastructure.Logger
	authorization middleware.FirebaeAuth
}

func NewPostRoute(controller controller.PostController, handler infrastructure.RequestHandler, logger infrastructure.Logger, a middleware.FirebaeAuth) PostRoute {
	return PostRoute{
		controller:    controller,
		handler:       handler,
		Logger:        logger,
		authorization: a,
	}
}

func (p PostRoute) Setup() {
	p.Logger.Zap.Info("post routes")
	post := p.handler.Gin.Group("/api/post/")
	{
		post.POST("create", p.controller.CreatePost())
		post.GET("read", p.controller.GetAllPost())
		post.GET("read/:id", p.controller.GetPostById())
		post.PUT("upload/:id", p.controller.UploadPostImage())
		post.GET("user-post", p.authorization.Handle(), p.controller.GetPostByUserId())
	}
}
