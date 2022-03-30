package mobility

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/vincepg13/bp-sdk/beyond/x/mobility/tags"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgInitOrder:
			return handleMsgInitOrder(ctx, k, msg)
		case MsgFinalizeOrder:
			return handleMsgFinalizeOrder(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized Msg type: %v", reflect.TypeOf(msg).Name())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgInitOrder(ctx sdk.Context, k Keeper, msg MsgInitOrder) sdk.Result {

	var lastOrderNumber uint64
	lastOrderNumber = k.GetOrderCount(ctx, msg.InitiatorAddress)

	lastOrderNumber++

	k.SetInitOrder(ctx, msg.InitiatorAddress, msg)
	k.SetOrderCount(ctx, msg.InitiatorAddress, lastOrderNumber)

	resTags := sdk.NewTags(
		tags.Action, tags.ActionInitOrder,
		tags.Buyer, []byte(msg.InitiatorAddress.String()),
		tags.Seller, []byte(msg.RecipientAddress.String()),
		tags.OrderNumber, []byte(strconv.FormatUint(lastOrderNumber, 10)),
	)

	return sdk.Result{
		Code: sdk.ABCICodeOK,
		Tags: resTags,
	}
}

func handleMsgFinalizeOrder(ctx sdk.Context, k Keeper, msg MsgFinalizeOrder) sdk.Result {

	//Retrieve last InitOrderNumber from Initiator address
	lastOrderNumber := k.GetOrderCount(ctx, msg.InitiatorAddress)
	//Link initOrder and finalizeOrder in tags

	//TODO: Perform secure payment, move logic from light node
	/*k.PayForCharging(ctx, msg.InitiatorAddress, msg.RecipientAddress, int64(msg.TotalCharge))*/

	resTags := sdk.NewTags(
		tags.Action, tags.ActionFinalizeOrder,
		tags.OrderNumber, []byte(strconv.FormatUint(lastOrderNumber, 10)),
	)

	return sdk.Result{
		Code: sdk.ABCICodeOK,
		Tags: resTags,
	}
}
