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
	// Validate fee collector account exists
	feeCollector := k.accountKeeper.GetModuleAccount(ctx, k.feeCollectorName)
	if feeCollector == nil {
		return fmt.Errorf("fee collector module account %s not found", k.feeCollectorName)
	}

	feesCollectedInt := k.bank.GetAllBalances(ctx, feeCollector.GetAddress())

	// Check if there are any fees to process
	if feesCollectedInt.IsZero() || !feesCollectedInt.IsAllPositive() {
		k.Logger(ctx).Debug("No positive fees in fee collector")
		return nil
	}

	params := k.GetParams(ctx)
	btcStakingReward := types.GetCoinsPortion(feesCollectedInt, params.BtcStakingPortion)

	if btcStakingReward.IsZero() {
		k.Logger(ctx).Debug("Calculated BTC staking reward is zero, skipping transfer")
		return nil
	}

	// Validate we have sufficient balance for the transfer
	if !feesCollectedInt.IsAllGTE(btcStakingReward) {
		return fmt.Errorf("insufficient fee collector balance")
	}

	contracts := k.GetBSNContracts(ctx)
	if contracts == nil || !contracts.IsSet() {
		return fmt.Errorf("BSN contracts are not set")
	}

	// Validate contract address
	finalityContractAddr, err := sdk.AccAddressFromBech32(contracts.BtcFinalityContract)
	if err != nil {
		return fmt.Errorf("invalid finality contract address %s: %w",
			contracts.BtcFinalityContract, err)
	}

	// Perform the transfer with error handling
	err = k.bank.SendCoinsFromModuleToAccount(ctx, k.feeCollectorName, finalityContractAddr, btcStakingReward)
	if err != nil {
		return fmt.Errorf("bank keeper failed to transfer funds to %s: %w", finalityContractAddr.String(), err)
	}

	k.Logger(ctx).Info("Successfully transferred BTC staking rewards",
		"amount", btcStakingReward,
		"to", finalityContractAddr.String(),
		"portion", params.BtcStakingPortion)

	return nil
}
