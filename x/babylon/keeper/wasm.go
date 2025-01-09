package keeper

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/contract"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) InstantiateBabylonContracts(
	ctx sdk.Context,
	babylonContractCodeId uint64,
	btcStakingContractCodeId uint64,
	btcFinalityContractCodeId uint64,
	babylonInitMsg []byte,
	btcStakingInitMsg []byte,
	btcFinalityInitMsg []byte,
) (string, string, string, error) {
	contractKeeper := wasmkeeper.NewGovPermissionKeeper(k.wasm)

	// gov address
	govAddr, err := sdk.AccAddressFromBech32(k.authority)
	if err != nil {
		panic(err)
	}

	// instantiate Babylon contract
	babylonContractAddr, _, err := contractKeeper.Instantiate(ctx, babylonContractCodeId, govAddr, govAddr, babylonInitMsg, "Babylon contract", nil)
	if err != nil {
		return "", "", "", types.ErrInvalid.Wrapf("failed to instantiate Babylon contract: %v", err)
	}

	// instantiate BTC staking contract
	btcStakingContractAddr, _, err := contractKeeper.Instantiate(ctx, btcStakingContractCodeId, govAddr, govAddr, btcStakingInitMsg, "BTC staking contract", nil)
	if err != nil {
		return "", "", "", types.ErrInvalid.Wrapf("failed to instantiate BTC staking contract: %v", err)
	}

	// instantiate BTC finality contract
	btcFinalityContractAddr, _, err := contractKeeper.Instantiate(ctx, btcFinalityContractCodeId, govAddr, govAddr, btcFinalityInitMsg, "BTC finality contract", nil)
	if err != nil {
		return "", "", "", types.ErrInvalid.Wrapf("failed to instantiate BTC finality contract: %v", err)
	}

	return babylonContractAddr.String(), btcStakingContractAddr.String(), btcFinalityContractAddr.String(), nil
}

func (k Keeper) getBTCStakingContractAddr(ctx sdk.Context) sdk.AccAddress {
	// get address of the BTC staking contract
	addrStr := k.GetParams(ctx).BtcStakingContractAddress
	if len(addrStr) == 0 {
		// the BTC staking contract address is not set yet, skip sending BeginBlockMsg
		return nil
	}
	addr, err := sdk.AccAddressFromBech32(addrStr)
	if err != nil {
		// Although this is a programming error so we should panic, we emit
		// a warning message to minimise the impact on the consumer chain's operation
		k.Logger(ctx).Warn("the BTC staking contract address is malformed", "contract", addrStr, "error", err)
		return nil
	}
	if !k.wasm.HasContractInfo(ctx, addr) {
		// NOTE: it's possible that the default contract address does not correspond to
		// any contract. We emit a warning message rather than panic to minimise the
		// impact on the consumer chain's operation
		k.Logger(ctx).Warn("the BTC staking contract address is not on-chain", "contract", addrStr)
		return nil
	}

	return addr
}

func (k Keeper) getBTCFinalityContractAddr(ctx sdk.Context) sdk.AccAddress {
	// get address of the BTC finality contract
	addrStr := k.GetParams(ctx).BtcFinalityContractAddress
	if len(addrStr) == 0 {
		// the BTC finality contract address is not set yet, skip sending BeginBlockMsg
		return nil
	}
	addr, err := sdk.AccAddressFromBech32(addrStr)
	if err != nil {
		// Although this is a programming error so we should panic, we emit
		// a warning message to minimise the impact on the consumer chain's operation
		k.Logger(ctx).Warn("the BTC finality contract address is malformed", "contract", addrStr, "error", err)
		return nil
	}
	if !k.wasm.HasContractInfo(ctx, addr) {
		// NOTE: it's possible that the default contract address does not correspond to
		// any contract. We emit a warning message rather than panic to minimise the
		// impact on the consumer chain's operation
		k.Logger(ctx).Warn("the BTC finality contract address is not on-chain", "contract", addrStr)
		return nil
	}

	return addr
}

// SendBeginBlockMsg sends a BeginBlock sudo message to the BTC finality contract via sudo
func (k Keeper) SendBeginBlockMsg(c context.Context) error {
	ctx := sdk.UnwrapSDKContext(c)

	// try to get and parse BTC finality contract
	addr := k.getBTCFinalityContractAddr(ctx)
	if addr == nil {
		return nil
	}

	// construct the sudo message
	headerInfo := ctx.HeaderInfo()
	msg := contract.SudoMsg{
		BeginBlockMsg: &contract.BeginBlock{
			HashHex:    hex.EncodeToString(headerInfo.Hash),
			AppHashHex: hex.EncodeToString(headerInfo.AppHash),
		},
	}

	// send the sudo call
	return k.doSudoCall(ctx, addr, msg)
}

// SendEndBlockMsg sends a EndBlock sudo message to the BTC finality contract via sudo
func (k Keeper) SendEndBlockMsg(c context.Context) error {
	ctx := sdk.UnwrapSDKContext(c)

	// try to get and parse BTC finality contract
	addr := k.getBTCFinalityContractAddr(ctx)
	if addr == nil {
		return nil
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
	return k.doSudoCall(ctx, addr, msg)
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
