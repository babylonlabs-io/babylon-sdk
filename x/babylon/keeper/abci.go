package keeper

import (
	"context"
	"time"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
)

func (k *Keeper) BeginBlocker(ctx context.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	return k.SendBeginBlockMsg(ctx)
}

// EndBlocker is called after every block
func (k *Keeper) EndBlocker(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	if err := k.SendEndBlockMsg(ctx); err != nil {
		return []abci.ValidatorUpdate{}, err
	}

	return []abci.ValidatorUpdate{}, nil
}
