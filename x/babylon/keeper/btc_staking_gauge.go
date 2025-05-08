package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
)

func (k Keeper) accumulateBTCStakingReward(ctx context.Context, btcStakingReward sdk.Coins) {
	// update BTC staking gauge
	height := uint64(sdk.UnwrapSDKContext(ctx).HeaderInfo().Height)
	gauge := types.NewGauge(btcStakingReward...)
	k.SetBTCStakingGauge(ctx, height, gauge)

	// transfer the BTC staking reward from fee collector account to incentive module account
	err := k.bank.SendCoinsFromModuleToModule(ctx, k.feeCollectorName, types.ModuleName, btcStakingReward)
	if err != nil {
		// this can only be programming error and is unrecoverable
		panic(err)
	}
}

func (k Keeper) SetBTCStakingGauge(ctx context.Context, height uint64, gauge *types.Gauge) {
	store := k.btcStakingGaugeStore(ctx)
	gaugeBytes := k.cdc.MustMarshal(gauge)
	store.Set(sdk.Uint64ToBigEndian(height), gaugeBytes)
}

func (k Keeper) GetBTCStakingGauge(ctx context.Context, height uint64) *types.Gauge {
	store := k.btcStakingGaugeStore(ctx)
	gaugeBytes := store.Get(sdk.Uint64ToBigEndian(height))
	if gaugeBytes == nil {
		return nil
	}

	var gauge types.Gauge
	k.cdc.MustUnmarshal(gaugeBytes, &gauge)
	return &gauge
}

// btcStakingGaugeStore returns the KVStore of the gauge of total reward for
// BTC staking at each height
// prefix: BTCStakingGaugeKey
// key: gauge height
// value: gauge of rewards for BTC staking at this height
func (k Keeper) btcStakingGaugeStore(ctx context.Context) prefix.Store {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return prefix.NewStore(storeAdapter, types.BTCStakingGaugeKey)
}
