package v2_0_0

import "cosmossdk.io/math"

type UnbondingFail struct {
	Delegator string
	Validator string

	CreationHeight int64
	Balance        math.Int

	Index       int
	UnbondingId uint64
}
