package app

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	pruningtypes "cosmossdk.io/store/pruning/types"
	"cosmossdk.io/store/snapshots"
	snapshottypes "cosmossdk.io/store/snapshots/types"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/ed25519"
	tmjson "github.com/cometbft/cometbft/libs/json"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	simsutils "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	appparams "github.com/babylonlabs-io/babylon-sdk/demo/app/params"
)

// SetupOptions defines arguments that are passed into `Simapp` constructor.
type SetupOptions struct {
	Logger             log.Logger
	DB                 *dbm.MemDB
	InvCheckPeriod     uint
	SkipUpgradeHeights map[int64]bool
	AppOpts            servertypes.AppOptions
}

func setup(t testing.TB, chainID string, withGenesis bool, invCheckPeriod uint, opts ...wasmkeeper.Option) (*ConsumerApp, GenesisState) {
	db := dbm.NewMemDB()
	nodeHome := t.TempDir()
	snapshotDir := filepath.Join(nodeHome, "data", "snapshots")

	snapshotDB, err := dbm.NewDB("metadata", dbm.GoLevelDBBackend, snapshotDir)
	require.NoError(t, err)
	t.Cleanup(func() { snapshotDB.Close() })
	snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	require.NoError(t, err)

	appOptions := make(simsutils.AppOptionsMap, 0)
	appOptions[flags.FlagHome] = nodeHome // ensure unique folder
	appOptions[server.FlagTrace] = true
	appOptions[server.FlagInvCheckPeriod] = invCheckPeriod
	app := NewConsumerApp(log.NewNopLogger(), db, nil, true, appOptions, opts, bam.SetChainID(chainID), bam.SetSnapshot(snapshotStore, snapshottypes.SnapshotOptions{KeepRecent: 2}))
	if withGenesis {
		return app, app.DefaultGenesis()
	}
	return app, GenesisState{}
}

// NewAppWithCustomOptions initializes a new ConsumerApp with custom options.
func NewAppWithCustomOptions(t *testing.T, isCheckTx bool, options SetupOptions) *ConsumerApp {
	t.Helper()

	// Ensure address prefixes are set before any test operations
	appparams.SetAddressPrefixes()

	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)
	// create validator set with single validator
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100000000000000))),
	}

	app := NewConsumerApp(options.Logger, options.DB, nil, true, options.AppOpts, emptyWasmOptions)
	genesisState := NewDefaultGenesisState(app.appCodec)
	genesisState, err = GenesisStateWithValSet(app.AppCodec(), genesisState, valSet, []authtypes.GenesisAccount{acc}, balance)
	require.NoError(t, err)

	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		stateBytes, err := tmjson.MarshalIndent(genesisState, "", " ")
		require.NoError(t, err)

		// Initialize the chain
		consensusParams := simsutils.DefaultConsensusParams
		initialHeight := app.LastBlockHeight() + 1
		consensusParams.Abci = &cmtproto.ABCIParams{VoteExtensionsEnableHeight: initialHeight}
		_, err = app.InitChain(
			&abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: consensusParams,
				AppStateBytes:   stateBytes,
				InitialHeight:   initialHeight,
			},
		)
		require.NoError(t, err)
	}

	return app
}

// Setup initializes a new ConsumerApp. A Nop logger is set in ConsumerApp.
func Setup(t *testing.T, opts ...wasmkeeper.Option) *ConsumerApp {
	t.Helper()

	// Ensure address prefixes are set before any test operations
	appparams.SetAddressPrefixes()

	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)

	// create validator set with single validator
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100000000000000))),
	}
	chainID := "testing"
	app := SetupWithGenesisValSet(t, valSet, []authtypes.GenesisAccount{acc}, chainID, opts, balance)

	return app
}

