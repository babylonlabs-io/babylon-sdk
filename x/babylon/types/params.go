package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"fmt"

	"cosmossdk.io/math"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters
func DefaultParams(denom string) Params {
	return Params{
		MaxGasBeginBlocker: 500_000,
		BtcStakingPortion:  math.LegacyNewDecWithPrec(6, 1), // 6 * 10^{-1} = 0.6
	}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// TotalPortion calculates the sum of portions of all stakeholders
func (p *Params) TotalPortion() math.LegacyDec {
	sum := p.BtcStakingPortion
	return sum
}

// BTCStakingPortion calculates the sum of portions of all BTC staking stakeholders
func (p *Params) BTCStakingPortion() math.LegacyDec {
	return p.BtcStakingPortion
}

// Validate validates the set of params
func (p Params) Validate() error {
	if p.BtcStakingPortion.IsNil() {
		return fmt.Errorf("BtcStakingPortion should not be nil")
	}

	// sum of all portions should be less than 1
	if p.TotalPortion().GTE(math.LegacyOneDec()) {
		return fmt.Errorf("sum of all portions should be less than 1")
	}

	if p.MaxGasBeginBlocker == 0 {
		return fmt.Errorf("empty max gas end-blocker setting")
	}

	return nil
}

func (p Params) GetContractAddresses() (sdk.AccAddress, sdk.AccAddress, sdk.AccAddress, sdk.AccAddress, error) {
	if !p.IsCodeStored() {
		return nil, nil, nil, nil, errors.New("contracts are not instantiated")
	}

	babylonAddr, err := sdk.AccAddressFromBech32(p.BabylonContractAddress)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	btcLightClientAddr, err := sdk.AccAddressFromBech32(p.BtcLightClientContractAddress)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	btcStakingAddr, err := sdk.AccAddressFromBech32(p.BtcStakingContractAddress)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	btcFinalityAddr, err := sdk.AccAddressFromBech32(p.BtcFinalityContractAddress)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return babylonAddr, btcLightClientAddr, btcStakingAddr, btcFinalityAddr, nil
}

func (p Params) IsCodeStored() bool {
	return p.BabylonContractCodeId != 0 &&
		p.BtcLightClientContractCodeId != 0 &&
		p.BtcStakingContractCodeId != 0 &&
		p.BtcFinalityContractCodeId != 0
}

func (p Params) IsContractInstantiated() bool {
	return p.IsCodeStored() &&
		len(p.BabylonContractAddress) > 0 &&
		len(p.BtcLightClientContractAddress) > 0 &&
		len(p.BtcStakingContractAddress) > 0 &&
		len(p.BtcFinalityContractAddress) > 0
}
