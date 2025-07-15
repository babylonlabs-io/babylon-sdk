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
			fmt.Sprintf(`{"admin":"%s"}`, babylonKeeper.GetAuthority()),
			fmt.Sprintf(`{"admin":"%s"}`, babylonKeeper.GetAuthority()),
			`{"header": {"version": 536870912, "prev_blockhash": "000000c0a3841a6ae64c45864ae25314b40fd522bfb299a4b6bd5ef288cae74d", "merkle_root": "e666a9797b7a650597098ca6bf500bd0873a86ada05189f87073b6dfdbcaf4ee", "time": 1599332844, "bits": 503394215, "nonce": 9108535}, "height": 2016, "total_work": "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAkY98OU="}`,
			"test-consumer",
			"test-consumer-description",
		},
		"",
		babylonKeeper.GetAuthority(),
		babylonKeeper.GetAuthority(),
	)
	require.NoError(t, err)

	// instantiate Babylon contract
	msgResp, err := babylonMsgServer.InstantiateBabylonContracts(ctx, initMsg)
	require.NoError(t, err)
	require.NotNil(t, msgResp)

	// Verify that the contracts were instantiated successfully
	params := babylonKeeper.GetParams(ctx)
	babylonAddr, btcLightClientAddr, btcStakingAddr, btcFinalityAddr, err := params.GetContractAddresses()
	require.NoError(t, err)
	require.NotEmpty(t, babylonAddr)
	require.NotEmpty(t, btcLightClientAddr)
	require.NotEmpty(t, btcStakingAddr)
	require.NotEmpty(t, btcFinalityAddr)

	// Verify that the contract code IDs are set correctly
	require.Equal(t, babylonContractCodeID, params.BabylonContractCodeId)
	require.Equal(t, btcLightClientContractCodeID, params.BtcLightClientContractCodeId)
	require.Equal(t, btcStakingContractCodeID, params.BtcStakingContractCodeId)
	require.Equal(t, btcFinalityContractCodeID, params.BtcFinalityContractCodeId)

	// Verify that the contracts are instantiated
	require.True(t, wasmKeeper.HasContractInfo(ctx, babylonAddr))
	require.True(t, wasmKeeper.HasContractInfo(ctx, btcLightClientAddr))
	require.True(t, wasmKeeper.HasContractInfo(ctx, btcStakingAddr))
	require.True(t, wasmKeeper.HasContractInfo(ctx, btcFinalityAddr))
}
