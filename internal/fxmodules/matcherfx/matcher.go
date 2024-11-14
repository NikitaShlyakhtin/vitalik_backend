package matcherfx

import (
	"go.uber.org/fx"
	"vitalik_backend/internal/pkg/services/matcher"
)

var Module = fx.Module("matcherfx",
	fx.Provide(matcher.NewMatcher),
	fx.Invoke(startMatcher),
)
