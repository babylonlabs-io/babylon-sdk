package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Option is an extension point to instantiate keeper with non default values
type Option interface {
	apply(*Keeper)
}

type Keeper struct {
	storeKey storetypes.StoreKey
	memKey   storetypes.StoreKey
	cdc      codec.Codec
	bank     types.BankKeeper
	Staking  types.StakingKeeper
	wasm     types.WasmKeeper
	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper constructor with vanilla sdk keepers
func NewKeeper(
	cdc codec.Codec,
	storeKey storetypes.StoreKey,
	memoryStoreKey storetypes.StoreKey,
	bank types.BankKeeper,
	staking types.StakingKeeper,
	wasm types.WasmKeeper,
	authority string,
	opts ...Option,
) *Keeper {
	k := &Keeper{
		storeKey:  storeKey,
		memKey:    memoryStoreKey,
		cdc:       cdc,
		bank:      bank,
		Staking:   staking,
		wasm:      wasm,
		authority: authority,
	}
	for _, o := range opts {
		o.apply(k)
	}

	return k
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetTest(ctx sdk.Context, actor sdk.AccAddress) string {
	return "placeholder"
}
