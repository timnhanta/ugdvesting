package types

import (
	"encoding/json"
	fmt "fmt"
	io "io"
	"os"

	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/go-bip39"
)

const (
	mnemonicEntropySize = 256
)

type Vesting struct {
	Amount   string `json:"amount"`
	Start    string `json:"start"`
	Duration string `json:"duration"`
	Parts    int64  `json:"parts"`
}
type PeriodicVestingAccount struct {
	FromAddress sdk.AccAddress
	ToAddress   sdk.AccAddress
	Start       int64
	Periods     []vestingtypes.Period
}

func UgdVesting() map[string]Vesting {
	//fmt.Println("\n-- UgdVesting --")
	jsonData, err := os.Open("/home/team9413/Projects/ccosmos/ugdvesting/vesting.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer jsonData.Close()

	// Read the file contents into a byte slice
	fileContents, err := io.ReadAll(jsonData)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	var data map[string]Vesting
	err = json.Unmarshal(fileContents, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil
	}

	// Print the values in the map
	/*for address, vesting := range data {
		fmt.Println("Address: ", address)
		fmt.Println("Start: ", vesting.Start)
		fmt.Println("Duration: ", vesting.Duration)
		fmt.Println("Parts: ", vesting.Parts)
	}*/

	return data
}

func CosmosVesting(data map[string]Vesting, fromAddress sdk.AccAddress) PeriodicVestingAccount {
	fmt.Println("\n-- CosmosVesting --")
	var vestingAccount PeriodicVestingAccount
	for address, vesting := range data {
		fmt.Println("\nAddress: ", address)
		fmt.Println("Amount: ", vesting.Amount)
		fmt.Println("Start: ", vesting.Start)
		fmt.Println("Duration: ", vesting.Duration)
		fmt.Println("Parts: ", vesting.Parts)
		//toAddr := CreateCosmosKey(address)
		toAddr, err := sdk.AccAddressFromBech32("cosmos1dcm49f8ypyfdj67k0uqecaxu9cye9eyg05etu8")
		if err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Printf("cosmos adr: %s\n", toAddr)
		/*t, _ := time.Parse(time.RFC3339, vesting.Start)
		startTime := t.Unix()*/
		startTime := time.Now().Unix() + 10
		periods := getPeriod(vesting.Amount, vesting.Duration, vesting.Parts)
		fmt.Println("\nunix time:")
		fmt.Println(startTime)
		fmt.Println("\nPeriod:")
		fmt.Println(periods)
		vestingAccount := PeriodicVestingAccount{FromAddress: fromAddress, ToAddress: toAddr, Start: startTime, Periods: periods}
		return vestingAccount
		//msg := vestingtypes.NewMsgCreatePeriodicVestingAccount(fromAddress, toAddr /*to newAddress*/, startTime, periods)
	}

	return vestingAccount
}

func divideRoundUp(dividend, divisor int64) int64 {
	quotient := dividend / divisor
	if dividend%divisor != 0 {
		quotient++
	}
	return quotient
}

func getPeriod(address string, duration string, parts int64) []vestingtypes.Period {
	fmt.Println("\n-- getPeriod --")
	parse, _ := time.ParseDuration(strings.ReplaceAll(strings.ToLower(duration), "p", ""))
	seconds := int64(parse.Seconds())

	length := divideRoundUp(seconds, parts)
	coin, err := sdk.ParseCoinNormalized(address)
	newAmount := divideRoundUp(coin.Amount.Int64(), parts) /* TODO: improve the diveded amount*/
	amount := sdk.Coins{sdk.NewInt64Coin(coin.Denom, newAmount)}
	if err != nil {
		fmt.Println("Error parsing coins:", err)
	}
	fmt.Println("Denom: ", coin.Denom)
	fmt.Println("newAmount: ", newAmount)

	var periods []vestingtypes.Period

	for i := 0; i < int(parts); i++ {
		period := vestingtypes.Period{Length: int64(length), Amount: amount}
		periods = append(periods, period)
	}

	return periods
}

func getCodec() codec.Codec {
	registry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(registry)
	return codec.NewProtoCodec(registry)
}

func CreateCosmosKey(ugdAddress string) sdk.AccAddress {
	cdc := getCodec()
	kb := keyring.NewInMemory(cdc)
	entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
	if err != nil {
		return nil
	}

	mnemonic, err := bip39.NewMnemonic(entropySeed)
	if err != nil {
		return nil
	}
	//mnemonic := "annual snake shine weather jeans rain bless keen uncover prize salute thunder car want speak abandon either sea orchard dice solid bitter satisfy jar"

	k, err := kb.NewAccount(
		ugdAddress,
		mnemonic,
		"", hd.CreateHDPath(118, 0, 0).String(),
		hd.Secp256k1,
	)
	if err != nil {
		panic(err)
	}

	out, err := keyring.MkAccKeyOutput(k)
	if err != nil {
		panic(err)
	}

	addr, err := sdk.AccAddressFromBech32(out.Address)
	if err != nil {
		panic(err)
	}

	return addr
}
