package cosmwasm

import (
	"context"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
	"github.com/cosmos/cosmos-sdk/client"
)

// MustQueryBabylonContracts queries the Babylon module for all contract addresses and panics on error.
type BabylonContracts struct {
	BabylonContract        string
	BtcLightClientContract string
	BtcStakingContract     string
	BtcFinalityContract    string
}

func (cc *CosmwasmConsumerController) MustQueryBabylonContracts() *BabylonContracts {
	ctx := context.Background()

	clientCtx := client.Context{Client: cc.cwClient.RPCClient}
	queryClient := types.NewQueryClient(clientCtx)

	resp, err := queryClient.BSNContracts(ctx, &types.QueryBSNContractsRequest{})
	if err != nil {
		panic(err)
	}

	return &BabylonContracts{
		BabylonContract:        resp.BsnContracts.BabylonContract,
		BtcLightClientContract: resp.BsnContracts.BtcLightClientContract,
		BtcStakingContract:     resp.BsnContracts.BtcStakingContract,
		BtcFinalityContract:    resp.BsnContracts.BtcFinalityContract,
	}
}
