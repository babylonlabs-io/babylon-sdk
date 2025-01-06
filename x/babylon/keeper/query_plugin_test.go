package keeper_test

import (
	"encoding/json"
	"testing"

	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/contract"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/keeper"
	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
	"github.com/cometbft/cometbft/libs/rand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestChainedCustomQuerier(t *testing.T) {
	myContractAddr := sdk.AccAddress(rand.Bytes(32))
	keepers := NewTestKeepers(t)

	defaultParams := contract.ParamsResponse{
		BabylonContractCodeId:      0,
		BtcStakingContractCodeId:   0,
		BtcFinalityContractCodeId:  0,
		BabylonContractAddress:     "",
		BtcStakingContractAddress:  "",
		BtcFinalityContractAddress: "",
		MaxGasBeginBlocker:         500_000,
	}
	expData, err := json.Marshal(defaultParams)
	require.NoError(t, err)

	specs := map[string]struct {
		src           wasmvmtypes.QueryRequest
		viewKeeper    keeper.ViewKeeper
		expData       []byte
		expErr        bool
		expNextCalled bool
	}{
		"non custom query": {
			src: wasmvmtypes.QueryRequest{
				Bank: &wasmvmtypes.BankQuery{},
			},
			viewKeeper:    keepers.BabylonKeeper,
			expNextCalled: true,
		},
		"unexpected babylon query": {
			src: wasmvmtypes.QueryRequest{
				Custom: []byte(`{"foo":{}}`),
			},
			viewKeeper:    keepers.BabylonKeeper,
			expNextCalled: true,
		},
		"expected babylon query": {
			src: wasmvmtypes.QueryRequest{
				Custom: []byte(`{"params":{}}`),
			},
			viewKeeper:    keepers.BabylonKeeper,
			expNextCalled: false,
			expData:       expData,
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			var nextCalled bool
			next := keeper.QueryHandlerFn(func(ctx sdk.Context, caller sdk.AccAddress, request wasmvmtypes.QueryRequest) ([]byte, error) {
				nextCalled = true
				return nil, nil
			})

			ctx, _ := keepers.Ctx.CacheContext()
			gotData, gotErr := keeper.ChainedCustomQuerier(spec.viewKeeper, next).HandleQuery(ctx, myContractAddr, spec.src)
			if spec.expErr {
				require.Error(t, gotErr)
				return
			}
			require.NoError(t, gotErr)
			assert.Equal(t, spec.expData, gotData, string(gotData))
			assert.Equal(t, spec.expNextCalled, nextCalled)
		})
	}
}

var _ keeper.ViewKeeper = &MockViewKeeper{}

type MockViewKeeper struct {
	GetParamsFn func(ctx sdk.Context) types.Params
}

func (m MockViewKeeper) GetParams(ctx sdk.Context) types.Params {
	if m.GetParamsFn == nil {
		panic("not expected to be called")
	}
	return m.GetParamsFn(ctx)
}
