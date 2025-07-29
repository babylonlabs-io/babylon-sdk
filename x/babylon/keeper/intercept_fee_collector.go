package keeper

import (
	"fmt"

	"github.com/babylonlabs-io/babylon/v3/x/incentive/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// HandleCoinsInFeeCollector intercepts a portion of coins in fee collector and distributes
// to the cosmos BSN contracts.
// It is invoked upon every `BeginBlock`.
// https://github.com/babylonlabs-io/babylon/blob/1a05ecd8dfc69691b6c17637ef520ce9ec302113/x/incentive/keeper/intercept_fee_collector.go#L13
func (k Keeper) HandleCoinsInFeeCollector(ctx sdk.Context) error {
	feeCollector := k.accountKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	feesCollectedInt := k.bank.GetAllBalances(ctx, feeCollector.GetAddress())

	if !feesCollectedInt.IsAllPositive() {
		k.Logger(ctx).Info("not fees in fee collector")
		return nil
	}

	params := k.GetParams(ctx)
	btcStakingPortion := params.BtcStakingPortion

	contracts := k.GetBSNContracts(ctx)
	if contracts == nil || !contracts.IsSet() {
		return fmt.Errorf("BSN contracts are not set")
	}

	// send the collected fee to the finality contracts which will handle ICS20 transfer
	contractAddr, err := sdk.AccAddressFromBech32(contracts.BtcFinalityContract)
	if err != nil {
		panic(err)
	}

	btcStakingReward := types.GetCoinsPortion(feesCollectedInt, btcStakingPortion)
	err = k.bank.SendCoinsFromModuleToAccount(ctx, k.feeCollectorName, contractAddr, btcStakingReward)
	if err != nil {
		return fmt.Errorf("bank keeper failed to transfer funds: %w", err)
	}

	return nil
}
