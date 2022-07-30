package logic

import (
	"github.com/Carina-hackatom/coordinator/client/base"
	"github.com/Carina-hackatom/coordinator/client/base/query"
	novaTx "github.com/Carina-hackatom/coordinator/client/nova/msgs"
	"github.com/Carina-hackatom/coordinator/config"
	"github.com/Carina-hackatom/coordinator/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc"
	"log"
	"os"
	"reflect"
	"time"
)

var (
	Host = &config.HostChainInfo{}
)

func UpdateChainState(host string, ctx client.Context, txf tx.Factory, botInfo keyring.Info, interval int, errLogger *os.File) {

	Host.Set(host)

	conn, err := grpc.Dial(
		Host.GrpcAddr,
		grpc.WithInsecure(),
	)
	utils.CheckErr(err, "cannot create gRPC connection", 0)
	defer func(c *grpc.ClientConn) {
		err = c.Close()
		utils.CheckErr(err, "", 1)
	}(conn)
	cq := &query.CosmosQueryClient{ClientConn: conn}

	i := 0
	intv := time.Duration(interval)
	for {
		botTickLog("Oracle", int(intv)*i)

		delegatedToken, height, apphash := OracleInfo(cq, Host.Validator)
		msg1 := novaTx.MakeMsgUpdateChainState(botInfo.GetAddress(), host, Host.Denom, Host.Decimal, delegatedToken, height, apphash)
		msgs := []sdktypes.Msg{msg1}
		log.Println("----> MsgUpdateChainState was sent")
		base.GenTxWithFactory(errLogger, ctx, txf, false, msgs...)
		time.Sleep(intv * time.Second)
		i++
	}
}

func IcaAutoStake(host string, ctx client.Context, txf tx.Factory, botInfo keyring.Info, interval int, errLogger *os.File) {

	Host.Set(host)

	conn, err := grpc.Dial(
		Host.GrpcAddr,
		grpc.WithInsecure(),
	)
	utils.CheckErr(err, "cannot create gRPC connection", 0)
	defer func(c *grpc.ClientConn) {
		err = c.Close()
		utils.CheckErr(err, "", 1)
	}(conn)
	cq := &query.CosmosQueryClient{ClientConn: conn}

	i := 0
	intv := time.Duration(interval)
	for {
		botTickLog("Re-Staking", int(intv)*i)

		r := RewardsWithAddr(cq, Host.HostAccount, Host.Validator)
		if reflect.DeepEqual(r, sdktypes.DecCoin{}) {
			time.Sleep(intv * time.Second)
			i++
			continue
		}

		msg1 := novaTx.MakeMsgIcaAutoStaking(host, Host.HostAccount, botInfo.GetAddress(), r)
		msgs := []sdktypes.Msg{msg1}
		log.Println("----> MsgIcaAutoStaking was sent")
		base.GenTxWithFactory(errLogger, ctx, txf, false, msgs...)
		time.Sleep(intv * time.Second)
		i++
	}
}

func IcaStake(host string, ctx client.Context, txf tx.Factory, botInfo keyring.Info, chanID string, interval int, errLogger *os.File) {

	i := 0
	intv := time.Duration(interval)
	for {
		botTickLog("ICA-Staking", int(intv)*i)

		msg1 := novaTx.MakeMsgDelegate(host, botInfo.GetAddress(), "transfer", chanID)
		msgs := []sdktypes.Msg{msg1}
		log.Println("----> MsgDelegate was sent")
		base.GenTxWithFactory(errLogger, ctx, txf, false, msgs...)
		time.Sleep(intv * time.Second)
		i++
	}
}

func UndelegateAndWithdraw(host string, ctx client.Context, txf tx.Factory, botInfo keyring.Info, chanID string, interval int, errLogger *os.File) {

	Host.Set(host)

	conn, err := grpc.Dial(
		Host.GrpcAddr,
		grpc.WithInsecure(),
	)
	utils.CheckErr(err, "cannot create gRPC connection", 0)
	defer func(c *grpc.ClientConn) {
		err = c.Close()
		utils.CheckErr(err, "", 1)
	}(conn)
	cq := &query.CosmosQueryClient{ClientConn: conn}

	isStart := true
	i := 0
	intv := time.Duration(interval)
	for {
		botTickLog("Undelegate & Withdraw", int(intv)*i)

		blkTS := LatestBlockTS(cq)
		delegatedToken, height, apphash := OracleInfo(cq, Host.Validator)
		msg1 := novaTx.MakeMsgUpdateChainState(botInfo.GetAddress(), host, Host.Denom, Host.Decimal, delegatedToken, height, apphash)

		if isStart {
			msg2 := novaTx.MakeMsgUndelegate(host, botInfo.GetAddress())
			msgs := []sdktypes.Msg{msg1, msg2}
			log.Println("----> MsgUndelegate was sent")
			base.GenTxWithFactory(errLogger, ctx, txf, false, msgs...)
			isStart = false
		} else {
			msg2 := novaTx.MakeMsgUndelegate(host, botInfo.GetAddress())
			msgs := []sdktypes.Msg{msg1, msg2}
			log.Println("----> MsgUndelegate was sent")
			base.GenTxWithFactory(errLogger, ctx, txf, false, msgs...)

			time.Sleep(60 * time.Second)

			msg3 := novaTx.MakeMsgPendingWithdraw(host, botInfo.GetAddress(), "transfer", chanID, blkTS)
			msgs = []sdktypes.Msg{msg3}
			log.Println("----> MsgPendingWithdraw was sent")
			base.GenTxWithFactory(errLogger, ctx, txf, false, msgs...)
		}

		time.Sleep(intv * time.Second)
		i++
	}
}
