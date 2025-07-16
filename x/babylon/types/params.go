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
