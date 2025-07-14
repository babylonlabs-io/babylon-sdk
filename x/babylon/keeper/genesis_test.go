package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/babylonlabs-io/babylon-sdk/x/babylon/types"
)

func TestInitGenesis(t *testing.T) {
	testAddr1 := "cosmos10ak4gg0cy6puxjed9sj58pwek7rms0cqmdma2w"
	testAddr2 := "cosmos1vsegxma779druznaruxsr0emnfjgdpfrvtwr8v"
	testAddr3 := "cosmos1336gfnzgwjr43a5mj7q756y84hc7krppn2uq3y"
	testAddr4 := "cosmos169g8wamxt26v3r86pm0gyxykganntgeldaaztc"

	specs := map[string]struct {
		state  types.GenesisState
		expErr bool
	}{
		"custom param, should pass": {
			state: types.GenesisState{
				Params: types.Params{
					MaxGasBeginBlocker: 600_000,
				},
				BsnContracts: &types.BSNContracts{
					BabylonContract:        testAddr1,
					BtcLightClientContract: testAddr2,
					BtcStakingContract:     testAddr3,
					BtcFinalityContract:    testAddr4,
				},
			},
			expErr: false,
		},
		"custom small value param, should pass": {
			state: types.GenesisState{
				Params: types.Params{
					MaxGasBeginBlocker: 10_000,
				},
				BsnContracts: &types.BSNContracts{
					BabylonContract:        testAddr1,
					BtcLightClientContract: testAddr2,
					BtcStakingContract:     testAddr3,
					BtcFinalityContract:    testAddr4,
				},
			},
			expErr: false,
		},
	}
	specs["invalid babylon contract address, should panic"] = struct {
		state  types.GenesisState
		expErr bool
	}{
		state: types.GenesisState{
			Params: types.Params{
				MaxGasBeginBlocker: 600_000,
			},
			BsnContracts: &types.BSNContracts{
				BabylonContract:        "not-a-valid-address",
				BtcLightClientContract: testAddr2,
				BtcStakingContract:     testAddr3,
				BtcFinalityContract:    testAddr4,
			},
		},
		expErr: true,
	}
	specs["bsn contracts being nil, should pass"] = struct {
		state  types.GenesisState
		expErr bool
	}{
		state: types.GenesisState{
			Params: types.Params{
				MaxGasBeginBlocker: 600_000,
			},
			BsnContracts: nil,
		},
		expErr: false,
	}
	specs["bsn contract addresses being nil, should pass"] = struct {
		state  types.GenesisState
		expErr bool
	}{
		state: types.GenesisState{
			Params: types.Params{
				MaxGasBeginBlocker: 600_000,
			},
			BsnContracts: &types.BSNContracts{},
		},
		expErr: false,
	}

	for name, spec := range specs {
		t.Run(name, func(t *testing.T) {
			keepers := NewTestKeepers(t)
			k := keepers.BabylonKeeper

			if spec.expErr {
				assert.Panics(t, func() {
					k.InitGenesis(keepers.Ctx, spec.state)
				})
				return
			}
			k.InitGenesis(keepers.Ctx, spec.state)

			p := k.GetParams(keepers.Ctx)
			assert.Equal(t, spec.state.Params.MaxGasBeginBlocker, p.MaxGasBeginBlocker)
			// Check contract addresses
			contracts := k.GetBSNContracts(keepers.Ctx)
			if spec.state.BsnContracts != nil && spec.state.BsnContracts.IsSet() {
				assert.Equal(t, testAddr1, contracts.BabylonContract)
				assert.Equal(t, testAddr2, contracts.BtcLightClientContract)
				assert.Equal(t, testAddr3, contracts.BtcStakingContract)
				assert.Equal(t, testAddr4, contracts.BtcFinalityContract)
			}
		})
	}
}

func TestExportGenesis(t *testing.T) {
	keepers := NewTestKeepers(t)
	k := keepers.BabylonKeeper
	params := types.DefaultParams(sdk.DefaultBondDenom)
	testAddr1 := "cosmos10ak4gg0cy6puxjed9sj58pwek7rms0cqmdma2w"
	testAddr2 := "cosmos1vsegxma779druznaruxsr0emnfjgdpfrvtwr8v"
	testAddr3 := "cosmos1336gfnzgwjr43a5mj7q756y84hc7krppn2uq3y"
	testAddr4 := "cosmos169g8wamxt26v3r86pm0gyxykganntgeldaaztc"
	bsnContracts := types.BSNContracts{
		BabylonContract:        testAddr1,
		BtcLightClientContract: testAddr2,
		BtcStakingContract:     testAddr3,
		BtcFinalityContract:    testAddr4,
	}

	err := k.SetParams(keepers.Ctx, params)
	require.NoError(t, err)
	// Set contract addresses
	err = k.SetBSNContracts(keepers.Ctx, &bsnContracts)
	require.NoError(t, err)

	exported := k.ExportGenesis(keepers.Ctx)
	assert.Equal(t, params.MaxGasBeginBlocker, exported.Params.MaxGasBeginBlocker)
	assert.Equal(t, testAddr1, exported.BsnContracts.BabylonContract)
	assert.Equal(t, testAddr2, exported.BsnContracts.BtcLightClientContract)
	assert.Equal(t, testAddr3, exported.BsnContracts.BtcStakingContract)
	assert.Equal(t, testAddr4, exported.BsnContracts.BtcFinalityContract)
}
func TestExportGenesisEmptyContracts(t *testing.T) {
	keepers := NewTestKeepers(t)
	k := keepers.BabylonKeeper
	params := types.DefaultParams(sdk.DefaultBondDenom)

	err := k.SetParams(keepers.Ctx, params)
	require.NoError(t, err)

	exported := k.ExportGenesis(keepers.Ctx)
	assert.Equal(t, params.MaxGasBeginBlocker, exported.Params.MaxGasBeginBlocker)
	assert.Nil(t, exported.BsnContracts)
}
