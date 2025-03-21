package keeper

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/contract"
	types "github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) InstantiateBabylonContracts(
	ctx sdk.Context,
	babylonContractCodeId uint64,
	initMsg []byte,
) (string, string, string, string, error) {
	contractKeeper := wasmkeeper.NewGovPermissionKeeper(k.wasm)

	// gov address
	govAddr, err := sdk.AccAddressFromBech32(k.authority)
	if err != nil {
		panic(err)
	}

	// instantiate Babylon contract
	babylonContractAddr, _, err := contractKeeper.Instantiate(ctx, babylonContractCodeId, govAddr, govAddr, initMsg, "Babylon contract", nil)
	if err != nil {
		return "", "", "", "", err
	}

	// get contract addresses
	configQuery := []byte(`{"config":{}}`)
	res, err := k.wasm.QuerySmart(ctx, babylonContractAddr, configQuery)
	if err != nil {
		return "", "", "", "", err
	}
	var config types.BabylonContractConfig
	err = json.Unmarshal(res, &config)
	if err != nil {
		return "", "", "", "", err
	}
	if len(config.BTCLightClient) == 0 {
		return "", "", "", "", errorsmod.Wrap(types.ErrInvalid, "failed to instantiate BTC light client contract")
	}
	if len(config.BTCStaking) == 0 {
		return "", "", "", "", errorsmod.Wrap(types.ErrInvalid, "failed to instantiate BTC staking contract")
	}
	if len(config.BTCFinality) == 0 {
		return "", "", "", "", errorsmod.Wrap(types.ErrInvalid, "failed to instantiate BTC finality contract")
	}

	return babylonContractAddr.String(), config.BTCLightClient, config.BTCStaking, config.BTCFinality, nil
}

func (k Keeper) getBTCLightClientContractAddr(ctx sdk.Context) sdk.AccAddress {
	// get address of the BTC light client contract
	addrStr := k.GetParams(ctx).BtcLightClientContractAddress
	if len(addrStr) == 0 {
		// the BTC light client contract address is not set yet, skip sending BeginBlockMsg
		return nil
	}
	addr, err := sdk.AccAddressFromBech32(addrStr)
	if err != nil {
		// Although this is a programming error so we should panic, we emit
		// a warning message to minimise the impact on the consumer chain's operation
		k.Logger(ctx).Warn("the BTC light client contract address is malformed", "contract", addrStr, "error", err)
		return nil
	}
	if !k.wasm.HasContractInfo(ctx, addr) {
		// NOTE: it's possible that the default contract address does not correspond to
		// any contract. We emit a warning message rather than panic to minimise the
		// impact on the consumer chain's operation
		k.Logger(ctx).Warn("the BTC light client contract address is not on-chain", "contract", addrStr)
		return nil
	}
	return addr
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

// SendBeginBlockMsg sends a BeginBlock sudo message to the BTC staking and finality contracts via sudo.
// NOTE: This is a design decision to be made by consumer chains - in this reference implementation,
// if the sudo call fails it will cause consensus failure/chain halt. Consumer chains may want to
// handle sudo call failures differently based on their requirements.
func (k Keeper) SendBeginBlockMsg(c context.Context) error {
	ctx := sdk.UnwrapSDKContext(c)
	headerInfo := ctx.HeaderInfo()

	stakingAddr := k.getBTCStakingContractAddr(ctx)
	finalityAddr := k.getBTCFinalityContractAddr(ctx)
	if stakingAddr == nil || finalityAddr == nil {
		k.Logger(ctx).Info("Skipping begin block processing: contract addresses are missing")
		return nil
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
