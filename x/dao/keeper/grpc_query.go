package keeper

import (
	"github.com/onomyprotocol/onomy/x/dao/types"
)

var _ types.QueryServer = Keeper{}
