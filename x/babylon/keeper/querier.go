package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
)

var _ types.QueryServer = Keeper{}

// Params implements the gRPC service handler for querying the babylon parameters.
func (k Keeper) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	params := k.GetParams(sdk.UnwrapSDKContext(ctx))
	return &types.QueryParamsResponse{Params: params}, nil
}

// BSNContracts implements the gRPC service handler for querying the babylon contract addresses.
func (k Keeper) BSNContracts(ctx context.Context, req *types.QueryBSNContractsRequest) (*types.QueryBSNContractsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	contracts := k.GetBSNContracts(sdkCtx)
	return &types.QueryBSNContractsResponse{
		BsnContracts: contracts,
	}, nil
}
