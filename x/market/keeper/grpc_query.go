package keeper

import (
	"github.com/onomyprotocol/onomy/x/market/types"
)

var _ types.QueryServer = Keeper{}
