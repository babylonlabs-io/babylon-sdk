package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	if err := k.SetParams(ctx, data.Params); err != nil {
		panic(fmt.Errorf("failed to set params in genesis: %w, params: %v", err, data.BsnContracts))
	}
	// Set BSN contracts if provided
	if data.BsnContracts != nil && data.BsnContracts.IsSet() {
		if err := k.SetBSNContracts(ctx, data.BsnContracts); err != nil {
			panic(err)
		}
	}
}

func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params := k.GetParams(ctx)
	contracts := k.GetBSNContracts(ctx)
	return types.NewGenesisState(params, contracts)
}
