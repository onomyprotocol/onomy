// Package simapp contains utils to bootstrap the chain.
package simapp

import (
	"bytes"
	"encoding/json"
	"testing"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/starport/starport/pkg/cosmoscmd"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	gravitytypes "github.com/onomyprotocol/cosmos-gravity-bridge/module/x/gravity/types"
	"github.com/onomyprotocol/onomy/app"
)

// The SimApp is OnomyApp wrapper with the advance testing capabilities.
type SimApp struct {
	onomyApp *app.OnomyApp
}

// GenesisState of the blockchain is represented here as a map of raw json
// messages key'd by a identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type GenesisState map[string]json.RawMessage

// Option is an option pattern function used fot the test simapp customisations.
type Option struct {
	before func(*SimApp, GenesisState) (*SimApp, GenesisState)
	after  func(*SimApp) *SimApp
}

// ValReq is simplified struct for the validator creation.
type ValReq struct {
	Balance      sdk.Coins
	SelfBondCoin sdk.Coin
	Commission   stakingtypes.CommissionRates
	Reward       sdk.Coin
}

// WithGenesisAccountsAndBalances returns genesis override Option for initial balances.
func WithGenesisAccountsAndBalances(balances ...banktypes.Balance) Option {
	return Option{
		before: func(simApp *SimApp, genState GenesisState) (*SimApp, GenesisState) {
			accounts := make([]authtypes.GenesisAccount, 0, len(balances))
			totalSupply := sdk.NewCoins()
			for _, balance := range balances {
				accounts = append(accounts, &authtypes.BaseAccount{
					Address: balance.Address,
				})
				totalSupply = totalSupply.Add(balance.Coins...)
			}

			authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), accounts)
			genState[authtypes.ModuleName] = simApp.onomyApp.AppCodec().MustMarshalJSON(authGenesis)

			s := totalSupply.String()
			_ = s

			bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{})
			genState[banktypes.ModuleName] = simApp.onomyApp.AppCodec().MustMarshalJSON(bankGenesis)

			return simApp, genState
		},
	}
}

// WithGenesisOverride returns genesis override ConfigOption.
func WithGenesisOverride(override func(map[string]json.RawMessage) map[string]json.RawMessage) Option {
	return Option{
		before: func(simApp *SimApp, genState GenesisState) (*SimApp, GenesisState) {
			genState = override(genState)
			return simApp, genState
		},
	}
}

// WithAppCommit commits the app state after the initialisation.
func WithAppCommit() Option {
	return Option{
		after: func(simApp *SimApp) *SimApp {
			simApp.onomyApp.Commit()
			return simApp
		},
	}
}

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState(cdc codec.JSONCodec) GenesisState {
	return app.ModuleBasics.DefaultGenesis(cdc)
}

// SetupWithValidators create new simApp with the defined list of the validators.
func SetupWithValidators(t *testing.T, vals map[string]ValReq, opts ...Option) (*SimApp, map[string]*secp256k1.PrivKey) {
	t.Helper()

	// prepare account genesis params
	balances := make([]banktypes.Balance, 0, len(vals))
	privateKeys := make(map[string]*secp256k1.PrivKey, len(vals))
	for moniker := range vals {
		privateKey := secp256k1.GenPrivKey()
		privateKeys[moniker] = privateKey
		address := sdk.AccAddress(privateKey.PubKey().Address())
		balance := vals[moniker].Balance
		if balance.AmountOf(vals[moniker].SelfBondCoin.Denom).LT(vals[moniker].SelfBondCoin.Amount) {
			balance = balance.Add(vals[moniker].SelfBondCoin)
		}
		balances = append(balances, banktypes.Balance{
			Address: address.String(),
			Coins:   balance,
		})
	}

	opts = append(opts, WithGenesisAccountsAndBalances(balances...), WithAppCommit())

	simApp := Setup(opts...)

	simApp.BeginNextBlock()
	for moniker, val := range vals {
		description := stakingtypes.Description{Moniker: moniker}
		simApp.CreateValidator(t, val.SelfBondCoin, description, val.Commission, sdk.OneInt(), privateKeys[moniker])
	}
	return simApp, privateKeys
}

// Setup initializes a new SimApp. A Nop logger is set in SimApp.
func Setup(opts ...Option) *SimApp {
	onomyApp, genesisState := setup(5) // nolint:gomnd //test constant

	simApp := &SimApp{
		onomyApp: onomyApp,
	}

	for _, opt := range opts {
		if opt.before != nil {
			simApp, genesisState = opt.before(simApp, genesisState)
		}
	}

	// init chain must be called to stop deliverState from being nil
	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		panic(err)
	}

	// Initialize the chain
	simApp.onomyApp.InitChain(
		abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: simapp.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	for _, opt := range opts {
		if opt.after != nil {
			simApp = opt.after(simApp)
		}
	}

	return simApp
}

// OnomyApp returns OnomyApp from the SimApp.
func (s *SimApp) OnomyApp() *app.OnomyApp {
	return s.onomyApp
}

// BeginNextBlock begins new SimApp block.
func (s *SimApp) BeginNextBlock() {
	s.onomyApp.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: s.onomyApp.LastBlockHeight() + 1}})
}

// EndBlockAndCommit ends the current block and commit the state.
func (s *SimApp) EndBlockAndCommit(ctx sdk.Context) {
	s.onomyApp.EndBlocker(ctx, abci.RequestEndBlock{Height: ctx.BlockHeight()})
	s.onomyApp.Commit()
}

