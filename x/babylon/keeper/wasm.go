package keeper

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/contract"
)

// SendBeginBlockMsg sends a BeginBlock sudo message to the BTC staking and finality contracts via sudo.
func (k Keeper) SendBeginBlockMsg(c context.Context) error {
	ctx := sdk.UnwrapSDKContext(c)
	headerInfo := ctx.HeaderInfo()

	contracts := k.GetBSNContracts(ctx)
	if contracts == nil || !contracts.IsSet() {
		k.Logger(ctx).Info("Skipping begin block processing: contract addresses are missing")
		return nil
	}

	stakingAddr, err := sdk.AccAddressFromBech32(contracts.BtcStakingContract)
	if err != nil {
		return fmt.Errorf("invalid BTC staking contract address %s: %w", contracts.BtcStakingContract, err)
	}
	finalityAddr, err := sdk.AccAddressFromBech32(contracts.BtcFinalityContract)
	if err != nil {
		return fmt.Errorf("invalid BTC finality contract address %s: %w", contracts.BtcFinalityContract, err)
	}

	// Send the sudo call to the BTC staking contract with gas limits
	headerHashHex := hex.EncodeToString(headerInfo.Hash)
	appHashHex := hex.EncodeToString(headerInfo.AppHash)
	stakingMsg := contract.SudoMsg{
		BeginBlockMsg: &contract.BeginBlock{
			HashHex:    headerHashHex,
			AppHashHex: appHashHex,
		},
	}

	maxGas := k.GetMaxSudoGasBeginBlocker(ctx)

	gasConsumed, err := k.doSudoCallWithGasLimit(ctx, stakingAddr, stakingMsg, maxGas)
	if err != nil {
		return fmt.Errorf("failed to send BeginBlock message to BTC staking contract %s: %w",
			stakingAddr.String(), err)
	}
	k.Logger(ctx).Debug("BeginBlock sudo call to BTC staking contract successful",
		"contract", stakingAddr.String(),
		"gas_used", gasConsumed)

	// Send the sudo call to the finality contract with gas limits
	finalityMsg := contract.SudoMsg{
		BeginBlockMsg: &contract.BeginBlock{
			HashHex:    headerHashHex,
			AppHashHex: appHashHex,
		},
	}
	gasConsumed, err = k.doSudoCallWithGasLimit(ctx, finalityAddr, finalityMsg, maxGas)
	if err != nil {
		return fmt.Errorf("failed to send BeginBlock message to BTC finality contract %s: %w",
			finalityAddr.String(), err)
	}
	k.Logger(ctx).Debug("BeginBlock sudo call to BTC finality contract successful",
		"contract", finalityAddr.String(),
		"gas_used", gasConsumed)

	return nil
}

// SendEndBlockMsg sends a EndBlock sudo message to the BTC finality contract via sudo
func (k Keeper) SendEndBlockMsg(c context.Context) error {
	ctx := sdk.UnwrapSDKContext(c)

	contracts := k.GetBSNContracts(ctx)
	if contracts == nil || !contracts.IsSet() {
		k.Logger(ctx).Info("Skipping end block processing: contract addresses are missing")
		return nil
	}

	finalityAddr, err := sdk.AccAddressFromBech32(contracts.BtcFinalityContract)
	if err != nil {
		return fmt.Errorf("invalid BTC finality contract address %s: %w", contracts.BtcFinalityContract, err)
	}

	// construct the sudo message
	headerInfo := ctx.HeaderInfo()
	msg := contract.SudoMsg{
		EndBlockMsg: &contract.EndBlock{
			HashHex:    hex.EncodeToString(headerInfo.Hash),
			AppHashHex: hex.EncodeToString(headerInfo.AppHash),
		},
	}

	// send the sudo call with gas limits
	gasConsumed, err := k.doSudoCallWithGasLimit(ctx, finalityAddr, msg, k.GetMaxSudoGasEndBlocker(ctx))
	if err != nil {
		k.Logger(ctx).Error("Failed to send EndBlock message to BTC finality contract", "error", err)
		return fmt.Errorf("BTC finality contract EndBlock call failed: %w", err)
	}
	k.Logger(ctx).Debug("EndBlock sudo call to BTC finality contract successful",
		"contract", finalityAddr.String(),
		"gas_used", gasConsumed)

	return nil
}

// caller must ensure gas limits are set proper and handle panics
func (k Keeper) doSudoCall(ctx sdk.Context, contractAddr sdk.AccAddress, msg contract.SudoMsg) ([]byte, error) {
	bz, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sudo message: %w", err)
	}

	return k.wasm.Sudo(ctx, contractAddr, bz)
}

// doSudoCallWithGasLimit performs a sudo call with gas limit protection and error recovery
func (k Keeper) doSudoCallWithGasLimit(ctx sdk.Context, contractAddr sdk.AccAddress, msg contract.SudoMsg, maxGas storetypes.Gas) (gasConsumed storetypes.Gas, err error) {
	gasConsumed = 0

	if maxGas == 0 {
		err = fmt.Errorf("max gas cannot be zero")
		return
	}

	// Create a cache context to ensure atomicity; commit only on success
	cacheCtx, write := ctx.CacheContext()
	// Create a gas-limited context on top of the cache context
	gasCtx := cacheCtx.WithGasMeter(storetypes.NewGasMeter(maxGas))

	// Use defer to recover from panics that might occur during contract execution
	defer func() {
		if r := recover(); r != nil {
			gasConsumed = gasCtx.GasMeter().GasConsumed()
			err = fmt.Errorf("contract call to %s panicked: %v, gas_used: %d",
				contractAddr.String(), r, gasConsumed)
		}
	}()

	if _, err = k.doSudoCall(gasCtx, contractAddr, msg); err != nil {
		gasConsumed = gasCtx.GasMeter().GasConsumed()
		k.Logger(ctx).Error("Sudo call failed",
			"contract", contractAddr.String(),
			"error", err,
			"gas_used", gasConsumed)
		err = errorsmod.Wrapf(err, "sudo call to contract failed")
		return
	}

	// Commit cached writes only after successful execution
	write()
	gasConsumed = gasCtx.GasMeter().GasConsumed()

	return
}
