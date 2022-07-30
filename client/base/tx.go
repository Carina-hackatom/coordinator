package base

import (
	"errors"
	"fmt"
	"github.com/Carina-hackatom/coordinator/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type (
	AccAddr = sdktypes.AccAddress
)

var (
	wm sync.RWMutex
)

func CheckAccAddr(target any) (AccAddr, error) {
	switch target := target.(type) {
	case AccAddr:
		return target, nil
	case string:
		addr, err := sdktypes.AccAddressFromBech32(target)
		if err != nil {
			return nil, err
		}
		return addr, nil
	case []byte:
		return target, nil
	default:
		return nil, errors.New("cannot covert target to AccAddress")
	}
}

func GenTxWithFactory(errFd *os.File, ctx client.Context, txf tx.Factory, onlyGen bool, msgs ...sdktypes.Msg) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Get panic while generating tx")
		}
	}()
	wm.Lock()
	defer wm.Unlock()
	if onlyGen {
		ctx = ctx.WithGenerateOnly(true)
	}

	err := tx.GenerateOrBroadcastTxWithFactory(ctx, txf, msgs...)
	if err != nil {
		if strings.Contains(err.Error(), "account sequence mismatch") {
			utils.LogErrWithFd(errFd, err, "", 1)
			for {
				txseq := txf.Sequence()
				txf.WithSequence(txseq)
				err = tx.GenerateOrBroadcastTxWithFactory(ctx, txf, msgs...)
				if !strings.Contains(err.Error(), "account sequence mismatch") {
					break
				}
				utils.LogErrWithFd(errFd, err, "", 1)
				time.Sleep(6 * time.Second)
			}

		} else {
			utils.LogErrWithFd(errFd, err, "something went wrong while make tx", 1)
		}
	} else {
		_, err = ctx.Output.Write([]byte(fmt.Sprintf("%v: Tx was generated\n\n", time.Now())))
		utils.CheckErr(err, "cannot write log on output", 1)
	}
}

func MakeTxFactory(ctx client.Context, gas string, gasPrice string, memo string, gasWeight float64) tx.Factory {
	gasSetting, _ := flags.ParseGasSetting(gas)

	initFac := tx.Factory{}.
		WithAccountNumber(0).
		WithSequence(0).
		WithTimeoutHeight(0).
		WithTxConfig(ctx.TxConfig).
		WithChainID(ctx.ChainID).
		WithKeybase(ctx.Keyring).
		WithAccountRetriever(ctx.AccountRetriever).
		WithGas(gasSetting.Gas).
		WithSimulateAndExecute(gasSetting.Simulate).
		WithGasAdjustment(flags.DefaultGasAdjustment * gasWeight).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)

	return initFac.
		WithGasPrices(gasPrice).
		WithMemo(memo)

}
