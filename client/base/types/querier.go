package types

import (
	tendermintv1beta1 "github.com/Carina-hackatom/nova/api/cosmos/base/tendermint/v1beta1"
	stakingv1beta1 "github.com/Carina-hackatom/nova/api/cosmos/staking/v1beta1"
	distv1beta1 "github.com/cosmos/cosmos-sdk/x/distribution/types"

	txv1beta1 "github.com/cosmos/cosmos-sdk/types/tx"
)

type BaseQuerier interface {
	txQuerier
	tmQuerier
	stakeQuerier
	distQuerier
}

type txQuerier interface {
	GetTx(hash string) *txv1beta1.GetTxResponse
}

type tmQuerier interface {
	GetNodeRes() *tendermintv1beta1.GetNodeInfoResponse
	GetBlockByHeight(height int64) *tendermintv1beta1.GetBlockByHeightResponse
	GetLatestBlock() *tendermintv1beta1.GetLatestBlockResponse
}

type stakeQuerier interface {
	GetValInfo(valAddr string) *stakingv1beta1.QueryValidatorResponse
	GetHistoricalInfo(height int64) *stakingv1beta1.QueryHistoricalInfoResponse
}

type distQuerier interface {
	GetRewards(delAddr string, valAddr string) *distv1beta1.QueryDelegationRewardsResponse
}
