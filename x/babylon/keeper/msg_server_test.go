package keeper_test

import (
	"encoding/json"
	"testing"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
	"github.com/stretchr/testify/require"
)

func TestStoreBabylonContractCodes(t *testing.T) {
	keepers := NewTestKeepers(t)
	msgServer := keepers.BabylonMsgServer

	babylonContractCode, btcStakingContractCode, btcFinalityContractCode := GetGZippedContractCodes()

	// store Babylon contract codes
	_, err := msgServer.StoreBabylonContractCodes(keepers.Ctx, &types.MsgStoreBabylonContractCodes{
		BabylonContractCode:     babylonContractCode,
		BtcStakingContractCode:  btcStakingContractCode,
		BtcFinalityContractCode: btcFinalityContractCode,
	})
	require.NoError(t, err)

	// ensure params are set
	params := keepers.BabylonKeeper.GetParams(keepers.Ctx)
	require.Positive(t, params.BabylonContractCodeId)
	require.Positive(t, params.BtcStakingContractCodeId)
	require.Positive(t, params.BtcFinalityContractCodeId)

	// ensure non-gov account cannot override
	_, err = msgServer.StoreBabylonContractCodes(keepers.Ctx, &types.MsgStoreBabylonContractCodes{
		BabylonContractCode:     babylonContractCode,
		BtcStakingContractCode:  btcStakingContractCode,
		BtcFinalityContractCode: btcFinalityContractCode,
	})
	require.Error(t, err)

	// gov can override
	_, err = msgServer.StoreBabylonContractCodes(keepers.Ctx, &types.MsgStoreBabylonContractCodes{
		Signer:                  keepers.BabylonKeeper.GetAuthority(),
		BabylonContractCode:     babylonContractCode,
		BtcStakingContractCode:  btcStakingContractCode,
		BtcFinalityContractCode: btcFinalityContractCode,
	})
	require.NoError(t, err)
}

// TODO: fix this test
func TestInstantiateBabylonContracts(t *testing.T) {
	keepers := NewTestKeepers(t)
	msgServer := keepers.BabylonMsgServer

	// store Babylon contract codes
	babylonContractCode, btcStakingContractCode, btcFinalityContractCode := GetGZippedContractCodes()
	_, err := msgServer.StoreBabylonContractCodes(keepers.Ctx, &types.MsgStoreBabylonContractCodes{
		BabylonContractCode:     babylonContractCode,
		BtcStakingContractCode:  btcStakingContractCode,
		BtcFinalityContractCode: btcFinalityContractCode,
	})
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
		BabylonTag:                    "01020304",
		BtcConfirmationDepth:          1,
		CheckpointFinalizationTimeout: 2,
		NotifyCosmosZone:              false,
		BtcStakingMsg:                 btcStakingInitMsgBytes,
		BtcFinalityMsg:                btcFinalityInitMsgBytes,
		ConsumerName:                  "test-consumer",
		ConsumerDescription:           "test-consumer-description",
	})
	require.NoError(t, err)
}
