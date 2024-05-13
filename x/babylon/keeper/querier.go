package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/babylonchain/babylon-sdk/x/babylon/types"
)

var _ types.QueryServer = &querier{}

type querier struct {
	cdc codec.Codec
	k   *Keeper
}

// NewQuerier constructor
func NewQuerier(cdc codec.Codec, k *Keeper) *querier {
	return &querier{cdc: cdc, k: k}
}

// Params implements the gRPC service handler for querying the babylon parameters.
func (q querier) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params := q.k.GetParams(sdk.UnwrapSDKContext(ctx))
	return &types.QueryParamsResponse{Params: params}, nil
}
