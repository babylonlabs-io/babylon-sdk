package types

import (
	"context"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BankKeeper interface {
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx context.Context, addr sdk.AccAddress) sdk.Coins
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amounts sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	UndelegateCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
}

// StakingKeeper expected staking keeper.
type StakingKeeper interface {
	BondDenom(ctx context.Context) (string, error)
	StakingTokenSupply(ctx context.Context) (sdkmath.Int, error)
}

// AccountKeeper interface contains functions for getting accounts and the module address
type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(context context.Context, name string) sdk.ModuleAccountI
}

// WasmKeeper abstract wasm keeper
type WasmKeeper interface {
	Sudo(context context.Context, contractAddress sdk.AccAddress, msg []byte) ([]byte, error)
	HasContractInfo(context context.Context, contractAddress sdk.AccAddress) bool
}
