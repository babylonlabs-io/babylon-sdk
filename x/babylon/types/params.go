package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultParams returns default babylon parameters
func DefaultParams(denom string) Params {
	return Params{
		MaxGasBeginBlocker: 500_000,
	}
}

// ValidateBasic performs basic validation on babylon parameters.
func (p Params) ValidateBasic() error {
	if p.MaxGasBeginBlocker == 0 {
		return ErrInvalid.Wrap("empty max gas end-blocker setting")
	}
	return nil
}

func (p Params) GetContractAddresses() (sdk.AccAddress, sdk.AccAddress, sdk.AccAddress, error) {
	if !p.IsCodeStored() {
		return nil, nil, nil, errors.New("contracts are not instantiated")
	}

	babylonAddr, err := sdk.AccAddressFromBech32(p.BabylonContractAddress)
	if err != nil {
		return nil, nil, nil, err
	}
	btcStakingAddr, err := sdk.AccAddressFromBech32(p.BtcStakingContractAddress)
	if err != nil {
		return nil, nil, nil, err
	}
	btcFinalityAddr, err := sdk.AccAddressFromBech32(p.BtcFinalityContractAddress)
	if err != nil {
		return nil, nil, nil, err
	}

	return babylonAddr, btcStakingAddr, btcFinalityAddr, nil
}

func (p Params) IsCodeStored() bool {
	return p.BabylonContractCodeId != 0 &&
		p.BtcStakingContractCodeId != 0 &&
		p.BtcFinalityContractCodeId != 0
}

func (p Params) IsContractInstantiated() bool {
	return p.IsCodeStored() &&
		len(p.BabylonContractAddress) > 0 &&
		len(p.BtcStakingContractAddress) > 0 &&
		len(p.BtcFinalityContractAddress) > 0
}