// EndBlock ends the current block.
func (s *SimApp) EndBlock(ctx sdk.Context) {
	s.onomyApp.EndBlocker(ctx, abci.RequestEndBlock{Height: ctx.BlockHeight()})
}

// NewContext returns empty sdk context for the SimApp.
func (s *SimApp) NewContext() sdk.Context {
	return s.onomyApp.BaseApp.NewContext(true, tmproto.Header{})
}

// CurrentContext returns current context for the SimApp.
func (s *SimApp) CurrentContext() sdk.Context {
	return s.onomyApp.BaseApp.NewContext(true, tmproto.Header{Height: s.onomyApp.LastBlockHeight()})
}

// NewNextContext creates next block sdk context for the SimApp.
func (s *SimApp) NewNextContext() sdk.Context {
	header := tmproto.Header{Height: s.onomyApp.LastBlockHeight() + 1}
	return s.onomyApp.BaseApp.NewContext(false, header)
}

// CreateValidator creates the validator.
func (s *SimApp) CreateValidator(
	t *testing.T,
	selfDelegation sdk.Coin,
	description stakingtypes.Description,
	commission stakingtypes.CommissionRates,
	minSelfDelegation sdk.Int,
	priv cryptotypes.PrivKey,
) {
	t.Helper()

	address := sdk.AccAddress(priv.PubKey().Address())
	valAddress := sdk.ValAddress(address)
	messages := make([]sdk.Msg, 0)
	createValidatorMsg, err := stakingtypes.NewMsgCreateValidator(
		valAddress, ed25519.GenPrivKey().PubKey(), selfDelegation, description, commission, minSelfDelegation,
	)
	require.NoError(t, err)
	messages = append(messages, createValidatorMsg)

	setOrchestratorAddressMsg := &gravitytypes.MsgSetOrchestratorAddress{
		Validator:    createValidatorMsg.ValidatorAddress,
		Orchestrator: address.String(),
		EthAddress:   gethcommon.BytesToAddress(bytes.Repeat([]byte{byte(1)}, 20)).String(), // nolint:gomnd // eth address
	}

	if selfDelegation.Amount.GTE(sdk.DefaultPowerReduction) {
		messages = append(messages, setOrchestratorAddressMsg)
	}

	s.sendTx(t, priv, messages...)
}

// CreateProposal creates a new proposal with the provided contant.
func (s *SimApp) CreateProposal(
	t *testing.T,
	content govtypes.Content,
	deposit sdk.Coin,
	priv cryptotypes.PrivKey,
) {
	t.Helper()

	address := sdk.AccAddress(priv.PubKey().Address())
	msg, err := govtypes.NewMsgSubmitProposal(content, sdk.NewCoins(deposit), address)
	require.NoError(t, err)

	s.sendTx(t, priv, msg)
}

// VoteProposal votes for the proposal.
func (s *SimApp) VoteProposal(
	t *testing.T,
	proposalID uint64,
	option govtypes.VoteOption,
	priv cryptotypes.PrivKey,
) {
	t.Helper()

	address := sdk.AccAddress(priv.PubKey().Address())

	msg := govtypes.NewMsgVote(address, proposalID, option)
	s.sendTx(t, priv, msg)
}

// GenAccountAddress generates random account.
func GenAccountAddress() sdk.AccAddress {
	pk := ed25519.GenPrivKey().PubKey()
	return sdk.AccAddress(pk.Address())
}

func (s *SimApp) sendTx(t *testing.T, priv cryptotypes.PrivKey, messages ...sdk.Msg) {
	t.Helper()

	address := sdk.AccAddress(priv.PubKey().Address())
	account := s.onomyApp.AccountKeeper.GetAccount(s.NewContext(), address)
	accountNum := account.GetAccountNumber()
	accountSeq := account.GetSequence()

	header := tmproto.Header{Height: s.onomyApp.LastBlockHeight() + 1}
	txGen := cosmoscmd.MakeEncodingConfig(app.ModuleBasics).TxConfig

	_, _, err := signCheckDeliver(t, txGen, s.onomyApp.BaseApp, header, messages, "", []uint64{accountNum}, []uint64{accountSeq}, true, true, priv)
	require.NoError(t, err)
}

func setup(invCheckPeriod uint) (*app.OnomyApp, GenesisState) {
	db := dbm.NewMemDB()
	encCdc := cosmoscmd.MakeEncodingConfig(app.ModuleBasics)
	simApp := app.New(log.NewNopLogger(), db, nil, true, map[int64]bool{}, app.DefaultNodeHome, invCheckPeriod, encCdc, simapp.EmptyAppOptions{})
	return simApp.(*app.OnomyApp), NewDefaultGenesisState(encCdc.Marshaler)
}

func signCheckDeliver(
	t *testing.T, txCfg client.TxConfig, app *bam.BaseApp, header tmproto.Header, msgs []sdk.Msg,
	chainID string, accNums, accSeqs []uint64, expSimPass, expPass bool, priv ...cryptotypes.PrivKey,
) (sdk.GasInfo, *sdk.Result, error) {
	t.Helper()

	return simapp.SignCheckDeliver(t, txCfg, app, header, msgs, chainID, accNums, accSeqs, expSimPass, expPass, priv...)
}
