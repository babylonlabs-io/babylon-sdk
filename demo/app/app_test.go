package app

import (
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/client/cli"
	babylonkeeper "github.com/babylonlabs-io/babylon-sdk/x/babylon/keeper"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var emptyWasmOpts []wasm.Option

// adapted from https://github.com/cosmos/cosmos-sdk/blob/v0.50.6/simapp/app_test.go#L47-L48
func TestSimAppExportAndBlockedAddrs(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewTestLogger(t)
	app := NewAppWithCustomOptions(t, false, SetupOptions{
		Logger:  logger.With("instance", "first"),
		DB:      db,
		AppOpts: simtestutil.NewAppOptionsWithFlagHome(t.TempDir()),
	})

	// BlockedAddresses returns a map of addresses in app v1 and a map of modules name in app v2.
	for acc := range BlockedAddresses() {
		var addr sdk.AccAddress
		if modAddr, err := sdk.AccAddressFromBech32(acc); err == nil {
			addr = modAddr
		} else {
			addr = app.AccountKeeper.GetModuleAddress(acc)
		}

		require.True(
			t,
			app.BankKeeper.BlockedAddr(addr),
			fmt.Sprintf("ensure that blocked addresses are properly set in bank keeper: %s should be blocked", acc),
		)
	}

	// finalize block so we have CheckTx state set
	_, err := app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height: 1,
	})
	require.NoError(t, err)

	_, err = app.Commit()
	require.NoError(t, err)

	// Making a new app object with the db, so that initchain hasn't been called
	app2 := NewConsumerApp(logger, db, nil, true, simtestutil.NewAppOptionsWithFlagHome(t.TempDir()), emptyWasmOpts)
	_, err = app2.ExportAppStateAndValidators(false, []string{}, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

func TestGetMaccPerms(t *testing.T) {
	dup := GetMaccPerms()
	require.Equal(t, maccPerms, dup, "duplicated module account permissions differed from actual module account permissions")
}

const (
	TestDataPath                   = "../../tests/testdata"
	BabylonContractCodePath        = TestDataPath + "/babylon_contract.wasm"
	BtcLightClientContractCodePath = TestDataPath + "/btc_light_client.wasm"
	BtcStakingContractCodePath     = TestDataPath + "/btc_staking.wasm"
	BtcFinalityContractCodePath    = TestDataPath + "/btc_finality.wasm"
)

func GetGZippedContractCodes() ([]byte, []byte, []byte, []byte) {
	babylonContractCode, err := types.GetGZippedContractCode(BabylonContractCodePath)
	if err != nil {
		panic(err)
	}
	btcLightClientContractCode, err := types.GetGZippedContractCode(BtcLightClientContractCodePath)
	if err != nil {
		panic(err)
	}
	btcStakingContractCode, err := types.GetGZippedContractCode(BtcStakingContractCodePath)
	if err != nil {
		panic(err)
	}
	btcFinalityContractCode, err := types.GetGZippedContractCode(BtcFinalityContractCodePath)
	if err != nil {
		panic(err)
	}

	return babylonContractCode, btcLightClientContractCode, btcStakingContractCode, btcFinalityContractCode
}

func TestInstantiateBabylonContracts(t *testing.T) {
	consumerApp := Setup(t)
	ctx := consumerApp.NewContext(false)
	ctx = ctx.WithBlockHeader(cmtproto.Header{Time: time.Now()})
	babylonKeeper := consumerApp.BabylonKeeper
	babylonMsgServer := babylonkeeper.NewMsgServer(babylonKeeper)
	wasmKeeper := consumerApp.WasmKeeper
	wasmMsgServer := wasmkeeper.NewMsgServerImpl(&wasmKeeper)

	// store Babylon contract codes
	babylonContractCode, btcLightClientContractCode, btcStakingContractCode, btcFinalityContractCode := GetGZippedContractCodes()
	resp, err := wasmMsgServer.StoreCode(ctx, &wasmtypes.MsgStoreCode{
		Sender:       consumerApp.BabylonKeeper.GetAuthority(),
		WASMByteCode: babylonContractCode,
	})
	babylonContractCodeID := resp.CodeID
	require.NoError(t, err)
	resp, err = wasmMsgServer.StoreCode(ctx, &wasmtypes.MsgStoreCode{
		Sender:       consumerApp.BabylonKeeper.GetAuthority(),
		WASMByteCode: btcLightClientContractCode,
	})
	btcLightClientContractCodeID := resp.CodeID
	require.NoError(t, err)
	resp, err = wasmMsgServer.StoreCode(ctx, &wasmtypes.MsgStoreCode{
		Sender:       consumerApp.BabylonKeeper.GetAuthority(),
		WASMByteCode: btcStakingContractCode,
	})
	btcStakingContractCodeID := resp.CodeID
	require.NoError(t, err)
	resp, err = wasmMsgServer.StoreCode(ctx, &wasmtypes.MsgStoreCode{
		Sender:       consumerApp.BabylonKeeper.GetAuthority(),
		WASMByteCode: btcFinalityContractCode,
	})
	btcFinalityContractCodeID := resp.CodeID
	require.NoError(t, err)

	initMsg, err := cli.ParseInstantiateArgs(
		[]string{
			fmt.Sprintf("%d", babylonContractCodeID),
			fmt.Sprintf("%d", btcLightClientContractCodeID),
			fmt.Sprintf("%d", btcStakingContractCodeID),
			fmt.Sprintf("%d", btcFinalityContractCodeID),
			"regtest",
			"01020304",
			"1",
			"2",
			"false",
			"test-consumer",
			"test-consumer-description",
		},
		"",
		babylonKeeper.GetAuthority(),
		babylonKeeper.GetAuthority(),
	)
	require.NoError(t, err)

	// instantiate Babylon contract
	_, err = babylonMsgServer.InstantiateBabylonContracts(ctx, initMsg)
	require.NoError(t, err, initMsg)
}
