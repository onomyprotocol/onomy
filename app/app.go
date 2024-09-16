package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	autocli "cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/x/tx/signing"
	abci "github.com/cometbft/cometbft/abci/types"
	tmjson "github.com/cometbft/cometbft/libs/json"
	tmos "github.com/cometbft/cometbft/libs/os"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	sigtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/version"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	txmodule "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"

	// govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	// paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/onomyprotocol/onomy/app/keepers"
	"github.com/spf13/cast"

	"github.com/cosmos/cosmos-sdk/x/auth/ante"

	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	// authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	ibctesting "github.com/cosmos/ibc-go/v8/testing"
	ibctestingtypes "github.com/cosmos/ibc-go/v8/testing/types"
	providertypes "github.com/cosmos/interchain-security/v5/x/ccv/provider/types"

	// psm "github.com/onomyprotocol/reserve/x/psm/module"

	v1_0_1 "github.com/onomyprotocol/onomy/app/upgrades/v1.0.1"
	v1_0_3 "github.com/onomyprotocol/onomy/app/upgrades/v1.0.3"
	v1_0_3_4 "github.com/onomyprotocol/onomy/app/upgrades/v1.0.3.4"
	v1_0_3_5 "github.com/onomyprotocol/onomy/app/upgrades/v1.0.3.5"
	v1_1_1 "github.com/onomyprotocol/onomy/app/upgrades/v1.1.1"
	v1_1_2 "github.com/onomyprotocol/onomy/app/upgrades/v1.1.2"
	v1_1_4 "github.com/onomyprotocol/onomy/app/upgrades/v1.1.4"
	// "github.com/onomyprotocol/onomy/docs"
	// "github.com/onomyprotocol/onomy/x/dao"
	// daoclient "github.com/onomyprotocol/onomy/x/dao/client"
	// daokeeper "github.com/onomyprotocol/onomy/x/dao/keeper"
	// daotypes "github.com/onomyprotocol/onomy/x/dao/types"
)

const (
	// AccountAddressPrefix is cosmos-sdk accounts prefixes.
	AccountAddressPrefix = "onomy"
	// Name is the name of the onomy chain.
	AppName = "onomy"
)

// func getGovProposalHandlers() []govclient.ProposalHandler {
// 	var govProposalHandlers []govclient.ProposalHandler

// 	govProposalHandlers = append(govProposalHandlers,
// 		paramsclient.ProposalHandler,
// 		distrclient.ProposalHandler,
// 		upgradeclient.ProposalHandler,
// 		upgradeclient.CancelProposalHandler,
// 		ibcclientclient.UpdateClientProposalHandler,
// 		ibcclientclient.UpgradeProposalHandler,
// 		daoclient.FundTreasuryProposalHandler,
// 		daoclient.ExchangeWithTreasuryProposalProposalHandler,
// 		ibcproviderclient.ConsumerAdditionProposalHandler,
// 		ibcproviderclient.ConsumerRemovalProposalHandler,
// 		ibcproviderclient.EquivocationProposalHandler,
// 	)

// 	return govProposalHandlers
// }

var (
	// DefaultNodeHome default home directories for the application daemon.
	DefaultNodeHome string // nolint:gochecknoglobals // cosmos-sdk application style
)

var (
	_ runtime.AppI            = (*OnomyApp)(nil)
	_ servertypes.Application = (*OnomyApp)(nil)
	_ ibctesting.TestingApp   = (*OnomyApp)(nil)
)

func init() { // nolint: gochecknoinits // init funcs is are commonly used in cosmos
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, "."+AppName)

	// change default power reduction to 18 digits, since the onomy anom is 18 digits based.
	sdk.DefaultPowerReduction = math.NewIntWithDecimal(1, 18) // nolint: gomnd

	// Set prefixes
	accountPubKeyPrefix := AccountAddressPrefix + "pub"
	validatorAddressPrefix := AccountAddressPrefix + "valoper"
	validatorPubKeyPrefix := AccountAddressPrefix + "valoperpub"
	consNodeAddressPrefix := AccountAddressPrefix + "valcons"
	consNodePubKeyPrefix := AccountAddressPrefix + "valconspub"

	// Set and seal config
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(AccountAddressPrefix, accountPubKeyPrefix)
	config.SetBech32PrefixForValidator(validatorAddressPrefix, validatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(consNodeAddressPrefix, consNodePubKeyPrefix)
	config.Seal()
}

// OnomyApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type OnomyApp struct {
	*baseapp.BaseApp
	keepers.AppKeepers

	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// mm is the module manager
	mm           *module.Manager
	ModuleBasics module.BasicManager
	configurator module.Configurator
	// sm is the simulation manager
	sm *module.SimulationManager
}

// New returns a reference to an initialized blockchain app.
func NewOnomyApp( // nolint:funlen // app new cosmos func
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *OnomyApp {
	legacyAmino := codec.NewLegacyAmino()
	interfaceRegistry, err := types.NewInterfaceRegistryWithOptions(types.InterfaceRegistryOptions{
		ProtoFiles: proto.HybridResolver,
		SigningOptions: signing.Options{
			AddressCodec: address.Bech32Codec{
				Bech32Prefix: sdk.GetConfig().GetBech32AccountAddrPrefix(),
			},
			ValidatorAddressCodec: address.Bech32Codec{
				Bech32Prefix: sdk.GetConfig().GetBech32ValidatorAddrPrefix(),
			},
		},
	})
	if err != nil {
		panic(err)
	}

	appCodec := codec.NewProtoCodec(interfaceRegistry)
	txConfig := authtx.NewTxConfig(appCodec, authtx.DefaultSignModes)

	std.RegisterLegacyAminoCodec(legacyAmino)
	std.RegisterInterfaces(interfaceRegistry)

	// App Opts
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))
	invCheckPeriod := cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod))

	bApp := baseapp.NewBaseApp(
		AppName,
		logger,
		db,
		txConfig.TxDecoder(),
		baseAppOptions...,
	)

	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)
	bApp.SetTxEncoder(txConfig.TxEncoder())

	app := &OnomyApp{
		BaseApp:           bApp,
		legacyAmino:       legacyAmino,
		txConfig:          txConfig,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
	}

	moduleAccountAddresses := app.ModuleAccountAddrs()

	// Setup keepers
	app.AppKeepers = keepers.NewAppKeeper(
		appCodec,
		bApp,
		legacyAmino,
		maccPerms,
		moduleAccountAddresses,
		app.BlockedModuleAccountAddrs(moduleAccountAddresses),
		skipUpgradeHeights,
		homePath,
		invCheckPeriod,
		logger,
		appOpts,
	)

	// app.ProviderKeeper = ibcproviderkeeper.NewKeeper(
	// 	appCodec,
	// 	keys[providertypes.StoreKey],
	// 	app.GetSubspace(providertypes.ModuleName),
	// 	scopedIBCProviderKeeper,
	// 	app.IBCKeeper.ChannelKeeper,
	// 	&app.IBCKeeper.PortKeeper,
	// 	app.IBCKeeper.ConnectionKeeper,
	// 	app.IBCKeeper.ClientKeeper,
	// 	app.StakingKeeper,
	// 	app.SlashingKeeper,
	// 	app.AccountKeeper,
	// 	app.EvidenceKeeper,
	// 	app.DistrKeeper,
	// 	app.BankKeeper,
	// 	authtypes.FeeCollectorName,
	// )

	/****  Module Options ****/

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(appModules(app, appCodec, txConfig, skipGenesisInvariants)...)
	app.ModuleBasics = newBasicManagerFromManager(app)

	enabledSignModes := append([]sigtypes.SignMode(nil), authtx.DefaultSignModes...)
	enabledSignModes = append(enabledSignModes, sigtypes.SignMode_SIGN_MODE_TEXTUAL)

	txConfigOpts := authtx.ConfigOptions{
		EnabledSignModes:           enabledSignModes,
		TextualCoinMetadataQueryFn: txmodule.NewBankKeeperCoinMetadataQueryFn(app.BankKeeper),
	}
	txConfig, err = authtx.NewTxConfigWithOptions(
		appCodec,
		txConfigOpts,
	)
	if err != nil {
		panic(err)
	}
	app.txConfig = txConfig

	app.mm.SetOrderBeginBlockers(orderBeginBlockers()...)
	app.mm.SetOrderEndBlockers(orderEndBlockers()...)
	app.mm.SetOrderInitGenesis(orderInitBlockers()...)

	//
	app.mm.RegisterInvariants(app.CrisisKeeper)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	err = app.mm.RegisterServices(app.configurator)
	if err != nil {
		panic(err)
	}
	autocliv1.RegisterQueryServer(app.GRPCQueryRouter(), runtimeservices.NewAutoCLIQueryService(app.mm.Modules))

	reflectionSvc, err := runtimeservices.NewReflectionService()
	if err != nil {
		panic(err)
	}
	reflectionv1.RegisterReflectionServiceServer(app.GRPCQueryRouter(), reflectionSvc)

	// add test gRPC service for testing gRPC queries in isolation
	testdata.RegisterQueryServer(app.GRPCQueryRouter(), testdata.QueryImpl{})
	// create the simulation manager and define the order of the modules for deterministic simulations
	//
	// NOTE: this is not required apps that don't use the simulator for fuzz testing
	// transactions
	app.sm = module.NewSimulationManager(simulationModules(app, appCodec, skipGenesisInvariants)...)
	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(app.GetKVStoreKey())
	app.MountTransientStores(app.GetTransientStoreKey())
	app.MountMemoryStores(app.GetMemoryStoreKey())

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetPreBlocker(app.PreBlocker)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	anteHandler, err := ante.NewAnteHandler(
		ante.HandlerOptions{
			AccountKeeper:   app.AccountKeeper,
			BankKeeper:      app.BankKeeper,
			SignModeHandler: txConfig.SignModeHandler(),
			FeegrantKeeper:  app.FeeGrantKeeper,
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
	)
	if err != nil {
		panic(err)
	}

	// and storeloader
	app.setupUpgradeHandlers()
	// app.setupUpgradeStoreLoaders()

	app.SetAnteHandler(anteHandler)

	// At startup, after all modules have been registered, check that all prot
	// annotations are correct.
	protoFiles, err := proto.MergedRegistry()
	if err != nil {
		panic(err)
	}
	err = msgservice.ValidateProtoAnnotations(protoFiles)
	if err != nil {
		// Once we switch to using protoreflect-based antehandlers, we might
		// want to panic here instead of logging a warning.
		fmt.Fprintln(os.Stderr, err.Error())
	}

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(fmt.Sprintf("failed to load latest version: %s", err))
		}
	}
	return app
}

