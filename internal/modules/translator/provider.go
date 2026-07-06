package translator

import "github.com/google/wire"

var Provider wire.ProviderSet = wire.NewSet(
	NewHandler,
	NewRoute,
)
