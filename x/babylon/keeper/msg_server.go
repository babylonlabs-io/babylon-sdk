package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	k *Keeper
}

// NewMsgServer constructor
func NewMsgServer(k *Keeper) *msgServer {
	return &msgServer{k: k}
}

func (ms msgServer) InstantiateBabylonContracts(goCtx context.Context, req *types.MsgInstantiateBabylonContracts) (*types.MsgInstantiateBabylonContractsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := ms.k.GetParams(ctx)

	// only the authority can override the instantiated contracts
	if params.IsContractInstantiated() && req.Signer != ms.k.authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "only authority can override instantiated contracts; expected %s, got %s", ms.k.authority, req.Signer)
	}

	// construct the init message
	initMsg, err := types.NewInitMsg(
		req.Network,
		req.BabylonTag,
		req.BtcConfirmationDepth,
		req.CheckpointFinalizationTimeout,
		req.NotifyCosmosZone,
		req.BtcStakingContractCodeId,
		req.BtcStakingMsg,
		req.BtcFinalityContractCodeId,
		req.BtcFinalityMsg,
		req.ConsumerName,
		req.ConsumerDescription,
		req.Admin,
	)
	if err != nil {
		return nil, err
	}

	// instantiate the contracts
	babylonContractAddr, btcStakingContractAddr, btcFinalityContractAddr, err := ms.k.InstantiateBabylonContracts(ctx, req.BabylonContractCodeId, initMsg)
	if err != nil {
		return nil, err
	}

	// update params
	params.BabylonContractCodeId = req.BabylonContractCodeId
	params.BtcStakingContractCodeId = req.BtcStakingContractCodeId
	params.BtcFinalityContractCodeId = req.BtcFinalityContractCodeId
	params.BabylonContractAddress = babylonContractAddr
	params.BtcStakingContractAddress = btcStakingContractAddr
	params.BtcFinalityContractAddress = btcFinalityContractAddr
	if err := ms.k.SetParams(ctx, params); err != nil {
		panic(err)
	}

	return &types.MsgInstantiateBabylonContractsResponse{}, nil
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
