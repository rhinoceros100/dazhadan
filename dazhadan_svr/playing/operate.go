package playing

import (
	"fmt"
	"dazhadan/dazhadan_svr/card"
)

type OperateType int

const (
	OperateEnterRoom	OperateType = iota + 1
	OperateReadyRoom
	OperateLeaveRoom

	OperateConfirmPlayAlone
	OperateSwitchOperator
	OperateDrop
	OperatePass
)

func (operateType OperateType) String() string {
	switch operateType {
	case OperateEnterRoom :
		return "OperateEnterRoom"
	case OperateReadyRoom :
		return "OperateReadyRoom"
	case OperateLeaveRoom:
		return "OperateLeaveRoom"
	case OperateConfirmPlayAlone:
		return "OperateConfirmPlayAlone"
	case OperateSwitchOperator:
		return "OperateSwitchOperator"
	case OperateDrop:
		return "OperateDrop"
	case OperatePass:
		return "OperatePass"
	}
	return "unknow OperateType"
}

type Operate struct {//玩家操作
	Op			OperateType
	Operator	*Player				//操作者
	Data		interface{}
	ResultCh		chan bool
}

func (op *Operate) String() string {
	if op == nil {
		return "{operator=nil, op=nil}"
	}
	return fmt.Sprintf("{operator=%v, op=%v}", op.Operator, op.Op)
}

func newOperate(op OperateType, operator *Player, data interface{}) *Operate{
	return &Operate{
		Op:	op,
		Data: data,
		Operator: operator,
		ResultCh: make(chan bool, 1),
	}
}

type OperateEnterRoomData struct {
}
func NewOperateEnterRoom(operator *Player, data *OperateEnterRoomData) *Operate {
	return newOperate(OperateEnterRoom, operator, data)
}

type OperateReadyRoomData struct {
}
func NewOperateReadyRoom(operator *Player, data *OperateReadyRoomData) *Operate {
	return newOperate(OperateReadyRoom, operator, data)
}

type OperateLeaveRoomData struct {
}
func NewOperateLeaveRoom(operator *Player, data *OperateLeaveRoomData) *Operate {
	return newOperate(OperateLeaveRoom, operator, data)
}

type OperateConfirmPlayAloneData struct {
	IsPlayAlone bool
}
func NewOperateConfirmPlayAlone(operator *Player, data *OperateConfirmPlayAloneData) *Operate {
	return newOperate(OperateConfirmPlayAlone, operator, data)
}

type OperateSwitchOperatorData struct {
	MustDrop bool
	CanDrop bool
}
func NewSwitchOperator(operator *Player, data *OperateSwitchOperatorData) *Operate {
	return newOperate(OperateSwitchOperator, operator, data)
}

type OperateDropData struct {
	whatGroup []*card.Card
	cardsType int
	planeNum int
	weight int
}
func NewOperateDrop(operator *Player, data *OperateDropData) *Operate {
	return newOperate(OperateDrop, operator, data)
}

type OperatePassData struct {}
func NewOperatePass(operator *Player, data *OperatePassData) *Operate {
	return newOperate(OperatePass, operator, data)
}
