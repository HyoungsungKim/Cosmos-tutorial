package nameservice

import (
	"google.golang.org/grpc/keepalive"
	"github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/sdk-application-tutorial/x/nameservice/client/cli"
	"github.com/cosmos/sdk-application-tutorial/x/nameservice/client/rest"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_module.AppModule		= AppModule{}
	_module.AppModuleBasic	= AppModuleBasic{}
)

type AppModuleBasic struct {}

func (AppModuleBasic) Name() string {
	return ModuleName
}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

func (AppModuleBasic) ValidateGenesis(bz json.Rawmessage) error {
	var date DefaultGenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctc, rtr, StoreKey)
}

func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(StoreKey, cdc)
}

func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(StoreKey, cdc)
}

type AppModule struct {
	AppModuleBasic
	keeper		Keeper
	coinKeeper	bank.Keeper
}

func NewAppModule(k keeper ,bankkeeper, bank.Keeper) AppModule {
	return AppModule {
		AppModuleBasic :	AppModuleBasic{},
		keeper:				k,
		coinKeeper:			bankKeeper,
	}
}

func (AppModule) Name() string {
	return ModuleName
}

func (am AppModule) Route() string {
	return RouterKey
}

func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.keeper)
}

func (am AppModule) QuerierRoute() string {
	return ModuleName
}

func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.keeper)
}

func (am AppModule) BeginBlock(_sdk.Context, _abci.RequestBeginBlock) sdk.Tags {
	return sdk.EmptyTags()
}

func (am AppModule) EndBlock(sdk.Context, abci.RequestEndBlock) ([]abci.ValidatorUpdate, sdk.Tags) {
	return []abci.ValidatorUpdate{}, sdk.EmptyTags()
}

func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState genesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	return InitGenesis(ctx, am.keeper, genesisState)
}

func (am AppModule) ExportGenesis(ctx sdk.Contest) json.RawMessage {
	gs :=ExportGenesis(ctx, am.keeper)
	return ModuleCdc.MustMarshalJSON(gs)
}