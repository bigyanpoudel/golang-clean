package infrastructure

import "go.uber.org/fx"

// Module exported for initializing application
var Module = fx.Options(
	fx.Provide(NewEnv),
	fx.Provide(NewDatabase),
	fx.Provide(NewLogger),
	fx.Provide(NewRequestHandler),
	fx.Provide(NewMigration),
	fx.Provide(NewFBApp),
	fx.Provide(NewFBAuth),
	fx.Provide(ConnectAws),
)
