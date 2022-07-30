package logic

import (
	"github.com/Carina-hackatom/coordinator/client/base/query"
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
