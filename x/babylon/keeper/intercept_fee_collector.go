package keeper

import (
	"context"
	"fmt"

	"github.com/babylonlabs-io/babylon/v3/x/incentive/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandleCoinsInFeeCollector intercepts a portion of coins in fee collector and distributes
// to the cosmos BSN contracts.
// It is invoked upon every `BeginBlock`.
// https://github.com/babylonlabs-io/babylon/blob/1a05ecd8dfc69691b6c17637ef520ce9ec302113/x/incentive/keeper/intercept_fee_collector.go#L13
func (k Keeper) HandleCoinsInFeeCollector(ctx context.Context) error {
	feeCollector := k.accountKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	feesCollectedInt := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())

	if !feesCollectedInt.IsAllPositive() {
		k.Logger(ctx).Info("not fees in fee collector")
		return nil
	}

	params := k.GetParams(ctx)
	btcStakingPortion := params.BtcStakingPortion

	// TODO: contract address should be parsed during initiation
	contractAddr, err := sdk.AccAddressFromBech32(params.BabylonContractAddress)
	if err != nil {
		panic(err)
	}

	btcStakingReward := types.GetCoinsPortion(feesCollectedInt, btcStakingPortion)
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.feeCollectorName, contractAddr, btcStakingReward)
	if err != nil {
		return fmt.Errorf("bank keeper failed to transfer funds: %w", err)
	}

	return nil
}
