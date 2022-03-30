package mobility

import (
	"bytes"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgInitOrder is a Msg type for initiating an Order when buying conditions are agreed on.
// Extend it to add additional fields (Order conditions, etc)
type MsgInitOrder struct {
	InitiatorAddress sdk.AccAddress
	RecipientAddress sdk.AccAddress
	AgreedPrice      uint64
	EstimatedCharge  uint64
}

// Construct new NewMsgInitOrder.
func NewMsgInitOrder(initiatorAddress sdk.AccAddress, recipientAddress sdk.AccAddress, price uint64, estimatedCharge uint64) MsgInitOrder {
	return MsgInitOrder{
		InitiatorAddress: initiatorAddress,
		RecipientAddress: recipientAddress,
		AgreedPrice:      price,
		EstimatedCharge:  estimatedCharge,
	}
}

// enforce the msg type at compile time
var _ sdk.Msg = MsgInitOrder{}

//nolint
func (msg MsgInitOrder) Type() string                 { return "mobility" }
func (msg MsgInitOrder) Route() string                { return "order" }
func (msg MsgInitOrder) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.InitiatorAddress} }
func (msg MsgInitOrder) String() string {
	return fmt.Sprintf("MsgInitOrder{InitiatorAddress: %v, Price: %v, Amount: %v}", msg.InitiatorAddress, msg.AgreedPrice, msg.EstimatedCharge)
}

// validate MsgInitOrder
func (msg MsgInitOrder) ValidateBasic() sdk.Error {
	if len(msg.InitiatorAddress) == 0 {
		return sdk.ErrUnknownAddress(msg.InitiatorAddress.String()).TraceSDK("")
	}

	if len(msg.RecipientAddress) == 0 {
		return sdk.ErrUnknownAddress(msg.RecipientAddress.String()).TraceSDK("")
	}

	if bytes.Equal(msg.InitiatorAddress, msg.RecipientAddress) {
		return sdk.ErrInvalidAddress("Initiator and recipient have the same address")
	}
	return nil
}

// GetSignBytes returns the canonical byte representation of the Msg.
func (msg MsgInitOrder) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bz
}

//_______________________________________________________________________

type MsgFinalizeOrder struct {
	InitiatorAddress sdk.AccAddress
	RecipientAddress sdk.AccAddress
	TotalAmount      uint64
	TotalCharge      uint64
}

// Construct new NewMsgInitOrder.
func NewMsgFinalizeOrder(initiatorAddress sdk.AccAddress, recipientAddress sdk.AccAddress, totalAmount uint64, totalCharge uint64) MsgFinalizeOrder {
	return MsgFinalizeOrder{
		InitiatorAddress: initiatorAddress,
		RecipientAddress: recipientAddress,
		TotalAmount:      totalAmount,
		TotalCharge:      totalCharge,
	}
}

var _ sdk.Msg = MsgFinalizeOrder{}

//nolint
func (msg MsgFinalizeOrder) Type() string  { return "mobility" }
func (msg MsgFinalizeOrder) Route() string { return "order" }
func (msg MsgFinalizeOrder) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.InitiatorAddress}
}
func (msg MsgFinalizeOrder) String() string {
	return fmt.Sprintf("MsgFinalizeOrder{InitiatorAddress: %v, TotalAmount: %v, TotalCharge: %v}", msg.InitiatorAddress, msg.TotalAmount, msg.TotalCharge)
}

// validate MsgFinalizeOrder
func (msg MsgFinalizeOrder) ValidateBasic() sdk.Error {
	if len(msg.InitiatorAddress) == 0 {
		return sdk.ErrUnknownAddress(msg.InitiatorAddress.String()).TraceSDK("")
	}

	if len(msg.RecipientAddress) == 0 {
		return sdk.ErrUnknownAddress(msg.RecipientAddress.String()).TraceSDK("")
	}

	if bytes.Equal(msg.InitiatorAddress, msg.RecipientAddress) {
		return sdk.ErrInvalidAddress("Initiator and recipient have the same address")
	}

	if msg.TotalCharge <= 0 {
		return ErrNoChargeAmountProvided()
	}
	return nil
}

// GetSignBytes returns the canonical byte representation of the Msg.
func (msg MsgFinalizeOrder) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bz
}
