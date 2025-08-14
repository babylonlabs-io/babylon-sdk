package keeper_test

import (
	"testing"
	"time"

	"cosmossdk.io/core/header"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"cosmossdk.io/x/evidence"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	"cosmossdk.io/x/upgrade"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/keeper/wasmtesting"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/ibc-go/modules/capability"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v10/modules/core"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"
	"github.com/stretchr/testify/require"

	appparams "github.com/babylonlabs-io/babylon-sdk/demo/app/params"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/keeper"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
)

type encodingConfig struct {
	InterfaceRegistry codectypes.InterfaceRegistry
	Marshaler         codec.Codec
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}

var moduleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distribution.AppModuleBasic{},
	gov.NewAppModuleBasic([]govclient.ProposalHandler{
		paramsclient.ProposalHandler,
	}),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	ibc.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	vesting.AppModuleBasic{},
)

func makeEncodingConfig(_ testing.TB) encodingConfig {
	amino := codec.NewLegacyAmino()
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	std.RegisterInterfaces(interfaceRegistry)
	std.RegisterLegacyAminoCodec(amino)

	moduleBasics.RegisterLegacyAminoCodec(amino)
	moduleBasics.RegisterInterfaces(interfaceRegistry)
	// add wasm types
	wasmtypes.RegisterInterfaces(interfaceRegistry)
	wasmtypes.RegisterLegacyAminoCodec(amino)
	// add babylon types
	types.RegisterInterfaces(interfaceRegistry)
	types.RegisterLegacyAminoCodec(amino)

	return encodingConfig{
		InterfaceRegistry: interfaceRegistry,
		Marshaler:         marshaler,
		TxConfig:          tx.NewTxConfig(marshaler, tx.DefaultSignModes),
		Amino:             amino,
	}
}

func NewTestBabylonKeeperWithStore(
	t testing.TB,
	db dbm.DB,
	stateStore store.CommitMultiStore,
	storeKey *storetypes.KVStoreKey,
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
	wasmKeeper types.WasmKeeper,
	stakingKeeper types.StakingKeeper,
) (*keeper.Keeper, sdk.Context) {
	if storeKey == nil {
		storeKey = storetypes.NewKVStoreKey(types.StoreKey)
	}

	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)

	memKeys := storetypes.NewMemoryStoreKeys(types.MemStoreKey)
	k := keeper.NewKeeper(
		cdc,
		storeKey,
		memKeys[types.MemStoreKey],
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		wasmKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	ctx := sdk.NewContext(
		stateStore,
		cmtproto.Header{
			Time: time.Now().UTC(),
		},
		false,
		log.NewNopLogger(),
	)
	ctx = ctx.WithHeaderInfo(header.Info{})

	return k, ctx
}

func NewTestBabylonKeeper(t testing.TB, bankKeeper types.BankKeeper, accountKeeper types.AccountKeeper, wasmKeeper types.WasmKeeper, stakingKeeper types.StakingKeeper) (*keeper.Keeper, sdk.Context) {
	return NewTestBabylonKeeperWithStoreKey(t, nil, bankKeeper, accountKeeper, wasmKeeper, stakingKeeper)
}

func NewTestBabylonKeeperWithStoreKey(
	t testing.TB,
	storeKey *storetypes.KVStoreKey,
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
	wasmKeeper types.WasmKeeper,
	stakingKeeper types.StakingKeeper,
) (*keeper.Keeper, sdk.Context) {
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewTestLogger(t), storemetrics.NewNoOpMetrics())

	k, ctx := NewTestBabylonKeeperWithStore(t, db, stateStore, storeKey, bankKeeper, accountKeeper, wasmKeeper, stakingKeeper)

	// Initialize params
	if err := k.SetParams(ctx, types.DefaultParams()); err != nil {
		panic(err)
	}

	return k, ctx
}

type TestKeepers struct {
	Ctx              sdk.Context
	StakingKeeper    *stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	BankKeeper       bankkeeper.Keeper
	StoreKey         *storetypes.KVStoreKey
	EncodingConfig   encodingConfig
	BabylonKeeper    *keeper.Keeper
	BabylonMsgServer types.MsgServer
	AccountKeeper    authkeeper.AccountKeeper
	WasmKeeper       *wasmkeeper.Keeper
	WasmMsgServer    wasmtypes.MsgServer
	Faucet           *wasmkeeper.TestFaucet
}

