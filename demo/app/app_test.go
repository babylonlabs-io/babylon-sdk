package app

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/babylonlabs-io/babylon-sdk/tests/e2e/datagen"
	babylonkeeper "github.com/babylonlabs-io/babylon-sdk/x/babylon/keeper"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
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

// Config represents the configuration for the Babylon contract
type babylonContractConfig struct {
	Network                       string `json:"network"`
	BabylonTag                    []byte `json:"babylon_tag"`
	BTCConfirmationDepth          uint32 `json:"btc_confirmation_depth"`
	CheckpointFinalizationTimeout uint32 `json:"checkpoint_finalization_timeout"`
	NotifyCosmosZone              bool   `json:"notify_cosmos_zone"`
	BTCLightClient                string `json:"btc_light_client,omitempty"`
	BTCStaking                    string `json:"btc_staking,omitempty"`
	BTCFinality                   string `json:"btc_finality,omitempty"`
	ConsumerName                  string `json:"consumer_name,omitempty"`
	ConsumerDescription           string `json:"consumer_description,omitempty"`
}

func TestInstantiateBabylonContracts(t *testing.T) {
	consumerApp := Setup(t)
	ctx := consumerApp.NewContext(false)
	ctx = ctx.WithBlockHeader(cmtproto.Header{Time: time.Now()})
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

	network := "regtest"
	btcConfirmationDepth := 1
	btcFinalizationTimeout := 2
	babylonAdmin := consumerApp.BabylonKeeper.GetAuthority()
	btcLightClientInitMsg := fmt.Sprintf(`{"network":"%s","btc_confirmation_depth":%d,"checkpoint_finalization_timeout":%d,"initial_header":%s}`,
		network, btcConfirmationDepth, btcFinalizationTimeout, datagen.MustGetInitialHeaderInStr())
	btcFinalityInitMsg := fmt.Sprintf(`{"admin":"%s"}`, babylonAdmin)
	btcStakingInitMsg := fmt.Sprintf(`{"admin":"%s"}`, babylonAdmin)

	// Base64 encode the init messages as required by the contract schemas
	btcLightClientInitMsgBz := base64.StdEncoding.EncodeToString([]byte(btcLightClientInitMsg))
	btcFinalityInitMsgBz := base64.StdEncoding.EncodeToString([]byte(btcFinalityInitMsg))
	btcStakingInitMsgBz := base64.StdEncoding.EncodeToString([]byte(btcStakingInitMsg))

	babylonInitMsg := map[string]interface{}{
		"network":                         network,
		"babylon_tag":                     "01020304",
		"btc_confirmation_depth":          btcConfirmationDepth,
		"checkpoint_finalization_timeout": btcFinalizationTimeout,
		"notify_cosmos_zone":              false,
		"btc_light_client_code_id":        btcLightClientContractCodeID,
		"btc_light_client_msg":            btcLightClientInitMsgBz,
		"btc_staking_code_id":             btcStakingContractCodeID,
		"btc_staking_msg":                 btcStakingInitMsgBz,
		"btc_finality_code_id":            btcFinalityContractCodeID,
		"btc_finality_msg":                btcFinalityInitMsgBz,
		"consumer_name":                   "test-consumer",
		"btc_light_client_initial_header": datagen.MustGetInitialHeaderInHex(),
		"consumer_description":            "test-consumer-description",
		"ics20_channel_id":                "channel-0",
		"destination_module":              "btcstaking",
	}
	babylonInitMsgBz, err := json.Marshal(babylonInitMsg)
	require.NoError(t, err)

	instResp, err := wasmMsgServer.InstantiateContract(ctx, &wasmtypes.MsgInstantiateContract{
		Sender: babylonAdmin,
		Admin:  babylonAdmin,
		CodeID: babylonContractCodeID,
		Label:  "test-contract",
		Msg:    babylonInitMsgBz,
		Funds:  nil,
	})
	require.NoError(t, err)
	require.NotEmpty(t, instResp.Address)
	// Debug: print the contract address
	babylonAddr := string(instResp.Address)
	t.Logf("Instantiated Babylon contract address (string): %s", babylonAddr)
	t.Logf("Instantiated Babylon contract address (bytes): %x", instResp.Address)

	babylonAccAddr, err := sdk.AccAddressFromBech32(babylonAddr)
	require.NoError(t, err)

	// Check if the contract info exists in the keeper before querying
	if !wasmKeeper.HasContractInfo(ctx, babylonAccAddr) {
		t.Fatalf("Wasm keeper does not have contract info for address: %s", babylonAddr)
	}

	// Extract contract addresses from SDK context events
	var btcLightClientAddr, btcStakingAddr, btcFinalityAddr string
	events := ctx.EventManager().Events()
	for _, event := range events {
		if event.Type == "instantiate" {
			var addr string
			var codeID string
			for _, attr := range event.Attributes {
				if string(attr.Key) == "_contract_address" || string(attr.Key) == "contract_address" {
					addr = string(attr.Value)
				}
				if string(attr.Key) == "code_id" {
					codeID = string(attr.Value)
				}
			}
			// Map by code ID to get the sub-contract addresses
			switch codeID {
			case fmt.Sprintf("%d", btcLightClientContractCodeID):
				btcLightClientAddr = addr
				t.Logf("BTC Light Client contract address: %s", addr)
			case fmt.Sprintf("%d", btcStakingContractCodeID):
				btcStakingAddr = addr
				t.Logf("BTC Staking contract address: %s", addr)
			case fmt.Sprintf("%d", btcFinalityContractCodeID):
				btcFinalityAddr = addr
				t.Logf("BTC Finality contract address: %s", addr)
			}
		}
	}

	// Verify all contract addresses were found
	require.NotEmpty(t, btcLightClientAddr, "BTC Light Client contract address not found in events")
	require.NotEmpty(t, btcStakingAddr, "BTC Staking contract address not found in events")
	require.NotEmpty(t, btcFinalityAddr, "BTC Finality contract address not found in events")

	// Set all contract addresses atomically using the new governance message
	contracts := &types.BSNContracts{
		BabylonContract:        babylonAddr,
		BtcLightClientContract: btcLightClientAddr,
		BtcStakingContract:     btcStakingAddr,
		BtcFinalityContract:    btcFinalityAddr,
	}
	setMsg := &types.MsgSetBSNContracts{
		Authority: consumerApp.BabylonKeeper.GetAuthority(),
		Contracts: contracts,
	}
	babylonMsgServer := babylonkeeper.NewMsgServer(consumerApp.BabylonKeeper)
	_, err = babylonMsgServer.SetBSNContracts(ctx, setMsg)
	require.NoError(t, err)

	// Verify that the contracts are set and retrievable via the new unified object
	bsnContracts := consumerApp.BabylonKeeper.GetBSNContracts(ctx)
	require.NotNil(t, bsnContracts)
	require.Equal(t, babylonAddr, bsnContracts.BabylonContract)
	require.Equal(t, btcLightClientAddr, bsnContracts.BtcLightClientContract)
	require.Equal(t, btcStakingAddr, bsnContracts.BtcStakingContract)
	require.Equal(t, btcFinalityAddr, bsnContracts.BtcFinalityContract)

	// Verify that the contracts are instantiated
	require.True(t, wasmKeeper.HasContractInfo(ctx, babylonAccAddr))
	btcLightClientAccAddress, err := sdk.AccAddressFromBech32(btcLightClientAddr)
	require.NoError(t, err)
	require.True(t, wasmKeeper.HasContractInfo(ctx, btcLightClientAccAddress))
	btcStakingAccAddress, err := sdk.AccAddressFromBech32(btcStakingAddr)
	require.NoError(t, err)
	require.True(t, wasmKeeper.HasContractInfo(ctx, btcStakingAccAddress))
	btcFinalityAccAddress, err := sdk.AccAddressFromBech32(btcFinalityAddr)
	require.NoError(t, err)
	require.True(t, wasmKeeper.HasContractInfo(ctx, btcFinalityAccAddress))
}
