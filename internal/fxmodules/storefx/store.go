package storefx

import (
	"go.uber.org/fx"
	"vitalik_backend/internal/dependencies"
	"vitalik_backend/internal/pkg/services/store"
)

var Module = fx.Module("storefx",
	fx.Provide(
		fx.Annotate(store.NewStore, fx.As(new(dependencies.IStore))),
	),
	fx.Provide(
		fx.Private,
		newPgxPool,
	),
	fx.Invoke(startPgxPool),
)
