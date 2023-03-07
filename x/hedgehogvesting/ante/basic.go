package app

import (
	"log"
	"os"
	"strings"
	"time"

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

type MyError struct{}

const (
	denom = "ugd"
)

func NewValidateBasicDecorator(bk bankkeeper.Keeper) ValidateBasicDecorator {
	return ValidateBasicDecorator{
		bankKeeper: bk,
	}
}

func (m *MyError) Error() string {
	return "Custom error for testing"
}

func (vbd ValidateBasicDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// no need to validate basic on recheck tx, call next antehandler
	if ctx.IsReCheckTx() {
		return next(ctx, tx, simulate)
	}
	// Test Start -  Create a log file
	file, errr := os.OpenFile("/home/team9413/Projects/ccosmos/ugdvesting/app_ante.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if errr != nil {
		log.Fatal(errr)
	}
	defer file.Close()
	log.SetOutput(file)

	for _, msg := range tx.GetMsgs() {
		if msgBank, ok := msg.(*banktypes.MsgSend); ok {
			log.Println("------------")
			log.Println("  FromAddress: ", msgBank.FromAddress)
			log.Println("  ToAddress: ", msgBank.ToAddress)
			log.Println("  Amount: ", msgBank.Amount)
			log.Println("------------")
			addr, err := sdk.AccAddressFromBech32(msgBank.FromAddress)
			log.Println("addr: ", addr)
			coin := vbd.bankKeeper.GetBalance(ctx, addr, denom)
			log.Println(coin)
			if err != nil {
				log.Fatal(err)
			}

			data := types.UgdVesting()
			log.Println("data: ", data)
			if val, ok := data[msgBank.FromAddress]; ok {
				unvestedAmount := getUnvestedAmount(val)
				log.Println("unvestedAmount: ", unvestedAmount)
			}
		}
	}
	// Test End -  Create a log file

	if err := tx.ValidateBasic(); err != nil {
		err := &MyError{}
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

func getUnvestedAmount(val types.Vesting) float64 {
	log.Println("------------")
	coinVesting, err := sdk.ParseCoinNormalized(val.Amount)
	if err != nil {
		log.Fatal(err)
	}
	dateStart, _ := time.Parse(time.RFC3339, val.Start)
	dateNow := time.Now()
	/*dateStringNow := "2023-03-02T16:40:54.378693Z"
	dateNow, _ := time.Parse(time.RFC3339, dateStringNow)*/
	log.Println("Now: ", dateNow)
	log.Println("Start: ", dateStart)
	timePassed := dateNow.Sub(dateStart)
	duration, _ := time.ParseDuration(strings.ReplaceAll(strings.ToLower(val.Duration), "p", ""))
	if timePassed.Milliseconds() < duration.Milliseconds() {
		//dateEnd := dateStart.Add(duration)
		durationPart := duration.Milliseconds() / 6
		var amountPart float64
		amountPart = float64(coinVesting.Amount.Int64() / 6)
		part := timePassed.Milliseconds() / durationPart
		amountVested := float64(part) * amountPart
		amountUnVested := float64(coinVesting.Amount.Int64()) - amountVested
		log.Println("Amount: ", amountVested)
		log.Println("Duration (hour): ", duration.Hours())
		log.Println("Duration passed (hour): ", timePassed.Hours())
		log.Println("Part: ", part)
		log.Println("amountVested: ", amountVested)
		log.Println("amountUnVested: ", amountUnVested)
		return amountUnVested
	}

	return 0
}
