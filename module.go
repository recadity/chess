package chess

import (
    "github.com/cosmos/cosmos-sdk/codec"
    "github.com/cosmos/cosmos-sdk/types/module"
    "github.com/recadity/chess/client/cli"
    "github.com/spf13/cobra"
)

type AppModuleBasic struct{}

func (AppModuleBasic) Name() string {
    return ModuleName
}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
    RegisterCodec(cdc)
}

func (AppModuleBasic) DefaultGenesis() json.RawMessage {
    return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
    var data GenesisState
    if err := ModuleCdc.UnmarshalJSON(bz, &data); err != nil {
        return err
    }
    return ValidateGenesis(data)
}

func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
    return cli.GetTxCmd(StoreKey, cdc)
}

func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
    return cli.GetQueryCmd(StoreKey, cdc)
}

// AppModule implements an application module for the chess module.
type AppModule struct {
    AppModuleBasic
    keeper Keeper
}

func NewAppModule(k Keeper) AppModule {
    return AppModule{
        AppModuleBasic: AppModuleBasic{},
        keeper:         k,
    }
}

// RegisterInvariants does nothing, there are no invariants to enforce
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

// Route returns the message routing key for the chess module.
func (am AppModule) Route() string {
    return RouterKey
}

// NewHandler returns an sdk.Handler for the chess module.
func (am AppModule) NewHandler() sdk.Handler {
    return NewHandler(am.keeper)
}

// QuerierRoute returns the chess module's querier route name.
func (am AppModule) QuerierRoute() string {
    return QuerierRoute
}

// NewQuerierHandler returns the chess module sdk.Querier.
func (am AppModule) NewQuerierHandler() sdk.Querier {
    return NewQuerier(am.keeper)
}

// InitGenesis initializes the genesis state for the chess module.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc *codec.Codec, data json.RawMessage) []abci.ValidatorUpdate {
    var genesisState GenesisState
    cdc.MustUnmarshalJSON(data, &genesisState)
    return InitGenesis(ctx, am.keeper, genesisState)
}

// ExportGenesis exports the genesis state for the chess module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc *codec.Codec) json.RawMessage {
    gs := ExportGenesis(ctx, am.keeper)
    return cdc.MustMarshalJSON(gs)
}

