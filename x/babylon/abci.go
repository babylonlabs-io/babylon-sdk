package babylon

import (
	"context"
	"fmt"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/keeper"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
)

func BeginBlocker(ctx context.Context, k keeper.Keeper) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	// handle coins in the fee collector account
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if sdkCtx.HeaderInfo().Height > 0 {
		if err := k.HandleCoinsInFeeCollector(sdkCtx); err != nil {
			k.Logger(sdkCtx).Error("BeginBlocker failed to handle coins in fee collector", "error", err)
			// Emit an alert event for monitoring systems
			sdkCtx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeFeeCollectorError,
					sdk.NewAttribute(types.AttributeKeyError, err.Error()),
					sdk.NewAttribute(types.AttributeKeyHeight, fmt.Sprintf("%d", sdkCtx.HeaderInfo().Height)),
				),
			)
			// not return error to not cause panic
		}
	}

	// send BeginBlocker message to contracts and handle contract communication errors gracefully
	if err := k.SendBeginBlockMsg(ctx); err != nil {
		k.Logger(sdkCtx).Error("BeginBlocker failed to send message to contracts", "error", err)
		// Emit an alert event for monitoring systems
		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeContractCommunicationError,
				sdk.NewAttribute(types.AttributeKeyError, err.Error()),
				sdk.NewAttribute(types.AttributeKeyHeight, fmt.Sprintf("%d", sdkCtx.HeaderInfo().Height)),
				sdk.NewAttribute(types.AttributeKeyPhase, "BeginBlock"),
			),
		)
		// not return error to not cause panic
	}

	return nil
}

// EndBlocker is called after every block
func EndBlocker(ctx context.Context, k keeper.Keeper) ([]abci.ValidatorUpdate, error) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if err := k.SendEndBlockMsg(ctx); err != nil {
		k.Logger(sdkCtx).Error("EndBlocker failed to send message to contracts", "error", err)
		// Emit an alert event for monitoring systems
		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeContractCommunicationError,
				sdk.NewAttribute(types.AttributeKeyError, err.Error()),
				sdk.NewAttribute(types.AttributeKeyHeight, fmt.Sprintf("%d", sdkCtx.HeaderInfo().Height)),
				sdk.NewAttribute(types.AttributeKeyPhase, "EndBlock"),
			),
		)
		// not return error to not cause panic
	}

	return []abci.ValidatorUpdate{}, nil
}
