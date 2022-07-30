package msgs

import (
	"github.com/Carina-hackatom/nova/x/gal/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"time"
)

func MakeMsgDelegate(chainID string, operator sdktypes.AccAddress, portID string, chanID string) *types.MsgDelegate {
	return types.NewMsgDelegate(chainID, operator, portID, chanID)
}

func MakeMsgPendingWithdraw(chainID string, operator sdktypes.AccAddress, portID string, chanID string, blockTs time.Time) *types.MsgPendingWithdraw {
	return types.NewMsgPendingWithdraw(chainID, operator, portID, chanID, blockTs)
}

func MakeMsgUndelegate(chainID string, operator sdktypes.AccAddress) *types.MsgUndelegate {
	return types.NewMsgUndelegate(chainID, operator)
}
