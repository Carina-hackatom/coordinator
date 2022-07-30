package msgs

import (
	"github.com/Carina-hackatom/nova/x/ibcstaking/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
)

func MakeMsgIcaAutoStaking(chainID string, hostAddr string, operator sdktypes.AccAddress, decCoin sdktypes.DecCoin) *types.MsgIcaAutoStaking {
	coin, _ := sdktypes.NormalizeDecCoin(decCoin).TruncateDecimal()
	return types.NewMsgIcaAutoStaking(chainID, hostAddr, operator, coin)
}
