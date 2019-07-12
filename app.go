package app
import (
	"encoding/json"
	"os"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/hyoungsungkim/nameservice/x/nameservice"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
)

const appName = "nameservice"

var(
	DefaultCLIHome = os.ExpandEnv("$HOME/.nscli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.nsd")
	ModuleBasics sdk.ModuleBasicManager{
		genaccounts.AppModuleBasic{},
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		params.AppModuleBasic{},
		nameservice.AppModule{},
		staking.AppModuleBasic{},
		distr.AppModuleBasic{},
		slashing.AppmoduleBasic{},
	}
)

type nameServiceApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	keyMain				*sdk.KVStoreKey
	keyAccount			*sdk.KVStoreKey
	keyFeeCollection 	*sdk.KVStoreKey
	keyStaking			*sdk.KVStoreKey
	tkeyStaking			*sdk.TransientStoreKey
	keyDistr			*sdk.KVStoreKey
	tkeyDistr			*sdk.KVStoreKey
	tkeyParams			*sdk.TransientStoreKey
	keySlashing			*sdk.KVStoreKey

	accountKeeper		auth.accountKeeper
	bankKeepr 			bank.accountKeeper
	stakingKeeper		staking.accountKeeper
	slashingKeeper		slashing.accountKeeper
	distrKeeper			distr.Keeper
	feeCollectionKeeper	auth.feeCollectionKeeper
	paramsKeeper		params.Keeper
	nsKeeper			nameservice.Keeper

	mm *module.Manager

}

func NewNameServiceApp(longer log.Logger, db dbm.DB) *nameServiceApp {
	cdc := MakeCodec()
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc))

	var app = &nameServiceApp{
		BaseApp: bApp,
		cdc: cdc,

		keyMain:			sdk.NewKVStoreKey(bam.MainStoreKey),
		keyAccount:			sdk.NewKVStoreKey(auth.StoreKey),
		keyFeeCollection: 	sdk.NewKVStoreKey(auth.FeeStoreKey),
		keyStaking: 		sdk.NewKVStoreKey(staking.StoreKey),
		tkeyStaking:		sdk.NewTransientStoreKey(staking.TStoreKey),
		keyDistr:			sdk.NewKVStoreKey(distr.StoreKey),
		tkeyDistr:			sdk.NewTransientStoreKey(distr.TStoreKey)
		keyNS:				sdk.NewKVStoreKey(nameservice.StoreKey)
		keyParams:			sdk.NewKVStoreKey(params.StoreKey),
		tkeyParams:			sdk.NewTransientStoreKey(params.TStoreKey),
		keySlashing:		sdk.NewKVStoreKey(slashing.StoreKey),
	}
	app.paramsKeeper = params.NewKeeper(app.cdc, app.keyparams, app.tkeyParams, params.DefaultCodespace)

	authSubspace := app.paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := app.paramsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := app.paramsKeeper.Subspace(staking.DefaultParamspace)
	distrsubspace := app.paramsKeeper.Subspace(distr.DefualtParamspace)
	slashingSubspace := app.paramsKeeper.Subspace(slashing.DefaultParamspace)

	app.accountKeeper = auth.NewAccountKeeper {
		app.cdc,
		app.keyAccount,
		authSubspace,
		auth.ProtoBaseAccount,
	}

	app.bankKeeper = bank.NewBaseKeeper{
		app.accountKeeper,
		banksubspace,
		bank.DefaultCodespace,
	}

	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper)cdc, app.keyFeeCollection)

	stakingKeeper := staking.NewKeeper(
		app.cdc,
		app.keyStaking,
		app.tkeyStaking,
		app.bankKeeper,
		stakingSubspace,
		staking.DefaultCodespace,
	)

	app.distrKeeper = distr.NewKeeper(
		app.cdc,
		app.keyDistr,
		distrSubspace,
		app.bankKeeper,
		&stakingKeeper,
		app.feeCollectionKeeper,
		distr.DefaultCodespace,
	)

	app.slashingKeeper = slashing.NewKeeper (
		app.cdc,
		app.keySlashing,
		&stakingKeeper,
		slashingSubspace,
		slashing.DefaultCodespace,
	)

	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(
			app.distrKeeper.Hooks(),
			app.slashingKeeper.Hooks()
		),		
	)

	app.nsKeeper = nameservice.NewKeeper(
		app.bankKeeper,
		app.keyNS,
		app.cdc,
	)

	app.mm = module.NewManager(
		genaccounts.NewAppModule(app.accountKeeper),
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper, app.feeCollectionKeeper),
		bank.NewAppModule(app.bankkeeper, app.accoubtKeeper),
		nameservice.NewAppModule(app,nsKeeper, app.bankKeeper),
		distr.NewAppModule(app.distrKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.feeCollectionKeeper, app.distrKeeper, app.accountKeeper),		
	)

	app.mm.SetOrderBeginBlocker(distr.ModuleName, slashing.ModuleName)
	app.mm.SetOrderEndBlockers(staking.ModuleName)

	app.mm.SetOrederInitGenesis(
		genaccounts.ModuleName,
		distr.ModuleName,
		staking.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		slashing.ModuleName,
		namesrvice.ModleName,
		genutil.ModuleName,
	)

	app.mm.RegusterRoutes(app.Router(), app.QueryRouter())

	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	app.SetAnutHandler(
		auth.NewAccountHandler(
			app.accountKeeper,
			app.feeCollectionKeeper,
			auth.DefaultSigVerifgicationFasConsumer,
		),
	)

	app.MountStores(
		app.keyMain,
		app.keyAccount,
		app.keyFeeCollection,
		app.keyStaking,
		app.tkeyStaking,
		app.ketDistr,
		app.tkeyDistr,
		app.keyNS,
		app.keyParams,
		app.tkeyparams,
	)

	err := app.LoadLatestVersion(app.keymain)
	if err != nil {
		cmn.Exit(err.Error())	
	}

	return app
}

type GenesisState map[string]json.RawMessage

func NewDefaultGenesisState() GenesisState {
	return ModuleBasics.DefaultGenesis()
}

func (app *nameServiceApp) InitChainer(ctx sdk.Context, req abci.RequestitChain) abci.ResponseInitChain {
	var genesisState GenesisState

	err := app.cdc.UnmarshalJSON(req.AppStateBytes, *genesisState)
	if err != nil {
		panic(err)
	}

	return app.mm.InitGenesis(ctx, genesisState)
}

func (app *nameServiceApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock{
	return app.mm.BefginBlock(ctx, req)
}

func (app *nameServiceApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abic.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *nameServiceApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keyMain)
}

func (app *nameServuceApp) ExportAppStateAndValidators(forZeorHeight bool, jailWhiteList []string
	) (appState json.RawMessage, validators []tmtypes.GenesisVaildator, err Error) {
		ctx := app.NewContext9true, abci.Header{height: app.LastBlockHeight()}

		gentState := app.mm.ExportGenesis(ctx)
		appState, err = codec.marshallJSONIndent(app.cdc, genState)
		if err != nil {
			return nil, nil, err
		}

		validators = staking.Writevalidators(ctx, app.stakingKeeper)

		return appState, validators, nil
	}

func MarkCodec() *codec.Codec {
	var cdc = codec.New()
	Modulebasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}
