// Package simapp contains functions to run the chain as simapp.
package simapp

import (
	"time"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/tendermint/spm/cosmoscmd"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	tmdb "github.com/tendermint/tm-db"

	"github.com/onomyprotocol/onomy/app"
)

const (
	consensusParamsBlockMaxBytes       = 200000
	consensusParamsBlockMaxBytesMaxGas = 2000000

	consensusEvidenceMaxAgeNumBlocks     = 302400
	consensusEvidenceMaxAgeDurationHours = 504
	consensusEvidenceMaxBytes            = 10000
)

// New creates application instance with in-memory database and disabled logging.
func New(dir string) cosmoscmd.App {
	db := tmdb.NewMemDB()
	logger := log.NewNopLogger()

	encoding := cosmoscmd.MakeEncodingConfig(app.ModuleBasics)

	a := app.New(logger, db, nil, true, map[int64]bool{}, dir, 0, encoding,
		simapp.EmptyAppOptions{})

	defaultConsensusParams := &abci.ConsensusParams{
		Block: &abci.BlockParams{
			MaxBytes: consensusParamsBlockMaxBytes,
			MaxGas:   consensusParamsBlockMaxBytesMaxGas,
		},
		Evidence: &tmproto.EvidenceParams{
			MaxAgeNumBlocks: consensusEvidenceMaxAgeNumBlocks,
			MaxAgeDuration:  consensusEvidenceMaxAgeDurationHours * time.Hour, // 3 weeks is the max duration
			MaxBytes:        consensusEvidenceMaxBytes,
		},
		Validator: &tmproto.ValidatorParams{
			PubKeyTypes: []string{
				tmtypes.ABCIPubKeyTypeEd25519,
			},
		},
	}

	// InitChain updates deliverState which is required when app.NewContext is called
	a.InitChain(abci.RequestInitChain{
		ConsensusParams: defaultConsensusParams,
		AppStateBytes:   []byte("{}"),
	})
	return a
}
