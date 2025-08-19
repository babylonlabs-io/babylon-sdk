package keeper_test

import (
	"testing"

	"github.com/cometbft/cometbft/libs/rand"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
)

func TestGRPCQuery_Params(t *testing.T) {
	keepers := NewTestKeepers(t)
	ctx := keepers.Ctx
	k := keepers.BabylonKeeper

	params := types.DefaultParams()
	require.NoError(t, k.SetParams(ctx, params))

	resp, err := k.Params(ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, params, resp.Params)
}

func TestGRPCQuery_Params_InvalidRequest(t *testing.T) {
	keepers := NewTestKeepers(t)
	ctx := keepers.Ctx
	k := keepers.BabylonKeeper

	_, err := k.Params(ctx, nil)
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, st.Code())
}

func TestGRPCQuery_BSNContracts(t *testing.T) {
	keepers := NewTestKeepers(t)
	ctx := keepers.Ctx
	k := keepers.BabylonKeeper

	addr1 := sdk.AccAddress(rand.Bytes(20)).String()
	addr2 := sdk.AccAddress(rand.Bytes(20)).String()
	addr3 := sdk.AccAddress(rand.Bytes(20)).String()
	addr4 := sdk.AccAddress(rand.Bytes(20)).String()

	contracts := &types.BSNContracts{
		BabylonContract:        addr1,
		BtcLightClientContract: addr2,
		BtcStakingContract:     addr3,
		BtcFinalityContract:    addr4,
	}
	require.NoError(t, k.SetBSNContracts(ctx, contracts))

	resp, err := k.BSNContracts(ctx, &types.QueryBSNContractsRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp.BsnContracts)
	require.True(t, contracts.Equal(resp.BsnContracts))
}

func TestGRPCQuery_BSNContracts_NotSet(t *testing.T) {
	keepers := NewTestKeepers(t)
	ctx := keepers.Ctx
	k := keepers.BabylonKeeper

	resp, err := k.BSNContracts(ctx, &types.QueryBSNContractsRequest{})
	require.NoError(t, err)
	require.Nil(t, resp.BsnContracts)
}

func TestGRPCQuery_BSNContracts_InvalidRequest(t *testing.T) {
	keepers := NewTestKeepers(t)
	ctx := keepers.Ctx
	k := keepers.BabylonKeeper

	_, err := k.BSNContracts(ctx, nil)
	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, st.Code())
}