// SetupWithGenesisValSet initializes a new ConsumerApp with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit in the default token of the ConsumerApp from first genesis
// account. A Nop logger is set in ConsumerApp.
func SetupWithGenesisValSet(t *testing.T, valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount, chainID string, opts []wasmkeeper.Option, balances ...banktypes.Balance) *ConsumerApp {
	t.Helper()

	app, genesisState := setup(t, chainID, true, 5, opts...)
	genesisState, err := GenesisStateWithValSet(app.AppCodec(), genesisState, valSet, genAccs, balances...)
	require.NoError(t, err)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)

	// init chain will set the validator set and initialize the genesis accounts
	consensusParams := simsutils.DefaultConsensusParams
	consensusParams.Block.MaxGas = 100 * simsutils.DefaultGenTxGas
	// it is required that the VoteExtensionsEnableHeight > 0 to enable vote extension
	initialHeight := app.LastBlockHeight() + 1
	consensusParams.Abci = &cmtproto.ABCIParams{VoteExtensionsEnableHeight: initialHeight}
	_, err = app.InitChain(&abci.RequestInitChain{
		ChainId:         app.ChainID(),
		Time:            time.Now().UTC(),
		Validators:      []abci.ValidatorUpdate{},
		ConsensusParams: consensusParams,
		InitialHeight:   initialHeight,
		AppStateBytes:   stateBytes,
	})
	require.NoError(t, err)

	_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height:             initialHeight,
		Hash:               app.LastCommitID().Hash,
		NextValidatorsHash: valSet.Hash(),
	})
	require.NoError(t, err)

	return app
}

// SetupWithEmptyStore set up a wasmd app instance with empty DB
func SetupWithEmptyStore(t testing.TB) *ConsumerApp {
	app, _ := setup(t, "testing", false, 0)
	return app
}

// GenesisStateWithSingleValidator initializes GenesisState with a single validator and genesis accounts
// that also act as delegators.
func GenesisStateWithSingleValidator(t *testing.T, app *ConsumerApp) GenesisState {
	t.Helper()

	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)

	// create validator set with single validator
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balances := []banktypes.Balance{
		{
			Address: acc.GetAddress().String(),
			Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100000000000000))),
		},
	}

	genesisState := NewDefaultGenesisState(app.appCodec)
	genesisState, err = GenesisStateWithValSet(app.AppCodec(), genesisState, valSet, []authtypes.GenesisAccount{acc}, balances...)
	require.NoError(t, err)

	return genesisState
}

// AddTestAddrs constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrs(app *ConsumerApp, ctx sdk.Context, accNum int, accAmt math.Int) ([]sdk.AccAddress, error) {
	testAddrs := createRandomAccounts(accNum)

	bondDenom, err := app.StakingKeeper.BondDenom(ctx)
	if err != nil {
		return nil, err
	}
	initCoins := sdk.NewCoins(sdk.NewCoin(bondDenom, accAmt))

	for _, addr := range testAddrs {
		initAccountWithCoins(app, ctx, addr, initCoins)
	}

	return testAddrs, nil
}

// createRandomAccounts is a strategy used by addTestAddrs() in order to generated addresses in random order.
func createRandomAccounts(accNum int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}

func initAccountWithCoins(app *ConsumerApp, ctx sdk.Context, addr sdk.AccAddress, coins sdk.Coins) {
	err := app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, coins)
	if err != nil {
		panic(err)
	}

	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, coins)
	if err != nil {
		panic(err)
	}
}

// ModuleAccountAddrs provides a list of blocked module accounts from configuration in AppConfig
//
// Ported from ConsumerApp
func ModuleAccountAddrs() map[string]bool {
	return BlockedAddresses()
}

var emptyWasmOptions []wasmkeeper.Option

// NewTestNetworkFixture returns a new ConsumerApp AppConstructor for network simulation tests
func NewTestNetworkFixture() network.TestFixture {
	dir, err := os.MkdirTemp("", "simapp")
	if err != nil {
		panic(fmt.Sprintf("failed creating temporary directory: %v", err))
	}
	defer os.RemoveAll(dir)

	app := NewConsumerApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, simsutils.NewAppOptionsWithFlagHome(dir), emptyWasmOptions)
	appCtr := func(val network.ValidatorI) servertypes.Application {
		return NewConsumerApp(val.GetCtx().Logger, dbm.NewMemDB(), nil, true, simsutils.NewAppOptionsWithFlagHome(val.GetCtx().Config.RootDir), emptyWasmOptions, bam.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)), bam.SetMinGasPrices(val.GetAppConfig().MinGasPrices), bam.SetChainID(val.GetCtx().Viper.GetString(flags.FlagChainID)))
	}

	return network.TestFixture{
		AppConstructor: appCtr,
		GenesisState:   NewDefaultGenesisState(app.AppCodec()),
		EncodingConfig: testutil.TestEncodingConfig{
			InterfaceRegistry: app.InterfaceRegistry(),
			Codec:             app.AppCodec(),
			TxConfig:          app.TxConfig(),
			Amino:             app.LegacyAmino(),
		},
	}
}

