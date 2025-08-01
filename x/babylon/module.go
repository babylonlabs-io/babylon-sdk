package babylon

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/appmodule"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/client/cli"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/keeper"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
)

// ConsensusVersion defines the module's consensus version.
const ConsensusVersion = 1

var (
	_ appmodule.AppModule       = AppModule{}
	_ appmodule.HasBeginBlocker = AppModule{}
	_ module.HasABCIEndBlock    = AppModule{}
	_ module.AppModuleBasic     = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the babylon module.
type AppModuleBasic struct{}

func (b AppModuleBasic) RegisterLegacyAminoCodec(amino *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(amino)
}

func (b AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, serveMux *runtime.ServeMux) {
	err := types.RegisterQueryHandlerClient(context.Background(), serveMux, types.NewQueryClient(clientCtx))
	if err != nil {
		panic(err)
	}
}

// Name returns the babylon module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// DefaultGenesis returns default genesis state as raw bytes for the babylon
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the babylon module.
func (b AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}

	return types.ValidateGenesis(&data)
}

// GetTxCmd returns the root tx command for the babylon module.
func (b AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// GetQueryCmd returns no root query command for the babylon module.
func (b AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterInterfaces implements InterfaceModule
func (b AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// AppModule implements an application module interface
type AppModule struct {
	AppModuleBasic
	cdc codec.Codec
	k   *keeper.Keeper
}

// NewAppModule constructor with defaults
func NewAppModule(cdc codec.Codec, k *keeper.Keeper) *AppModule {
	return NewAppModuleX(cdc, k)
}

// NewAppModuleX extended constructor
func NewAppModuleX(cdc codec.Codec, k *keeper.Keeper) *AppModule {
	return &AppModule{cdc: cdc, k: k}
}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

// Name returns the module's name.
func (AppModule) Name() string { return types.ModuleName }

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServer(am.k))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.NewQuerier(am.cdc, am.k))
}

// RegisterInvariants registers the module's invariants.
func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
}

// InitGenesis performs genesis initialization for the babylon module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, bz json.RawMessage) []abci.ValidatorUpdate {
	var data types.GenesisState
	cdc.MustUnmarshalJSON(bz, &data)
	am.k.InitGenesis(ctx, data)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the babylon
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(am.k.ExportGenesis(ctx))
}

// QuerierRoute returns the bank module's querier route name.
func (AppModule) QuerierRoute() string {
	return types.RouterKey
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 {
	return ConsensusVersion
}

// BeginBlock executed before every block
func (am AppModule) BeginBlock(ctx context.Context) error {
	return am.k.BeginBlocker(ctx)
}

// EndBlock executed after every block. It returns no validator updates.
func (am AppModule) EndBlock(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	return am.k.EndBlocker(ctx)
}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() { // marker
}
