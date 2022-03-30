// nolint
package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	ActionInitOrder     = []byte("initOrder")
	ActionFinalizeOrder = []byte("finalizeOrder")

	Action      = sdk.TagAction
	Buyer       = "buyer"
	Seller      = "seller"
	OrderNumber = "orderNumber"
)
