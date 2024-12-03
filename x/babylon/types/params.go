package types

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
