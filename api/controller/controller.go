package controller

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewPostController),
	fx.Provide(NewUserController),
)
