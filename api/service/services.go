package service

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewPostService),
	fx.Provide(NewUserService),
	fx.Provide(NewFirebaseService),
	fx.Provide(NewAWSService),
)
