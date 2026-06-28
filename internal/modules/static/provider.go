package static

import "github.com/google/wire"

var Provider wire.ProviderSet = wire.NewSet(
	NewStaticHandler,
	NewStaticRoute,
)
