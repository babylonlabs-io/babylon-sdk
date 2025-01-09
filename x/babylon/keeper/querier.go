package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
)

var _ types.QueryServer = &querier{}

type querier struct {
	k *Keeper
}

// NewQuerier constructor
func NewQuerier(k *Keeper) *querier {
	return &querier{k: k}
}

// Params implements the gRPC service handler for querying the babylon parameters.
func (q querier) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params := q.k.GetParams(sdk.UnwrapSDKContext(ctx))
	return &types.QueryParamsResponse{Params: params}, nil
}
