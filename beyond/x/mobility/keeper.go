package mobility

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

// Keeper
type Keeper struct {
	ck        bank.Keeper
	storeKey  sdk.StoreKey // The (unexposed) key used to access the store from the Context.
	cdc       *codec.Codec
	codespace sdk.CodespaceType
}

func NewKeeper(key sdk.StoreKey, coinKeeper bank.Keeper, codespace sdk.CodespaceType) Keeper {
	cdc := codec.New()
	codec.RegisterCrypto(cdc)
	return Keeper{
		storeKey:  key,
		cdc:       cdc,
		ck:        coinKeeper,
		codespace: codespace,
	}
}

// Key to knowing the global orderId
var lastOrderKey = []byte("lastOrderKey")

// GetOrderCount - get the last count
func (k Keeper) GetOrderCount(ctx sdk.Context, orderInitiator sdk.AccAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(orderInitiator)

	if bz == nil {
		return 0
	}
	count, _ := strconv.ParseUint(string(bz), 0, 64)
	return count
}

// SetOrderCount set the last count
func (k Keeper) SetOrderCount(ctx sdk.Context, orderInitiator sdk.AccAddress, count uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(orderInitiator, []byte(strconv.FormatUint(count, 16)))
}

// Implements sdk.AccountMapper.
func (k Keeper) SetInitOrder(ctx sdk.Context, orderInitiator sdk.AccAddress, initOrder MsgInitOrder) {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.MarshalBinaryLengthPrefixed(initOrder)
	if err != nil {
		panic(err)
	}
	store.Set(KeyInitOrder(orderInitiator), bz)
}

// Implements sdk.AccountMapper.
func (k Keeper) GetInitOrder(ctx sdk.Context, orderInitiator sdk.AccAddress) MsgInitOrder {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(KeyInitOrder(orderInitiator))
	if bz != nil {
		return MsgInitOrder{}
	}

	var initOrder MsgInitOrder
	err := k.cdc.UnmarshalBinaryLengthPrefixed(bz, &initOrder)

	if err != nil {
		panic(err)
	}

	return initOrder
}

func (k Keeper) SetFinalizeOrder(ctx sdk.Context, orderInitiator sdk.AccAddress, finalizeOrder MsgFinalizeOrder) {
	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.MarshalBinaryLengthPrefixed(finalizeOrder)
	if err != nil {
		panic(err)
	}
	store.Set(orderInitiator, bz)
}

func (k Keeper) PayForCharging(ctx sdk.Context, buyer sdk.AccAddress, seller sdk.AccAddress, charge int64) (sdk.Tags, sdk.Tags) {

	_, tagsBuyer, _ := k.ck.SubtractCoins(ctx, buyer, sdk.Coins{sdk.NewInt64Coin("byndcoin", charge)})
	_, tagsSeller, _ := k.ck.AddCoins(ctx, seller, sdk.Coins{sdk.NewInt64Coin("byndcoin", charge)})

	return tagsBuyer, tagsSeller
}

// Keeper keys
// Key for getting a specific initOrder from the store

var (
	ByteKeyInitOrder     = []byte("initOrder")
	ByteKeyFinalizeOrder = []byte("finalizeOrder")
)

func KeyInitOrder(orderInitiator sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("%d:%d", ByteKeyInitOrder, orderInitiator))
}

func KeyFinalizeOrder(orderInitiator sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("%d:%d", ByteKeyFinalizeOrder, orderInitiator))
}
