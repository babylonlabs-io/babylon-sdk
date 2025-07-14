package keeper

import (
	"context"
	"fmt"

	corestoretypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
)

// Option is an extension point to instantiate keeper with non default values
type Option interface {
	apply(*Keeper)
}

type Keeper struct {
	cdc           codec.Codec
	bankKeeper    types.BankKeeper
	Staking       types.StakingKeeper
	wasm          *wasmkeeper.Keeper
	accountKeeper types.AccountKeeper
	storeService  corestoretypes.KVStoreService
	// name of the FeeCollector ModuleAccount
	feeCollectorName string
	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper constructor with vanilla sdk keepers
func NewKeeper(
	cdc codec.Codec,
	storeService corestoretypes.KVStoreService,
	accountKeeper types.AccountKeeper,
	bank types.BankKeeper,
	staking types.StakingKeeper,
	wasm *wasmkeeper.Keeper,
	feeCollectorName string,
	authority string,
	opts ...Option,
) *Keeper {
	k := &Keeper{
		cdc:              cdc,
		storeService:     storeService,
		bankKeeper:       bank,
		accountKeeper:    accountKeeper,
		Staking:          staking,
		wasm:             wasm,
		feeCollectorName: feeCollectorName,
		authority:        authority,
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

func (k Keeper) Logger(goCtx context.Context) log.Logger {
	ctx := sdk.UnwrapSDKContext(goCtx)
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetTest(ctx sdk.Context, actor sdk.AccAddress) string {
	return "placeholder"
}
