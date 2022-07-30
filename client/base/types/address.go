package types

import "strings"

const (
	Bech32MainPrefix = "nova"

	Purpose            = 44
	CoinType           = 118
	FullFundraiserPath = "m/44'/118'/0'/0/0"
	PrefixAccount      = "acc"

	PrefixValidator = "val"
	PrefixConsensus = "cons"
	PrefixPublic    = "pub"
	PrefixOperator  = "oper"
	PrefixAddress   = "addr"

	Bech32PrefixAccAddr  = Bech32MainPrefix
	Bech32PrefixAccPub   = Bech32MainPrefix + PrefixPublic
	Bech32PrefixValAddr  = Bech32MainPrefix + PrefixValidator + PrefixOperator
	Bech32PrefixValPub   = Bech32MainPrefix + PrefixValidator + PrefixOperator + PrefixPublic
	Bech32PrefixConsAddr = Bech32MainPrefix + PrefixValidator + PrefixConsensus
	Bech32PrefixConsPub  = Bech32MainPrefix + PrefixValidator + PrefixConsensus + PrefixPublic
)

func ConvertMainPrefix(addr string, old string, new string) string {
	return strings.Replace(addr, old, new, 1)
}