func NewTestKeepers(t testing.TB, opts ...keeper.Option) TestKeepers {
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db, log.NewTestLogger(t), storemetrics.NewNoOpMetrics())

	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distributiontypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, ibcexported.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, ibctransfertypes.StoreKey,
		capabilitytypes.StoreKey, feegrant.StoreKey, authzkeeper.StoreKey,
		wasmtypes.StoreKey, types.StoreKey,
	)
	for _, v := range keys {
		ms.MountStoreWithDB(v, storetypes.StoreTypeIAVL, db)
	}
	memKeys := storetypes.NewMemoryStoreKeys(capabilitytypes.MemStoreKey, types.MemStoreKey)
	for _, v := range memKeys {
		ms.MountStoreWithDB(v, storetypes.StoreTypeMemory, db)
	}
	tkeys := storetypes.NewTransientStoreKeys(paramstypes.TStoreKey)
	for _, v := range tkeys {
		ms.MountStoreWithDB(v, storetypes.StoreTypeTransient, db)
	}
	require.NoError(t, ms.LoadLatestVersion())

	encConfig := makeEncodingConfig(t)
	appCodec := encConfig.Marshaler

	maccPerms := map[string][]string{ // module account permissions
		authtypes.FeeCollectorName:     nil,
		distributiontypes.ModuleName:   nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		types.ModuleName:               {authtypes.Minter, authtypes.Burner},
	}

	// Ensure authority address uses bbnc prefix
	cfg := sdk.GetConfig()
	cfg.SetBech32PrefixForAccount(appparams.Bech32PrefixAccAddr, appparams.Bech32PrefixAccPub)
	authority := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	accountKeeper := authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]), // target store
		authtypes.ProtoBaseAccount,                          // prototype
		maccPerms,
		authcodec.NewBech32Codec(appparams.Bech32PrefixAccAddr),
		appparams.Bech32PrefixAccAddr,
		authority,
	)
	blockedAddrs := make(map[string]bool)
	for acc := range maccPerms {
		blockedAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}
	ctx := sdk.NewContext(ms, cmtproto.Header{
		Height: 1234567,
		Time:   time.Date(2020, time.April, 22, 12, 0, 0, 0, time.UTC),
	}, false, log.NewNopLogger())

	bankKeeper := bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		accountKeeper,
		blockedAddrs,
		authority,
		log.NewNopLogger(),
	)
	require.NoError(t, bankKeeper.SetParams(ctx, banktypes.DefaultParams()))

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[stakingtypes.StoreKey]),
		accountKeeper,
		bankKeeper,
		authority,
		authcodec.NewBech32Codec(appparams.Bech32PrefixValAddr),
		authcodec.NewBech32Codec(appparams.Bech32PrefixConsAddr),
	)
	require.NoError(t, stakingKeeper.SetParams(ctx, stakingtypes.DefaultParams()))

	slashingKeeper := slashingkeeper.NewKeeper(
		appCodec,
		encConfig.Amino,
		runtime.NewKVStoreService(keys[slashingtypes.StoreKey]),
		stakingKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	require.NoError(t, slashingKeeper.SetParams(ctx, slashingtypes.DefaultParams()))

	distKeeper := distributionkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[distributiontypes.StoreKey]),
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(distributiontypes.ModuleName).String(),
	)

	querier := baseapp.NewGRPCQueryRouter()
	querier.SetInterfaceRegistry(encConfig.InterfaceRegistry)
	msgRouter := baseapp.NewMsgServiceRouter()
	msgRouter.SetInterfaceRegistry(encConfig.InterfaceRegistry)

	paramsKeeper := paramskeeper.NewKeeper(
		appCodec,
		encConfig.Amino,
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.TStoreKey],
	)

	upgradeKeeper := upgradekeeper.NewKeeper(
		map[int64]bool{},
		runtime.NewKVStoreService(keys[upgradetypes.StoreKey]),
		appCodec,
		t.TempDir(),
		nil,
		authtypes.NewModuleAddress(upgradetypes.ModuleName).String(),
	)

	ibcKeeper := ibckeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[ibcexported.StoreKey]),
		paramsKeeper.Subspace(ibcexported.ModuleName),
		upgradeKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	cfg = sdk.GetConfig() // Restore original config
	cfg.SetAddressVerifier(wasmtypes.VerifyAddressLen())

	wasmKeeper := wasmkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[wasmtypes.StoreKey]),
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		distributionkeeper.NewQuerier(distKeeper),
		ibcKeeper.ChannelKeeper,
		ibcKeeper.ChannelKeeper,
		wasmtesting.MockIBCTransferKeeper{},
		msgRouter,
		querier,
		t.TempDir(),
		wasmtypes.DefaultNodeConfig(),
		wasmtypes.VMConfig{},
		[]string{"iterator", "staking", "stargate", "cosmwasm_1_1", "cosmwasm_1_2", "cosmwasm_1_3", "cosmwasm_1_4", "cosmwasm_2_0", "virtual_staking"},
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	require.NoError(t, wasmKeeper.SetParams(ctx, wasmtypes.DefaultParams()))
	wasmMsgServer := wasmkeeper.NewMsgServerImpl(&wasmKeeper)

	babylonKeeper := keeper.NewKeeper(
		appCodec,
		keys[types.StoreKey],
		memKeys[types.MemStoreKey],
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		&wasmKeeper,
		authtypes.FeeCollectorName,
		authority,
		opts...,
	)
	require.NoError(t, babylonKeeper.SetParams(ctx, types.DefaultParams()))
	babylonMsgServer := keeper.NewMsgServer(babylonKeeper)

	faucet := wasmkeeper.NewTestFaucet(t, ctx, bankKeeper, minttypes.ModuleName, sdk.NewInt64Coin(sdk.DefaultBondDenom, 1_000_000_000_000))
	return TestKeepers{
		Ctx:              ctx,
		AccountKeeper:    accountKeeper,
		StakingKeeper:    stakingKeeper,
		SlashingKeeper:   slashingKeeper,
		BankKeeper:       bankKeeper,
		StoreKey:         keys[types.StoreKey],
		EncodingConfig:   encConfig,
		BabylonKeeper:    babylonKeeper,
		BabylonMsgServer: babylonMsgServer,
		WasmKeeper:       &wasmKeeper,
		WasmMsgServer:    wasmMsgServer,
		Faucet:           faucet,
	}
}

const (
	TestDataPath                = "../../../tests/testdata"
	BabylonContractCodePath     = TestDataPath + "/babylon_contract.wasm"
	BtcStakingContractCodePath  = TestDataPath + "/btc_staking.wasm"
	BtcFinalityContractCodePath = TestDataPath + "/btc_finality.wasm"
)

func GetGZippedContractCodes() ([]byte, []byte, []byte) {
	babylonContractCode, err := types.GetGZippedContractCode(BabylonContractCodePath)
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

	return babylonContractCode, btcStakingContractCode, btcFinalityContractCode
}
