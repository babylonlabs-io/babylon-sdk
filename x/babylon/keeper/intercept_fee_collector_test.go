package keeper_test

import (
	"math/rand"
	"testing"
	"time"

	sdkmath "cosmossdk.io/math"
	testkeeper "github.com/babylonlabs-io/babylon-sdk/tests/e2e/testutils"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/keeper"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	feeCollectorAcc = authtypes.NewEmptyModuleAccount(authtypes.FeeCollectorName)
	fees            = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100)))
)

func AddRandomSeedsToFuzzer(f *testing.F, num uint) {
	// Seed based on the current time
	r := rand.New(rand.NewSource(time.Now().Unix()))
	var idx uint
	for idx = 0; idx < num; idx++ {
		f.Add(r.Int63())
	}
}

func WithCtxHeight(ctx sdk.Context, height uint64) sdk.Context {
	headerInfo := ctx.HeaderInfo()
	headerInfo.Height = int64(height)
	ctx = ctx.WithHeaderInfo(headerInfo)
	return ctx
}

func FuzzInterceptFeeCollector(f *testing.F) {
	AddRandomSeedsToFuzzer(f, 10)
	f.Fuzz(func(t *testing.T, seed int64) {
		r := rand.New(rand.NewSource(seed))

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// Create BSN contracts with valid bbnc addresses
		bsnContracts := &types.BSNContracts{
			BabylonContract:        "bbnc16t8qwnmdd8wk60enqjugk644ha4xwlqwlkqq70",
			BtcLightClientContract: "bbnc1578akpvpdr8mmr3pd4jw50zpyhv6xucxgdkggr",
			BtcStakingContract:     "bbnc1gev2cfu5fdfupwy6gum9qh9pd75g5f025kh4np",
			BtcFinalityContract:    "bbnc1wg94wzu9a62am7yzztvqqh4k2fqdaf5n9u6k40",
		}

		// Mock bank keeper
		bankKeeper := types.NewMockBankKeeper(ctrl)
		bankKeeper.EXPECT().GetAllBalances(gomock.Any(), feeCollectorAcc.GetAddress()).Return(fees).Times(1)

		// Mock account keeper
		accountKeeper := types.NewMockAccountKeeper(ctrl)
		accountKeeper.EXPECT().GetModuleAccount(gomock.Any(), authtypes.FeeCollectorName).Return(feeCollectorAcc).Times(1)

		babylonKeeper, ctx := testkeeper.BabylonKeeper(t, bankKeeper, accountKeeper, nil, nil)
		height := uint64(r.Intn(1000))
		ctx = WithCtxHeight(ctx, height)

		// Set BSN contracts in the keeper
		err := babylonKeeper.SetBSNContracts(ctx, bsnContracts)
		require.NoError(t, err)

		// Get params and calculate expected fees for BTC staking
		params := babylonKeeper.GetParams(ctx)
		feesForBTCStaking := keeper.GetCoinsPortion(fees, params.BtcStakingPortion)

		// Convert contractAddr string to sdk.AccAddress for expectation
		contractAddrAcc, err := sdk.AccAddressFromBech32(bsnContracts.BtcFinalityContract)
		require.NoError(t, err)

		// Expect the exact portion to be sent to the BTC finality contract
		bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(),
			gomock.Eq(authtypes.FeeCollectorName),
			gomock.Eq(contractAddrAcc),
			gomock.Eq(feesForBTCStaking)).Return(nil).Times(1)

		// Handle coins in fee collector
		err = babylonKeeper.HandleCoinsInFeeCollector(ctx)

		require.NoError(t, err)
	})
}
