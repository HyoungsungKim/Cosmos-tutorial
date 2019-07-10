package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RouterKey = ModuleName
type MsgSetName struct {
	Name string 			'json:"name"'
	Value string			'json:"value"'
	Owner sdk.AccAddress 	'json:"owner"'
}

func NewMsgSetName(name string, value string, owen sdk.AccAddress) MsgSetName {
	return MsgSetName {
		Name: name,
		Value: value,
		Owner: owner,
	}
}

func (msg MsgSetName) Route() string {return RouterKey}

func (msg MsgSetName) Type() string {return "set_name"}

func (msg MsgSetName) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvakudAddress(msg.Owner.String())
	}
	if len(msg.Name) == 0 || len(msg.Value) == 0 {
		return sdk.ErrUnknownRequest("Name and/or Value cannot be empty")
	}
	return nil
}

func (msg MsgSetName) GetSignBytes() []byte {
	return sdk.MustSortJson(ModuleCdc.MustMarshalJson(msg))
}

func (msg MsgSetName)GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

