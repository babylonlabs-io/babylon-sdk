package types_test

import (
	sdkmath "cosmossdk.io/math"
	"testing"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestValidateGenesis(t *testing.T) {
	specs := map[string]struct {
		state  types.GenesisState
		expErr bool
	}{
		"default params": {
			state:  *types.DefaultGenesisState(sdk.DefaultBondDenom),
			expErr: false,
		},
		"custom param, should pass": {
			state: types.GenesisState{
				Params: types.Params{
					MaxGasBeginBlocker:    600_000,
					BlocksPerYear:         6_000_000,
					FinalityInflationRate: sdkmath.LegacyNewDecWithPrec(1, 1),
				},
			},
			expErr: false,
		},
		"custom small value param, should pass": {
			state: types.GenesisState{
				Params: types.Params{
					MaxGasBeginBlocker:    10_000,
					BlocksPerYear:         6_000,
					FinalityInflationRate: sdkmath.LegacyNewDecWithPrec(1, 10),
				},
			},
			expErr: false,
		},
		"invalid max gas length, should fail": {
			state: types.GenesisState{
				Params: types.Params{
					MaxGasBeginBlocker: 0,
				},
			},
			expErr: true,
		},
		"invalid blocks per year, should fail": {
			state: types.GenesisState{
				Params: types.Params{
					MaxGasBeginBlocker: 1,
					BlocksPerYear:      0,
				},
			},
			expErr: true,
		},
		"invalid finality inflation rate, should fail": {
			state: types.GenesisState{
				Params: types.Params{
					MaxGasBeginBlocker:    1,
					BlocksPerYear:         1,
					FinalityInflationRate: sdkmath.LegacyNewDec(-1),
				},
			},
			expErr: true,
		},
	}
	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			err := types.ValidateGenesis(&spec.state)
			if spec.expErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
