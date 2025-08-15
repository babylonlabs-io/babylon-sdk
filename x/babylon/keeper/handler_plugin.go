package keeper

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
)

// integrityHandlerSource interface for staking message access control
type integrityHandlerSource interface {
	CanInvokeStakingMsg(ctx sdk.Context, actor sdk.AccAddress) bool
}

// NewIntegrityHandler prevents any contract from using staking
// or stargate messages. This ensures that staked "virtual" tokens are not bypassing
// the instant undelegate and burn mechanism provided by babylon.
//
// This handler should be chained before any other.
func NewIntegrityHandler(k integrityHandlerSource) wasmkeeper.MessageHandlerFunc {
	return func(ctx sdk.Context, contractAddr sdk.AccAddress, _ string, msg wasmvmtypes.CosmosMsg) (events []sdk.Event, data [][]byte, err error) {
		if msg.Staking == nil || !k.CanInvokeStakingMsg(ctx, contractAddr) {
			return nil, nil, wasmtypes.ErrUnknownMsg // pass down the chain
		}
		// reject
		return nil, nil, types.ErrUnsupported
	}
}

func (k Keeper) CanInvokeStakingMsg(ctx sdk.Context, actor sdk.AccAddress) bool {
	// TODO: implement access control
	return true
}

// CustomMsgHandler handles custom messages from smart contracts
// Currently no custom messages are supported after mint rewards removal
type CustomMsgHandler struct{}

// NewDefaultCustomMsgHandler constructor for the empty custom message handler
func NewDefaultCustomMsgHandler(k *Keeper) *CustomMsgHandler {
	return &CustomMsgHandler{}
}

// DispatchMsg always returns ErrUnknownMsg since no custom messages are supported
func (h CustomMsgHandler) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, _ string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	if msg.Custom == nil {
		return nil, nil, wasmtypes.ErrUnknownMsg
	}

	// No custom messages are currently supported
	return nil, nil, wasmtypes.ErrUnknownMsg
}
