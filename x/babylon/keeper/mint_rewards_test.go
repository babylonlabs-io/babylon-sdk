package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMintBlockRewards(t *testing.T) {
	// myContractAddr := sdk.AccAddress(rand.Bytes(32))
	k := NewTestKeepers(t)

	tests := []struct {
		name          string
		setup         func(ctx sdk.Context)
		recipient     sdk.AccAddress
		inflationRate sdkmath.LegacyDec
		blocksPerYear int64
		expectAmount  sdkmath.Int
		expectErr     bool
	}{
		{
			name: "successful minting",
			setup: func(ctx sdk.Context) {
				// Setup extra initial supply (We already have a supply of 1_000_000_000_000 from the test setup)
				bondDenom, _ := k.StakingKeeper.BondDenom(ctx)
				initialSupply := sdkmath.NewInt(100_000_000_000) // 100_000 tokens
				coins := sdk.NewCoins(sdk.NewCoin(bondDenom, initialSupply))
				err := k.BankKeeper.MintCoins(ctx, types.ModuleName, coins)
				require.NoError(t, err)
			},
			recipient:     sdk.AccAddress("test_recipient"),
			inflationRate: sdkmath.LegacyNewDecWithPrec(5, 2), // 5%
			blocksPerYear: 5256000,                            // ~6 second blocks
			expectAmount:  sdkmath.NewInt(10464),              // Per block reward with these numbers
			expectErr:     false,
		},
		{
			name: "zero inflation rate",
			setup: func(ctx sdk.Context) {
				initialSupply := sdkmath.NewInt(100_000_000_000) // 100_000 tokens
				bondDenom, _ := k.StakingKeeper.BondDenom(ctx)
				coins := sdk.NewCoins(sdk.NewCoin(bondDenom, initialSupply))
				err := k.BankKeeper.MintCoins(ctx, types.ModuleName, coins)
				require.NoError(t, err)
			},
			recipient:     sdk.AccAddress("test_recipient"),
			inflationRate: sdkmath.LegacyZeroDec(),
			blocksPerYear: 5256000,
			expectAmount:  sdkmath.ZeroInt(),
			expectErr:     true,
		},
		{
			name: "negative inflation rate",
			setup: func(ctx sdk.Context) {
				initialSupply := sdkmath.NewInt(100_000_000_000) // 100_000 tokens
				bondDenom, _ := k.StakingKeeper.BondDenom(ctx)
				coins := sdk.NewCoins(sdk.NewCoin(bondDenom, initialSupply))
				err := k.BankKeeper.MintCoins(ctx, types.ModuleName, coins)
				require.NoError(t, err)
			},
			recipient:     sdk.AccAddress("test_recipient"),
			inflationRate: sdkmath.LegacyNewDecWithPrec(-5, 2),
			blocksPerYear: 5256000,
			expectAmount:  sdkmath.ZeroInt(),
			expectErr:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _ := k.Ctx.CacheContext()
			if tc.setup != nil {
				tc.setup(ctx)
			}

			amount, err := k.BabylonKeeper.MintBlockRewards(ctx, tc.recipient, tc.inflationRate, tc.blocksPerYear)

			if tc.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expectAmount, amount)

			// Verify recipient balance
			bondDenom, _ := k.StakingKeeper.BondDenom(ctx)
			balance := k.BankKeeper.GetBalance(ctx, tc.recipient, bondDenom)
			require.Equal(t, amount, balance.Amount)
		})
	}
}
