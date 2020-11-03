package app

import (
	"encoding/json"
	"github.com/bentaro/nftchain/x/wasm"
	"io"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/bentaro/nftchain/x/nftchain"
	nftchainkeeper "github.com/bentaro/nftchain/x/nftchain/keeper"
	nftchaintypes "github.com/bentaro/nftchain/x/nftchain/types"
    // this line is used by starport scaffolding # 1
    "path/filepath"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/spf13/viper"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	//nft module
	"github.com/cosmos/modules/incubator/nft"
)

const appName = "nftchain"

var (
	DefaultCLIHome = os.ExpandEnv("$HOME/.nftchaincli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.nftchaind")
	ModuleBasics = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		params.AppModuleBasic{},
		supply.AppModuleBasic{},
		nftchain.AppModuleBasic{},
    	// this line is used by starport scaffolding # 2
		distr.AppModuleBasic{},
		wasm.AppModuleBasic{},
		//add nft module basic
		nft.AppModuleBasic{},
	)

	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		// this line is used by starport scaffolding # 2.1
		distr.ModuleName: nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
	}
)

func MakeCodec() *codec.Codec {
	var cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	return cdc.Seal()
}

type NewApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	keys  map[string]*sdk.KVStoreKey
	tKeys map[string]*sdk.TransientStoreKey

	subspaces map[string]params.Subspace

	accountKeeper  auth.AccountKeeper
	bankKeeper     bank.Keeper
	stakingKeeper  staking.Keeper
	supplyKeeper   supply.Keeper
	paramsKeeper   params.Keeper
	nftchainKeeper nftchainkeeper.Keeper

	//add nft keeper
	nftKeeper      nft.Keeper

  // this line is used by starport scaffolding # 3
  	distrKeeper    distr.Keeper
	wasmKeeper    wasm.Keeper
	mm *module.Manager

	sm *module.SimulationManager
}

var _ simapp.App = (*NewApp)(nil)

func NewInitApp(
	logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp),
) *NewApp {
	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	keys := sdk.NewKVStoreKeys(
    bam.MainStoreKey,
    auth.StoreKey,
    staking.StoreKey,
		supply.StoreKey,
    params.StoreKey,
    //add nft store key
    nft.StoreKey,
    nftchaintypes.StoreKey,
    // this line is used by starport scaffolding # 5
    distr.StoreKey,
    wasm.StoreKey,
  )

	tKeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	var app = &NewApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tKeys:          tKeys,
		subspaces:      make(map[string]params.Subspace),
	}

	app.paramsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tKeys[params.TStoreKey])
	app.subspaces[auth.ModuleName] = app.paramsKeeper.Subspace(auth.DefaultParamspace)
	app.subspaces[bank.ModuleName] = app.paramsKeeper.Subspace(bank.DefaultParamspace)
	app.subspaces[staking.ModuleName] = app.paramsKeeper.Subspace(staking.DefaultParamspace)
	// this line is used by starport scaffolding # 5.1
	app.subspaces[distr.ModuleName] = app.paramsKeeper.Subspace(distr.DefaultParamspace)

	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		keys[auth.StoreKey],
		app.subspaces[auth.ModuleName],
		auth.ProtoBaseAccount,
	)

	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		app.subspaces[bank.ModuleName],
		app.ModuleAccountAddrs(),
	)

	app.supplyKeeper = supply.NewKeeper(
		app.cdc,
		keys[supply.StoreKey],
		app.accountKeeper,
		app.bankKeeper,
		maccPerms,
	)

	stakingKeeper := staking.NewKeeper(
		app.cdc,
		keys[staking.StoreKey],
		app.supplyKeeper,
		app.subspaces[staking.ModuleName],
	)

	// this line is used by starport scaffolding # 5.2
	app.distrKeeper = distr.NewKeeper(
		app.cdc, keys[distr.StoreKey], app.subspaces[distr.ModuleName], &stakingKeeper,
		app.supplyKeeper, auth.FeeCollectorName, app.ModuleAccountAddrs(),
			)

	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(
			// this line is used by starport scaffolding # 5.3
		app.distrKeeper.Hooks(),
		),
	)

	//instanciate nft keeper
	app.nftKeeper = nft.NewKeeper(
		app.cdc,
		keys[nft.StoreKey],
	)

	app.nftchainKeeper = nftchainkeeper.NewKeeper(
		app.bankKeeper,
		app.cdc,
		keys[nftchaintypes.StoreKey],
	)

  // this line is used by starport scaffolding # 4
type WasmWrapper struct { Wasm wasm.Config `mapstructure:"wasm"`}
		var wasmRouter = bApp.Router()
		homeDir := viper.GetString(cli.HomeFlag)
		wasmDir := filepath.Join(homeDir, "wasm")

		wasmWrap := WasmWrapper{Wasm: wasm.DefaultWasmConfig()}
		err := viper.Unmarshal(&wasmWrap)
		if err != nil {
			panic("error while reading wasm config: " + err.Error())
		}
		wasmConfig := wasmWrap.Wasm
		supportedFeatures := "staking"
		app.subspaces[wasm.ModuleName] = app.paramsKeeper.Subspace(wasm.DefaultParamspace)
		app.wasmKeeper = wasm.NewKeeper(app.cdc, keys[wasm.StoreKey], app.subspaces[wasm.ModuleName], app.accountKeeper, app.bankKeeper, app.stakingKeeper, app.distrKeeper, app.nftKeeper, wasmRouter, wasmDir, wasmConfig, supportedFeatures, nil, nil)

	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		nftchain.NewAppModule(app.nftchainKeeper, app.bankKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
    	// this line is used by starport scaffolding # 6
		distr.NewAppModule(app.distrKeeper, app.accountKeeper, app.supplyKeeper, app.stakingKeeper),
		wasm.NewAppModule(app.wasmKeeper),
		//add nft new app module
		nft.NewAppModule(app.nftKeeper, app.accountKeeper),
	)

	app.mm.SetOrderEndBlockers(
		staking.ModuleName,
		// this line is used by starport scaffolding # 6.1
		distr.ModuleName,
	)

	app.mm.SetOrderInitGenesis(
		// this line is used by starport scaffolding # 6.2
		distr.ModuleName,
		staking.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		nftchaintypes.ModuleName,
		supply.ModuleName,
		genutil.ModuleName,
        // this line is used by starport scaffolding # 7
		wasm.ModuleName,
		//nft module
		nft.ModuleName,
	)

	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	app.SetAnteHandler(
		auth.NewAnteHandler(
			app.accountKeeper,
			app.supplyKeeper,
			auth.DefaultSigVerificationGasConsumer,
		),
	)

	app.MountKVStores(keys)
	app.MountTransientStores(tKeys)

	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
		if err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

type GenesisState map[string]json.RawMessage

func NewDefaultGenesisState() GenesisState {
	return ModuleBasics.DefaultGenesis()
}

func (app *NewApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState

	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	return app.mm.InitGenesis(ctx, genesisState)
}

func (app *NewApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (app *NewApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *NewApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

func (app *NewApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func (app *NewApp) Codec() *codec.Codec {
	return app.cdc
}

func (app *NewApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

func GetMaccPerms() map[string][]string {
	modAccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		modAccPerms[k] = v
	}
	return modAccPerms
}