// Name returns the name of the OnomyApp.
func (app *OnomyApp) Name() string { return app.BaseApp.Name() }

// GetBaseApp returns the base app of the application.
func (app OnomyApp) GetBaseApp() *baseapp.BaseApp { return app.BaseApp }

func (app *OnomyApp) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	return app.mm.BeginBlock(ctx)
}

func (app *OnomyApp) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	return app.mm.EndBlock(ctx)
}

func (app *OnomyApp) InitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	if err := app.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap()); err != nil {
		panic(err)
	}

	response, err := app.mm.InitGenesis(ctx, app.appCodec, genesisState)
	if err != nil {
		panic(err)
	}

	return response, nil
}

func (app *OnomyApp) PreBlocker(ctx sdk.Context, _ *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	return app.mm.PreBlock(ctx)
}

// LoadHeight loads a particular height.
func (app *OnomyApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *OnomyApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlockedAddrs returns all the app's module account addresses that are not
// allowed to receive external tokens.
func (app *OnomyApp) BlockedAddrs() map[string]bool {
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = !allowedReceivingModAcc[acc]
	}

	// For ICS multiden fix
	delete(blockedAddrs, authtypes.NewModuleAddress(providertypes.ConsumerRewardsPool).String())

	return blockedAddrs
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *OnomyApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// AppCodec returns an app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *OnomyApp) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns an InterfaceRegistry.
func (app *OnomyApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
func (app *OnomyApp) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.AppKeepers.GetKey(storeKey)
}

// GetMemKey returns the MemoryStoreKey for the provided store key.
func (app *OnomyApp) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.AppKeepers.GetMemKey(storeKey)
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *OnomyApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *OnomyApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	cmtservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	app.ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register nodeservice grpc-gateway routes.
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if err := server.RegisterSwaggerAPI(apiSvr.ClientCtx, apiSvr.Router, apiConfig.Swagger); err != nil {
		panic(err)
	}
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *OnomyApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *OnomyApp) RegisterTendermintService(clientCtx client.Context) {
	cmtservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		app.Query,
	)
}

