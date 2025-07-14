package types_test

import (
	"testing"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestValidateGenesis(t *testing.T) {
	validAddr := "cosmos10ak4gg0cy6puxjed9sj58pwek7rms0cqmdma2w"
	invalidAddr := "test-invalid-addr"
	specs := map[string]struct {
		state  types.GenesisState
		expErr bool
	}{
		"default params": {
			state: types.GenesisState{
				Params: types.DefaultParams(sdk.DefaultBondDenom),
				BsnContracts: &types.BSNContracts{
					BabylonContract:        validAddr,
					BtcLightClientContract: validAddr,
					BtcStakingContract:     validAddr,
					BtcFinalityContract:    validAddr,
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
		"invalid babylon contract address, should fail": {
			state: types.GenesisState{
				Params: types.DefaultParams(sdk.DefaultBondDenom),
				BsnContracts: &types.BSNContracts{
					BabylonContract:        invalidAddr,
					BtcLightClientContract: validAddr,
					BtcStakingContract:     validAddr,
					BtcFinalityContract:    validAddr,
				},
			},
			expErr: true,
		},
		"invalid btc light client contract address, should fail": {
			state: types.GenesisState{
				Params: types.DefaultParams(sdk.DefaultBondDenom),
				BsnContracts: &types.BSNContracts{
					BabylonContract:        validAddr,
					BtcLightClientContract: invalidAddr,
					BtcStakingContract:     validAddr,
					BtcFinalityContract:    validAddr,
				},
			},
			expErr: true,
		},
		"invalid btc staking contract address, should fail": {
			state: types.GenesisState{
				Params: types.DefaultParams(sdk.DefaultBondDenom),
				BsnContracts: &types.BSNContracts{
					BabylonContract:        validAddr,
					BtcLightClientContract: validAddr,
					BtcStakingContract:     invalidAddr,
					BtcFinalityContract:    validAddr,
				},
			},
			expErr: true,
		},
		"invalid btc finality contract address, should fail": {
			state: types.GenesisState{
				Params: types.DefaultParams(sdk.DefaultBondDenom),
				BsnContracts: &types.BSNContracts{
					BabylonContract:        validAddr,
					BtcLightClientContract: validAddr,
					BtcStakingContract:     validAddr,
					BtcFinalityContract:    invalidAddr,
				},
			},
			expErr: true,
		},
		"invalid max cap coin denom, should fail": {
			state: types.GenesisState{
				Params: types.Params{
					MaxGasBeginBlocker: 0,
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
