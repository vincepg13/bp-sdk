package mobility

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Mobility errors reserve 300 ~ 399.
const (
	DefaultCodespace      sdk.CodespaceType = 4
	CodeNoChargeAmount    sdk.CodeType      = 398
	CodeNoOrderNumber     sdk.CodeType      = 399
	CodeEmptyEnergyAmount sdk.CodeType      = 400
)

// ErrNoEstimatedEnergyAmount
func ErrNoEstimatedEnergyAmount(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyEnergyAmount, fmt.Sprintf("Empty energy amount"))
}
func ErrNoOrderNumber(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNoOrderNumber, fmt.Sprintf("Error retrieving order number"))
}

func ErrNoChargeAmountProvided() sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeNoChargeAmount, fmt.Sprintf("Provide total charge amount"))
}