// SetOrderEndBlockers sets the order of set end-blocker calls.
func (app *OnomyApp) SetOrderEndBlockers(moduleNames ...string) {
	app.mm.SetOrderEndBlockers(moduleNames...)
}

// SimulationManager implements the SimulationApp interface.
func (app *OnomyApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

func (app *OnomyApp) setupUpgradeHandlers() {
	app.UpgradeKeeper.SetUpgradeHandler(v1_0_1.Name, v1_0_1.UpgradeHandler)
	app.UpgradeKeeper.SetUpgradeHandler(v1_0_3.Name, v1_0_3.UpgradeHandler)
	app.UpgradeKeeper.SetUpgradeHandler(v1_0_3_4.Name, v1_0_3_4.UpgradeHandler)
	app.UpgradeKeeper.SetUpgradeHandler(v1_0_3_5.Name, v1_0_3_5.UpgradeHandler)
	// we need to have the reference to `app` which is why we need this `func` here
	app.UpgradeKeeper.SetUpgradeHandler(
		v1_1_1.Name,
		func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			for moduleName, eachModule := range app.mm.Modules {
				m := eachModule.(module.GenesisOnlyAppModule)
				fromVM[moduleName] = m.ConsensusVersion()
			}

			// This is critical for the chain upgrade to work
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			app.ProviderKeeper.InitGenesis(sdkCtx, providertypes.DefaultGenesisState())

			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)
	app.UpgradeKeeper.SetUpgradeHandler(v1_1_2.Name, v1_1_2.UpgradeHandler)
	app.UpgradeKeeper.SetUpgradeHandler(v1_1_4.Name, v1_1_4.UpgradeHandler)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Errorf("failed to read upgrade info from disk: %w", err))
	}

	// configure store loader that checks if version == upgradeHeight and applies store upgrades
	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	var storeUpgrades *storetypes.StoreUpgrades

	switch upgradeInfo.Name {
	case v1_1_1.Name:
		storeUpgrades = &storetypes.StoreUpgrades{
			Added: []string{providertypes.ModuleName, providertypes.StoreKey},
		}
	case v1_1_4.Name:
		storeUpgrades = &storetypes.StoreUpgrades{
			Added: []string{authtypes.ModuleName, authzkeeper.StoreKey},
		}
	default:
		// no store upgrades
	}

	if storeUpgrades != nil {
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, storeUpgrades))
	}
}

func (app *OnomyApp) BlockedModuleAccountAddrs(modAccAddrs map[string]bool) map[string]bool {
	// remove module accounts that are ALLOWED to received funds
	delete(modAccAddrs, authtypes.NewModuleAddress(govtypes.ModuleName).String())
	return modAccAddrs
}

// func (app *OnomyApp) setupUpgradeStoreLoaders() {
// 	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
// 	if err != nil {
// 		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
// 	}

// 	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
// 		return
// 	}
// }

func (app *OnomyApp) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter(), cfg)
}

// GetIBCKeeper returns the IBC keeper.
func (app *OnomyApp) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper
}

// GetScopedIBCKeeper implements the TestingApp interface.
func (app *OnomyApp) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper { //nolint:nolintlint
	return app.ScopedIBCKeeper
}

// GetStakingKeeper implements the TestingApp interface. Needed for ICS.
func (app *OnomyApp) GetStakingKeeper() ibctestingtypes.StakingKeeper { //nolint:nolintlint
	return app.StakingKeeper
}

// GetTxConfig implements the TestingApp interface.
func (app *OnomyApp) GetTxConfig() client.TxConfig {
	return app.txConfig
}

// AutoCliOpts returns the autocli options for the app.
func (app *OnomyApp) AutoCliOpts() autocli.AppOptions {
	modules := make(map[string]appmodule.AppModule, 0)
	for _, m := range app.mm.Modules {
		if moduleWithName, ok := m.(module.HasName); ok {
			moduleName := moduleWithName.Name()
			if appModule, ok := moduleWithName.(appmodule.AppModule); ok {
				modules[moduleName] = appModule
			}
		}
	}

	return autocli.AppOptions{
		Modules:               modules,
		AddressCodec:          authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		ValidatorAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		ConsensusAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	}
}
