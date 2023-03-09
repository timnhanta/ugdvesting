package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/timnhanta/ugdvesting/x/hedgehogvesting/types"
)

// ValidateBasicDecorator will call tx.ValidateBasic and return any non-nil error.
// If ValidateBasic passes, decorator calls next AnteHandler in chain. Note,
// ValidateBasicDecorator decorator will not get executed on ReCheckTx since it
// is not dependent on application state.
type ValidateBasicDecorator struct {
	bankKeeper bankkeeper.Keeper
}

type MyError struct {
	message string
}

func NewValidateBasicDecorator(bk bankkeeper.Keeper) ValidateBasicDecorator {
	return ValidateBasicDecorator{
		bankKeeper: bk,
	}
}

func (err MyError) Error() string {
	return err.message
}

func (vbd ValidateBasicDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// no need to validate basic on recheck tx, call next antehandler
	if ctx.IsReCheckTx() {
		return next(ctx, tx, simulate)
	}
	// Test Start -  Create a log file
	allowTransaction := true
	wd, err := os.Getwd()
	if err != nil {
		return ctx, err
	}
	file, err := os.OpenFile(wd+"/_ante.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return ctx, err
	}
	defer file.Close()
	log.SetOutput(file)

	for _, msg := range tx.GetMsgs() {
		if msgBank, ok := msg.(*banktypes.MsgSend); ok {
			addr, err := sdk.AccAddressFromBech32(msgBank.FromAddress)
			account := vbd.bankKeeper.GetBalance(ctx, addr, types.Denom)
			if err != nil {
				return ctx, err
			}
			// Test hedgehog endpoint get vesting
			//types.HegdehogRequestGetVestingByAddr(addr.String())
			//log.Println("\n-		Hedgehog get is ok if we get here!!!")

			// Create a fake JSON response
			response := []byte(`{
				"amount":"7000ugd",
				"start": "2023-03-05T13:27:54.378693Z",
				"duration": "P168H30S",
				"parts": 7
			}`)

			var val types.Vesting

			err = json.Unmarshal(response, &val)
			if err != nil {
				return ctx, err
			}

			log.Println("Working directory:", wd)
			log.Println("Hedgehog response:", response)
			log.Println("------------")
			log.Println("  FromAddress: ", msgBank.FromAddress)
			log.Println("  ToAddress: ", msgBank.ToAddress)
			log.Println("  Amount: ", msgBank.Amount)
			log.Println("------------")
			log.Println("addr: ", addr.String())
			log.Println("Account coin: ", account)
			//log.Println("data: ", data)
			//if val, ok := data[msgBank.FromAddress]; ok {
			unvestedAmount := types.GetUnvestedAmount(val)
			messageAmount := msgBank.Amount.AmountOf(types.Denom).Int64()
			accountAmount := account.Amount.Int64()
			log.Println("val: ", val)
			log.Println("accountAmount: ", accountAmount)
			log.Println("messageAmount: ", messageAmount)
			log.Println("unvestedAmount: ", unvestedAmount)
			log.Println("ugd in addr after transation: ", (float64(accountAmount) - (float64(messageAmount))))

			// check if transaction is allowed based on unvested, transaction and account balance
			allowTransaction = unvestedAmount <= 0 || float64(accountAmount) >= unvestedAmount+float64(messageAmount)
			if !allowTransaction {
				err := MyError{
					fmt.Sprintf(
						"Vesting error: %v with unvested amount of %v, need to have at least %v, but has %v, when doing a transaction of %v.",
						addr.String(),
						unvestedAmount,
						unvestedAmount+float64(messageAmount),
						accountAmount,
						messageAmount,
					),
				}
				return ctx, err
			}
			//}
		}
	}
	// Test End -  Create a log file

	if err := tx.ValidateBasic(); err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}
