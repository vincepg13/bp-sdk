package mobility

import "github.com/cosmos/cosmos-sdk/codec"

// Register concrete types on wire codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgInitOrder{}, "mobility/InitOrder", nil)
	cdc.RegisterConcrete(MsgFinalizeOrder{}, "mobility/FinalizeOrder", nil)
}
