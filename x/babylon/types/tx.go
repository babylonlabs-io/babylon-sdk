package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ValidateBasic validate basic constraints
func (msg MsgSetBSNContracts) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
	}
	if msg.Contracts == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "contracts must be set")
	}
	if err := msg.Contracts.ValidateBasic(); err != nil {
		return err
	}
	return nil
}
