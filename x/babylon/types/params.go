package types

import (
	"fmt"

	"cosmossdk.io/math"
)

const DefaultGasBeginBlocker = 500_000

// DefaultParams returns default babylon parameters
func DefaultParams() Params {
	return Params{
		MaxGasBeginBlocker: DefaultGasBeginBlocker,
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

	if p.BtcStakingPortion.IsNegative() {
		return fmt.Errorf("BtcStakingPortion %v should not be negative", p.BtcStakingPortion)
	}

	if p.BtcStakingPortion.GT(math.LegacyOneDec()) {
		return fmt.Errorf("BtcStakingPortion %v should not be exceeding 100%", p.BtcStakingPortion)
	}

	return nil
}
