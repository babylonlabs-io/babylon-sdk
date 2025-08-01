package keeper

import (
	"testing"
	"time"

	"cosmossdk.io/core/header"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	appparams "github.com/babylonlabs-io/babylon/v3/app/params"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/keeper"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
)

func BabylonKeeperWithStore(
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
		appparams.AccGov.String(),
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

func BabylonKeeper(t testing.TB, bankKeeper types.BankKeeper, accountKeeper types.AccountKeeper, wasmKeeper types.WasmKeeper, stakingKeeper types.StakingKeeper) (*keeper.Keeper, sdk.Context) {
	return BabylonKeeperWithStoreKey(t, nil, bankKeeper, accountKeeper, wasmKeeper, stakingKeeper)
}

func BabylonKeeperWithStoreKey(
	t testing.TB,
	storeKey *storetypes.KVStoreKey,
	bankKeeper types.BankKeeper,
	accountKeeper types.AccountKeeper,
	wasmKeeper types.WasmKeeper,
	stakingKeeper types.StakingKeeper,
) (*keeper.Keeper, sdk.Context) {
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewTestLogger(t), storemetrics.NewNoOpMetrics())

	k, ctx := BabylonKeeperWithStore(t, db, stateStore, storeKey, bankKeeper, accountKeeper, wasmKeeper, stakingKeeper)

	// Initialize params
	if err := k.SetParams(ctx, types.DefaultParams()); err != nil {
		panic(err)
	}

	return k, ctx
}
