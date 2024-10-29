package types

import sdkmath "cosmossdk.io/math"

// DefaultParams returns default babylon parameters
func DefaultParams(denom string) Params {
	return Params{
		FinalityInflationRate: sdkmath.LegacyNewDecWithPrec(7, 2),
		BlocksPerYear:         60 * 60 * 24 * 365 / 5, // 5 seconds per block
		MaxGasBeginBlocker:    500_000,
	}
}

// ValidateBasic performs basic validation on babylon parameters.
func (p Params) ValidateBasic() error {
	if p.MaxGasBeginBlocker == 0 {
		return ErrInvalid.Wrap("empty max gas end-blocker setting")
	}
	if p.BlocksPerYear == 0 {
		return ErrInvalid.Wrap("empty blocks per year setting")
	}
	if p.FinalityInflationRate.IsNegative() {
		return ErrInvalid.Wrap("finality inflation rate cannot be negative")
	}
	return nil
}
