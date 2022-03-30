package cli

import (
	mob "github.com/vincepg13/bp-sdk/beyond/x/mobility"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/utils"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	bank "github.com/cosmos/cosmos-sdk/x/bank/client"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagTo              = "to"
	flagEstimatedAmount = "amount"
	flagAgreedPrice     = "price"
	flagChargeAmount    = "charge"
	coinDenom           = "byndcoin"
)

// SendInitOrderTxCmd will create a send tx and sign it with the given key.
func SendInitOrderTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "initOrder",
		Short: "Create and sign a initOrder tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			toStr := viper.GetString(flagTo)

			to, err := sdk.AccAddressFromBech32(toStr)
			if err != nil {
				return err
			}

			// get estimated energy amount
			amount := viper.GetInt64(flagEstimatedAmount)

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			account, err := cliCtx.GetAccount(from)
			if err != nil {
				return err
			}

			// ensure account has some coins
			// TODO:  calculate and verify if account has enough coins to pay = (EstimatedEnergyAmout * ElectricityPrice)
			// for now only positive balance is verified
			if !account.GetCoins().IsPositive() {
				return errors.Errorf("Address %s doesn't have enough coins to pay for this transaction.", from)
			}

			// build and sign the transaction, then broadcast to Tendermint
			// TODO: price is fixed for now (demo). 1KwH costs 2 bynd coins.
			msg := mob.NewMsgInitOrder(from, to, uint64(amount)*2, uint64(amount))

			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(flagTo, "", "Address of charging station or car (electricity source)")
	cmd.Flags().String(flagEstimatedAmount, "", "Estimated amount of energy to be used in (kWh) ")
	cmd.MarkFlagRequired(flagTo)
	cmd.MarkFlagRequired(flagEstimatedAmount)

	return cmd
}

/* -------------------------------------------------------------------------*/

// SendFinalizeOrderTxCmd will create a send tx and sign it with the given key.
func SendFinalizeOrderTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "finalizeOrder",
		Short: "Create and sign a finalizeOrder tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			toStr := viper.GetString(flagTo)

			to, err := sdk.AccAddressFromBech32(toStr)
			if err != nil {
				return err
			}

			// get charge amount from CLI
			charge := viper.GetInt64(flagChargeAmount)

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			account, err := cliCtx.GetAccount(from)
			if err != nil {
				return err
			}

			// ensure account has some coins
			// TODO:  calculate and verify if account has enough coins to pay = (EstimatedEnergyAmout * ElectricityPrice)
			// for now only positive balance is verified
			if !account.GetCoins().IsPositive() {
				return errors.Errorf("Address %s doesn't have enough coins to pay for this transaction.", from)
			}

			coins, err := sdk.ParseCoins((strconv.Itoa(int(charge * 2))) + coinDenom)
			if err != nil {
				return err
			}
			// build and sign the transaction, then broadcast to Tendermint
			// TODO: price is fixed for now (demo). 1KwH costs 2 bynd coins.
			msgFinalize := mob.NewMsgFinalizeOrder(from, to, uint64(charge)*2, uint64(charge))

			msgSend := bank.CreateMsg(from, to, coins)

			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msgFinalize, msgSend})
		},
	}
	cmd.Flags().String(flagTo, "", "Address of charging station or car (electricity source)")
	cmd.Flags().String(flagEstimatedAmount, "", "Estimated amount of energy to be used in (kWh) ")
	cmd.Flags().String(flagChargeAmount, "", "Actual amount of energy charged (kWh)")

	return cmd
}
