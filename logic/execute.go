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
