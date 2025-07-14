package keeper

import (
	"context"
)

// HandleCoinsInFeeCollector intercepts a portion of coins in fee collector and distributes
// to the cosmos BSN contracts.
// It is invoked upon every `BeginBlock`.
// https://github.com/babylonlabs-io/babylon/blob/1a05ecd8dfc69691b6c17637ef520ce9ec302113/x/incentive/keeper/intercept_fee_collector.go#L13
func (k Keeper) HandleCoinsInFeeCollector(ctx context.Context) {
	// find the fee collector account
	feeCollector := k.accountKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	// get all balances in the fee collector account,
	// where the balance includes minted tokens in the previous block
	feesCollectedInt := k.bankKeeper.GetAllBalances(ctx, feeCollector.GetAddress())

	// don't intercept if there is no fee in fee collector account
	if !feesCollectedInt.IsAllPositive() {
		return
	}

	// TODO: get reward portion and sends it to the cosmos BSN contracts
}
