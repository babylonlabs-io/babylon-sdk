package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
)

func TestInitGenesis(t *testing.T) {
	specs := map[string]struct {
		state  types.GenesisState
		expErr bool
	}{
		"custom param, should pass": {
			state: types.GenesisState{
				Params: types.Params{
					MaxGasBeginBlocker: 600_000,
				},
			},
			expErr: false,
		},
		"custom small value param, should pass": {
			state: types.GenesisState{
				Params: types.Params{
					MaxGasBeginBlocker: 10_000,
				},
			},
			expErr: false,
		},
	}

	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			keepers := NewTestKeepers(t)
			k := keepers.BabylonKeeper

			k.InitGenesis(keepers.Ctx, spec.state)

			p := k.GetParams(keepers.Ctx)
			assert.Equal(t, spec.state.Params.MaxGasBeginBlocker, p.MaxGasBeginBlocker)
		})
	}
}

func TestExportGenesis(t *testing.T) {
	keepers := NewTestKeepers(t)
	k := keepers.BabylonKeeper
	params := types.DefaultParams(sdk.DefaultBondDenom)

	err := k.SetParams(keepers.Ctx, params)
	require.NoError(t, err)

	exported := k.ExportGenesis(keepers.Ctx)
	assert.Equal(t, params.MaxGasBeginBlocker, exported.Params.MaxGasBeginBlocker)
}
