package types

import (
	"fmt"

	"cosmossdk.io/math"
)

// DefaultParams returns default babylon parameters
func DefaultParams(denom string) Params {
	return Params{
		MaxGasBeginBlocker: 500_000,
		BtcStakingPortion:  math.LegacyMustNewDecFromStr("0.1"),
	}
}

// ValidateBasic performs basic validation on babylon parameters.
func (p Params) ValidateBasic() error {
	if p.MaxGasBeginBlocker == 0 {
		return ErrInvalid.Wrap("empty max gas end-blocker setting")
	}

	if p.BtcStakingPortion.IsNil() {
		return fmt.Errorf("BtcStakingPortion should not be nil")
	}

	return nil
}
