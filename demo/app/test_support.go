package app

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	"github.com/cosmos/cosmos-sdk/baseapp"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
)

func (app *BSNApp) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper
}

func (app *BSNApp) GetScopedIBCKeeper() capabilitykeeper.ScopedKeeper {
	return app.ScopedIBCKeeper
}

func (app *BSNApp) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

func (app *BSNApp) GetBankKeeper() bankkeeper.Keeper {
	return app.BankKeeper
}

func (app *BSNApp) GetStakingKeeper() *stakingkeeper.Keeper {
	return app.StakingKeeper
}

func (app *BSNApp) GetAccountKeeper() authkeeper.AccountKeeper {
	return app.AccountKeeper
}

func (app *BSNApp) GetWasmKeeper() wasmkeeper.Keeper {
	return app.WasmKeeper
}
