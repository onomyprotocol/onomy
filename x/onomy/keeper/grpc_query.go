package keeper

import (
	"github.com/onomyprotocol/onomy/x/onomy/types"
)

var _ types.QueryServer = Keeper{}
