package logic

import (
	"github.com/Carina-hackatom/coordinator/client/base/query"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"log"
	"time"
)

func OracleInfo(cq *query.CosmosQueryClient, validatorAddr string) (string, int64, []byte) {
	h := cq.GetLatestBlock().GetBlock().GetHeader().GetHeight()
	hisInfo := cq.GetHistoricalInfo(h)
	apphash := hisInfo.GetHist().GetHeader().GetAppHash()

	var delegatedToken string
	for _, val := range hisInfo.GetHist().GetValset() {
		if val.GetOperatorAddress() != validatorAddr {
			continue
		} else {
			delegatedToken = val.GetTokens()
			break
		}
	}
	return delegatedToken, h, apphash
}

func RewardsWithAddr(cq *query.CosmosQueryClient, delegator string, validator string) sdktypes.DecCoin {
	var reward sdktypes.DecCoin
	defer func() {
		if err := recover(); err != nil {
			log.Println("There is no reward to handle")
			reward = sdktypes.DecCoin{}
		}
	}()
	reward = cq.GetRewards(delegator, validator).GetRewards()[0]
	return reward
}

func LatestBlockTS(cq *query.CosmosQueryClient) time.Time {
	secs := cq.GetLatestBlock().GetBlock().GetHeader().GetTime().GetSeconds()
	nanos := cq.GetLatestBlock().GetBlock().GetHeader().GetTime().GetNanos()
	currentTs := time.Unix(secs, int64(nanos)).UTC()
	return currentTs
}