// SignAndDeliverWithoutCommit signs and delivers a transaction. No commit
func SignAndDeliverWithoutCommit(
	t *testing.T, txCfg client.TxConfig, app *bam.BaseApp, msgs []sdk.Msg, fees sdk.Coins,
	chainID string, accNums, accSeqs []uint64, priv ...cryptotypes.PrivKey,
) (sdk.GasInfo, *sdk.Result, error) {
	tx, err := simsutils.GenSignedMockTx(
		rand.New(rand.NewSource(time.Now().UnixNano())),
		txCfg,
		msgs,
		fees,
		simsutils.DefaultGenTxGas,
		chainID,
		accNums,
		accSeqs,
		priv...,
	)
	require.NoError(t, err)

	// Simulate a sending a transaction and committing a block
	// app.BeginBlock(abci.RequestBeginBlock{Header: header})
	gInfo, res, err := app.SimDeliver(txCfg.TxEncoder(), tx)
	// app.EndBlock(abci.RequestEndBlock{})
	// app.Commit()

	return gInfo, res, err
}

// GenesisStateWithValSet returns a new genesis state with the validator set
// copied from simtestutil with delegation not added to supply
func GenesisStateWithValSet(
	codec codec.Codec,
	genesisState map[string]json.RawMessage,
	valSet *tmtypes.ValidatorSet,
	genAccs []authtypes.GenesisAccount,
	balances ...banktypes.Balance,
) (map[string]json.RawMessage, error) {
	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = codec.MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.DefaultPowerReduction

	for _, val := range valSet.Validators {
		pk, err := cryptocodec.FromCmtPubKeyInterface(val.PubKey)
		if err != nil {
			return nil, fmt.Errorf("failed to convert pubkey: %w", err)
		}

		pkAny, err := codectypes.NewAnyWithValue(pk)
		if err != nil {
			return nil, fmt.Errorf("failed to create new any: %w", err)
		}

		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   math.LegacyOneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
			MinSelfDelegation: math.ZeroInt(),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress().String(), sdk.ValAddress(val.Address).String(), math.LegacyOneDec()))
	}

	// set validators and delegations
	stakingGenesis := stakingtypes.NewGenesisState(stakingtypes.DefaultParams(), validators, delegations)
	genesisState[stakingtypes.ModuleName] = codec.MustMarshalJSON(stakingGenesis)

	signingInfos := make([]slashingtypes.SigningInfo, len(valSet.Validators))
	for i, val := range valSet.Validators {
		signingInfos[i] = slashingtypes.SigningInfo{
			Address:              sdk.ConsAddress(val.Address).String(),
			ValidatorSigningInfo: slashingtypes.ValidatorSigningInfo{},
		}
	}
	slashingParams := slashingtypes.DefaultParams()
	slashingParams.SlashFractionDowntime = math.LegacyNewDec(1).Quo(math.LegacyNewDec(10))
	slashingParams.SlashFractionDoubleSign = math.LegacyNewDec(1).Quo(math.LegacyNewDec(10))
	slashingGenesis := slashingtypes.NewGenesisState(slashingParams, signingInfos, nil)
	genesisState[slashingtypes.ModuleName] = codec.MustMarshalJSON(slashingGenesis)

	// add bonded amount to bonded pool module account
	balances = append(balances, banktypes.Balance{
		Address: authtypes.NewModuleAddress(stakingtypes.BondedPoolName).String(),
		Coins:   sdk.Coins{sdk.NewCoin(sdk.DefaultBondDenom, bondAmt.MulRaw(int64(len(valSet.Validators))))},
	})

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens to total supply
		totalSupply = totalSupply.Add(b.Coins...)
	}

	// update total supply
	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{}, []banktypes.SendEnabled{})
	genesisState[banktypes.ModuleName] = codec.MustMarshalJSON(bankGenesis)

	return genesisState, nil
}
