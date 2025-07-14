package keeper

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/contract"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SendBeginBlockMsg sends a BeginBlock sudo message to the BTC staking and finality contracts via sudo.
// NOTE: This is a design decision to be made by consumer chains - in this reference implementation,
// if the sudo call fails it will cause consensus failure/chain halt. Consumer chains may want to
// handle sudo call failures differently based on their requirements.
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
		return err
	}
	finalityAddr, err := sdk.AccAddressFromBech32(contracts.BtcFinalityContract)
	if err != nil {
		return err
	}

	// Send the sudo call to the BTC staking contract
	stakingMsg := contract.SudoMsg{
		BeginBlockMsg: &contract.BeginBlock{
			HashHex:    hex.EncodeToString(headerInfo.Hash),
			AppHashHex: hex.EncodeToString(headerInfo.AppHash),
		},
	}
	if err := k.doSudoCall(ctx, stakingAddr, stakingMsg); err != nil {
		return err
	}

	// Send the sudo call to the finality contract
	finalityMsg := contract.SudoMsg{
		BeginBlockMsg: &contract.BeginBlock{
			HashHex:    hex.EncodeToString(headerInfo.Hash),
			AppHashHex: hex.EncodeToString(headerInfo.AppHash),
		},
	}
	if err := k.doSudoCall(ctx, finalityAddr, finalityMsg); err != nil {
		return err
	}

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
		return err
	}

	// construct the sudo message
	headerInfo := ctx.HeaderInfo()
	msg := contract.SudoMsg{
		EndBlockMsg: &contract.EndBlock{
			HashHex:    hex.EncodeToString(headerInfo.Hash),
			AppHashHex: hex.EncodeToString(headerInfo.AppHash),
		},
	}

	// send the sudo call
	return k.doSudoCall(ctx, finalityAddr, msg)
}

// caller must ensure gas limits are set proper and handle panics
func (k Keeper) doSudoCall(ctx sdk.Context, contractAddr sdk.AccAddress, msg contract.SudoMsg) error {
	bz, err := json.Marshal(msg)
	if err != nil {
		return errorsmod.Wrap(err, "marshal sudo msg")
	}
	resp, err := k.wasm.Sudo(ctx, contractAddr, bz)
	k.Logger(ctx).Debug(fmt.Sprintf("response of sudo call %v to contract %s: %v", bz, contractAddr.String(), resp))
	return err
}
