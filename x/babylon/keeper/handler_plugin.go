package keeper

import (
	sdkmath "cosmossdk.io/math"
	"encoding/json"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/contract"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
)

// AuthSource abstract type that provides contract authorization.
// This is an extension point for custom implementations.
type AuthSource interface {
	// IsAuthorized returns if the contract authorized to execute a virtual stake message
	IsAuthorized(ctx sdk.Context, contractAddr sdk.AccAddress) bool
}

// abstract keeper
type babylonKeeper interface {
	GetParams(ctx sdk.Context) types.Params
	MintBlockRewards(ctx sdk.Context, recipient sdk.AccAddress, amount sdk.Coin) (sdkmath.Int, error)
}

type CustomMsgHandler struct {
	k    babylonKeeper
	auth AuthSource
}

// NewDefaultCustomMsgHandler constructor to set up the CustomMsgHandler with default max cap authorization
func NewDefaultCustomMsgHandler(k *Keeper) *CustomMsgHandler {
	return &CustomMsgHandler{k: k, auth: defaultMaxCapAuthorizator(k)}
}

// NewCustomMsgHandler constructor to set up CustomMsgHandler with an individual auth source.
// This is an extension point for non default contract authorization logic.
func NewCustomMsgHandler(k babylonKeeper, auth AuthSource) *CustomMsgHandler {
	return &CustomMsgHandler{k: k, auth: auth}
}

func defaultMaxCapAuthorizator(k *Keeper) AuthSourceFn {
	return func(ctx sdk.Context, contractAddr sdk.AccAddress) bool {
		return true
	}
}

// DispatchMsg handle contract message of type Custom in the babylon namespace
func (h CustomMsgHandler) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, _ string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, [][]*codectypes.Any, error) {
	if msg.Custom == nil {
		return nil, nil, nil, wasmtypes.ErrUnknownMsg
	}
	var customMsg contract.CustomMsg
	if err := json.Unmarshal(msg.Custom, &customMsg); err != nil {
		return nil, nil, nil, sdkerrors.ErrJSONUnmarshal.Wrap("custom message")
	}
	if customMsg.MintRewards == nil {
		// not our message type
		return nil, nil, nil, wasmtypes.ErrUnknownMsg
	}

	if !h.auth.IsAuthorized(ctx, contractAddr) {
		return nil, nil, nil, sdkerrors.ErrUnauthorized.Wrapf("contract has no permission for Babylon operations")
	}

	return h.handleMintRewardsMsg(ctx, contractAddr, customMsg.MintRewards)
}

func (h CustomMsgHandler) handleMintRewardsMsg(ctx sdk.Context, actor sdk.AccAddress, mintMsg *contract.MintRewardsMsg) ([]sdk.Event, [][]byte, [][]*codectypes.Any, error) {
	coin, err := wasmkeeper.ConvertWasmCoinToSdkCoin(mintMsg.Amount)
	if err != nil {
		return nil, nil, nil, err
	}
	params := h.k.GetParams(ctx)
	// Validate actor
	if actor.String() != params.BtcFinalityContractAddress {
		return nil, nil, nil, sdkerrors.ErrUnauthorized.Wrapf("minter must be the finality contract")
	}

	// Define recipient
	recipient, err := sdk.AccAddressFromBech32(mintMsg.Recipient)
	if err != nil {
		return nil, nil, nil, err
	}
	if err != nil {
		return nil, nil, nil, err
	}

	rewards, err := h.k.MintBlockRewards(ctx, recipient, coin)
	if err != nil {
		return nil, nil, nil, err
	}

	return []sdk.Event{sdk.NewEvent(
		types.EventTypeMintRewards,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(sdk.AttributeKeySender, actor.String()),
		sdk.NewAttribute(sdk.AttributeKeyAmount, rewards.String()),
	)}, nil, nil, nil
}

// AuthSourceFn is helper for simple AuthSource types
type AuthSourceFn func(ctx sdk.Context, contractAddr sdk.AccAddress) bool

// IsAuthorized returns if the contract authorized to execute a virtual stake message
func (a AuthSourceFn) IsAuthorized(ctx sdk.Context, contractAddr sdk.AccAddress) bool {
	return a(ctx, contractAddr)
}

// abstract keeper
type integrityHandlerSource interface {
	CanInvokeStakingMsg(ctx sdk.Context, actor sdk.AccAddress) bool
}

// NewIntegrityHandler prevents any contract with max cap set to use staking
// or stargate messages. This ensures that staked "virtual" tokens are not bypassing
// the instant undelegate and burn mechanism provided by babylon.
//
// This handler should be chained before any other.
// TODO: access control for msg call from contracts
func NewIntegrityHandler(k integrityHandlerSource) wasmkeeper.MessageHandlerFunc {
	return func(ctx sdk.Context, contractAddr sdk.AccAddress, _ string, msg wasmvmtypes.CosmosMsg) (events []sdk.Event, data [][]byte, msgResponses [][]*codectypes.Any, err error) {
		if msg.Staking == nil || !k.CanInvokeStakingMsg(ctx, contractAddr) {
			return nil, nil, nil, wasmtypes.ErrUnknownMsg // pass down the chain
		}
		// reject
		return nil, nil, nil, types.ErrUnsupported
	}
}

func (k Keeper) CanInvokeStakingMsg(ctx sdk.Context, actor sdk.AccAddress) bool {
	// TODO: implement access control
	return true
}
