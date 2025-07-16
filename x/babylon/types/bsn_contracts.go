package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Equal compares two BSNContracts for equality.
func (c *BSNContracts) Equal(other *BSNContracts) bool {
	if c == nil && other == nil {
		return true
	}
	if c == nil || other == nil {
		return false
	}
	return c.BabylonContract == other.BabylonContract &&
		c.BtcLightClientContract == other.BtcLightClientContract &&
		c.BtcStakingContract == other.BtcStakingContract &&
		c.BtcFinalityContract == other.BtcFinalityContract
}

// ValidateBasic validates the BSNContracts object
func (c *BSNContracts) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(c.BabylonContract); err != nil {
		return errorsmod.Wrap(err, "babylon contract")
	}
	if _, err := sdk.AccAddressFromBech32(c.BtcLightClientContract); err != nil {
		return errorsmod.Wrap(err, "btc light client contract")
	}
	if _, err := sdk.AccAddressFromBech32(c.BtcStakingContract); err != nil {
		return errorsmod.Wrap(err, "btc staking contract")
	}
	if _, err := sdk.AccAddressFromBech32(c.BtcFinalityContract); err != nil {
		return errorsmod.Wrap(err, "btc finality contract")
	}
	return nil
}

func (c *BSNContracts) IsSet() bool {
	return c.BabylonContract != "" &&
		c.BtcFinalityContract != "" &&
		c.BtcLightClientContract != "" &&
		c.BtcStakingContract != ""
}
