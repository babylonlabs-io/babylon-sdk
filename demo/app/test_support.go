package app

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v10/modules/core/keeper"

	"github.com/cosmos/cosmos-sdk/baseapp"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

func (app *ConsumerApp) GetIBCKeeper() *ibckeeper.Keeper {
	return app.IBCKeeper
}

func (app *ConsumerApp) GetBaseApp() *baseapp.BaseApp {
	return app.BaseApp
}

func (app *ConsumerApp) GetBankKeeper() bankkeeper.Keeper {
	return app.BankKeeper
}

func (app *ConsumerApp) GetStakingKeeper() *stakingkeeper.Keeper {
	return app.StakingKeeper
}

func (app *ConsumerApp) GetAccountKeeper() authkeeper.AccountKeeper {
	return app.AccountKeeper
}

func (app *ConsumerApp) GetWasmKeeper() wasmkeeper.Keeper {
	return app.WasmKeeper
}
