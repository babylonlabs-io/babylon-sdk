package types

import (
	"testing"

	"github.com/CosmWasm/wasmd/tests/ibctesting"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cometbft/cometbft/types"
	"github.com/stretchr/testify/require"

	"github.com/babylonlabs-io/babylon-sdk/demo/app"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// NewIBCCoordinator initializes Coordinator with N bcd TestChain instances
func NewIBCCoordinator(t *testing.T, opts ...[]wasmkeeper.Option) *ibctesting.Coordinator {
	return ibctesting.NewCoordinatorX(t, 2,
		func(
			t *testing.T,
			valSet *types.ValidatorSet,
			genAccs []authtypes.GenesisAccount,
			chainID string,
			opts []wasmkeeper.Option,
			balances ...banktypes.Balance,
		) ibctesting.ChainApp {
			return app.SetupWithGenesisValSet(t, valSet, genAccs, chainID, opts, balances...)
		},
		opts...,
	)
}

func InstantiateContract(t *testing.T, chain *ibctesting.TestChain, codeID uint64, initMsg []byte, funds ...sdk.Coin) sdk.AccAddress {
	instantiateMsg := &wasmtypes.MsgInstantiateContract{
		Sender: chain.SenderAccount.GetAddress().String(),
		Admin:  chain.SenderAccount.GetAddress().String(),
		CodeID: codeID,
		Label:  "ibc-test",
		Msg:    initMsg,
		Funds:  funds,
	}

	r, err := chain.SendMsgs(instantiateMsg)
	require.NoError(t, err)
	require.Zero(t, r.Code)
	require.NotEmpty(t, r.Data)

	// ensure there is only 1 contract under this code ID
	ctx := chain.GetContext()
	contractAddrs := []sdk.AccAddress{}
	chain.App.GetWasmKeeper().IterateContractsByCode(ctx, codeID, func(address sdk.AccAddress) bool {
		contractAddrs = append(contractAddrs, address)
		return false // keep iterating
	})
	require.Len(t, contractAddrs, 1)

	return contractAddrs[0]
}
