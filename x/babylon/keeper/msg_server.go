package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	k Keeper
}

// NewMsgServer constructor
func NewMsgServer(k Keeper) types.MsgServer {
	return &msgServer{k: k}
}

var _ types.MsgServer = msgServer{}

func (ms msgServer) SetBSNContracts(goCtx context.Context, req *types.MsgSetBSNContracts) (*types.MsgSetBSNContractsResponse, error) {
	if err := req.ValidateBasic(); err != nil {
		return nil, err
	}

	if authority := ms.k.GetAuthority(); authority != req.Authority {
		return nil, govtypes.ErrInvalidSigner.Wrapf("invalid authority; expected %s, got %s", authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := ms.k.SetBSNContracts(ctx, req.Contracts); err != nil {
		return nil, err
	}

	return &types.MsgSetBSNContractsResponse{}, nil
}

// UpdateParams updates the params.
func (ms msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if ms.k.authority != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", ms.k.authority, req.Authority)
	}
	if err := req.Params.ValidateBasic(); err != nil {
		return nil, govtypes.ErrInvalidProposalMsg.Wrapf("invalid parameter: %v", err)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := ms.k.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
