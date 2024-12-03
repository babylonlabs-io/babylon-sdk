package keeper_test

import (
	"encoding/json"
	"testing"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
	"github.com/stretchr/testify/require"
)

// TODO: fix this test
func TestInstantiateBabylonContracts(t *testing.T) {
	keepers := NewTestKeepers(t)
	msgServer := keepers.BabylonMsgServer
	wasmMsgServer := keepers.WasmMsgServer

	// store Babylon contract codes
	babylonContractCode, btcStakingContractCode, btcFinalityContractCode := GetGZippedContractCodes()
	resp, err := wasmMsgServer.StoreCode(keepers.Ctx, &wasmtypes.MsgStoreCode{
		Sender:       keepers.BabylonKeeper.GetAuthority(),
		WASMByteCode: babylonContractCode,
	})
	babylonContractCodeID := resp.CodeID
	require.NoError(t, err)
	resp, err = wasmMsgServer.StoreCode(keepers.Ctx, &wasmtypes.MsgStoreCode{
		Sender:       keepers.BabylonKeeper.GetAuthority(),
		WASMByteCode: btcStakingContractCode,
	})
	btcStakingContractCodeID := resp.CodeID
	require.NoError(t, err)
	resp, err = wasmMsgServer.StoreCode(keepers.Ctx, &wasmtypes.MsgStoreCode{
		Sender:       keepers.BabylonKeeper.GetAuthority(),
		WASMByteCode: btcFinalityContractCode,
	})
	btcFinalityContractCodeID := resp.CodeID
	require.NoError(t, err)

	// BTC staking init message
	btcStakingInitMsg := map[string]interface{}{
		"admin": keepers.BabylonKeeper.GetAuthority(),
	}
	btcStakingInitMsgBytes, err := json.Marshal(btcStakingInitMsg)
	require.NoError(t, err)
	// BTC finality init message
	btcFinalityInitMsg := map[string]interface{}{
		"admin": keepers.BabylonKeeper.GetAuthority(),
	}
	btcFinalityInitMsgBytes, err := json.Marshal(btcFinalityInitMsg)
	require.NoError(t, err)

	// instantiate Babylon contract
	_, err = msgServer.InstantiateBabylonContracts(keepers.Ctx, &types.MsgInstantiateBabylonContracts{
		Network:                       "regtest",
		BabylonContractCodeId:         babylonContractCodeID,
		BtcStakingContractCodeId:      btcStakingContractCodeID,
		BtcFinalityContractCodeId:     btcFinalityContractCodeID,
		BabylonTag:                    "01020304",
		BtcConfirmationDepth:          1,
		CheckpointFinalizationTimeout: 2,
		NotifyCosmosZone:              false,
		BtcStakingMsg:                 btcStakingInitMsgBytes,
		BtcFinalityMsg:                btcFinalityInitMsgBytes,
		ConsumerName:                  "test-consumer",
		ConsumerDescription:           "test-consumer-description",
		Admin:                         keepers.BabylonKeeper.GetAuthority(),
	})
	require.NoError(t, err)
}
