package keeper_test

import (
	"math/rand"
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/babylonlabs-io/babylon/v3/testutil/datagen"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	testkeeper "github.com/babylonlabs-io/babylon-sdk/tests/e2e/testutils"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
)

var (
	feeCollectorAcc = authtypes.NewEmptyModuleAccount(authtypes.FeeCollectorName)
	fees            = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100)))
)

// GetCoinsPortion calculates a portion of coins based on a decimal percentage
func GetCoinsPortion(coinsInt sdk.Coins, portion sdkmath.LegacyDec) sdk.Coins {
	// coins with decimal value
	coins := sdk.NewDecCoinsFromCoins(coinsInt...)
	// portion of coins with decimal values
	portionCoins := coins.MulDecTruncate(portion)
	// truncate back
	portionCoinsInt, _ := portionCoins.TruncateDecimal()
	return portionCoinsInt
}

func FuzzInterceptFeeCollector(f *testing.F) {
	datagen.AddRandomSeedsToFuzzer(f, 10)
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

		keeper, ctx := testkeeper.BabylonKeeper(t, bankKeeper, accountKeeper, nil, nil)
		height := datagen.RandomInt(r, 1000)
		ctx = datagen.WithCtxHeight(ctx, height)

		// Set BSN contracts in the keeper
		err := keeper.SetBSNContracts(ctx, bsnContracts)
		require.NoError(t, err)

		// Get params and calculate expected fees for BTC staking
		params := keeper.GetParams(ctx)
		feesForBTCStaking := GetCoinsPortion(fees, params.BtcStakingPortion)

		// Convert contractAddr string to sdk.AccAddress for expectation
		contractAddrAcc, err := sdk.AccAddressFromBech32(bsnContracts.BtcFinalityContract)
		require.NoError(t, err)

		// Expect the exact portion to be sent to the BTC finality contract
		bankKeeper.EXPECT().SendCoinsFromModuleToAccount(gomock.Any(),
			gomock.Eq(authtypes.FeeCollectorName),
			gomock.Eq(contractAddrAcc),
			gomock.Eq(feesForBTCStaking)).Return(nil).Times(1)

		// Handle coins in fee collector
		err = keeper.HandleCoinsInFeeCollector(ctx)

		require.NoError(t, err)
	})
}
