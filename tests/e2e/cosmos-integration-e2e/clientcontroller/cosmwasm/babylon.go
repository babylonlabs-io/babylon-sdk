package cosmwasm

import (
	"context"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
	"github.com/cosmos/cosmos-sdk/client"
)

func (cc *CosmwasmConsumerController) MustQueryBabylonParams() *types.Params {
	ctx := context.Background()

	clientCtx := client.Context{Client: cc.cwClient.RPCClient}
	queryClient := types.NewQueryClient(clientCtx)

	resp, err := queryClient.Params(ctx, &types.QueryParamsRequest{})
	if err != nil {
		panic(err)
	}

	return &resp.Params
}
