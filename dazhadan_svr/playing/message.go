package playing

import (
	"fmt"
	"dazhadan/dazhadan_svr/card"
)

type MsgType	int
const  (
	MsgGetInitCards	MsgType = iota + 1
	MsgWaitDadu
	MsgConfirmDadu
	MsgSwitchPosition
	MsgStartPlay

	MsgSwitchOperator
	MsgDrop
	MsgGuo
	MsgJiesuan

	MsgEnterRoom
	MsgReadyRoom
	MsgLeaveRoom
	MsgGameEnd
	MsgRoomClosed
)

func (msgType MsgType) String() string {
	switch msgType {
	case MsgGetInitCards:
		return "MsgGetInitCards"
	case MsgWaitDadu:
		return "MsgWaitDadu"
	case MsgConfirmDadu:
		return "MsgConfirmDadu"
	case MsgSwitchPosition:
		return "MsgSwitchPosition"
	case MsgStartPlay:
		return "MsgStartPlay"
	case MsgSwitchOperator:
		return "MsgSwitchOperator"
	case MsgDrop:
		return "MsgDrop"
	case MsgGuo:
		return "MsgGuo"
	case MsgJiesuan:
		return "MsgJiesuan"
	case MsgEnterRoom:
		return "MsgEnterRoom"
	case MsgReadyRoom:
		return "MsgReadyRoom"
	case MsgLeaveRoom:
		return "MsgEnterRoom"
	case MsgGameEnd:
		return "MsgGameEnd"
	case MsgRoomClosed:
		return "MsgRoomClosed"
	}
	return "unknow MsgType"
}

type Message struct {
	Type		MsgType
	Owner 	*Player
	Data 	interface{}
}

func (data *Message) String() string {
	if data == nil {
		return "{nil Message}"
	}
	return fmt.Sprintf("{type=%v, Owner=%v}", data.Type, data.Owner)
}

func newMsg(t MsgType, owner *Player, data interface{}) *Message {
	return &Message{
		Owner:	owner,
		Type: t,
		Data: data,
	}
}

//玩家获得初始牌的消息
type GetInitCardsMsgData struct {
	PlayingCards	*card.PlayingCards
}
func NewGetInitCardsMsg(owner *Player, data *GetInitCardsMsgData) *Message {
	return newMsg(MsgGetInitCards, owner, data)
}

//玩家等待打独的消息
type WaitDaduMsgData struct {
	WaitDaduPlayer *Player
	LeftSec int32
}
func NewWaitDaduMsg(owner *Player, data *WaitDaduMsgData) *Message {
	return newMsg(MsgWaitDadu, owner, data)
}

//玩家确认打独的消息
type ConfirmDaduMsgData struct {
	IsDadu bool
	DaduPlayer *Player
}
func NewConfirmDaduMsg(owner *Player, data *ConfirmDaduMsgData) *Message {
	return newMsg(MsgConfirmDadu, owner, data)
}

//玩家等待打独的消息
type SwitchPositionMsgData struct {
	OppUid uint64
	OppPos int32
	AssistUid uint64
	AssistPos int32
}
func NewSwitchPositionMsg(owner *Player, data *SwitchPositionMsgData) *Message {
	return newMsg(MsgSwitchPosition, owner, data)
}

//开始打牌的消息
type StartPlayMsgData struct {
	IsDadu bool
	Master *Player
	Assist *Player
}
func NewStartPlayMsg(owner *Player, data *StartPlayMsgData) *Message {
	return newMsg(MsgStartPlay, owner, data)
}

//切换玩家消息
type SwitchOperatorMsgData struct {
}
func NewSwitchOperatorMsg(owner *Player, data *SwitchOperatorMsgData) *Message {
	return newMsg(MsgSwitchOperator, owner, data)
}

type PlayerJiesuanData struct {
	P *Player
	Score int32
	Paixing int
}

//结算消息
type JiesuanMsgData struct {
	Scores []*PlayerJiesuanData
}
func NewJiesuanMsg(owner *Player, data *JiesuanMsgData) *Message {
	return newMsg(MsgJiesuan, owner, data)
}

//玩家进入房间的消息
type EnterRoomMsgData struct {
	EnterPlayer *Player
	AllPlayer 	[]*Player
}
func NewEnterRoomMsg(owner *Player, data *EnterRoomMsgData) *Message {
	return newMsg(MsgEnterRoom, owner, data)
}

//玩家进入房间的消息
type ReadyRoomMsgData struct {
	ReadyPlayer *Player
}
func NewReadyRoomMsg(owner *Player, data *ReadyRoomMsgData) *Message {
	return newMsg(MsgReadyRoom, owner, data)
}

//玩家离开房间的消息
type LeaveRoomMsgData struct {
	LeavePlayer *Player
	AllPlayer 	[]*Player
}
func NewLeaveRoomMsg(owner *Player, data *LeaveRoomMsgData) *Message {
	return newMsg(MsgLeaveRoom, owner, data)
}

//一盘游戏结束的消息
type GameEndMsgData struct {}
func NewGameEndMsg(owner *Player, data *GameEndMsgData) *Message{
	return newMsg(MsgGameEnd, owner, data)
}

//房间结束的消息
type RoomClosedMsgData struct {}
func NewRoomClosedMsg(owner *Player, data *RoomClosedMsgData) *Message{
	return newMsg(MsgRoomClosed, owner, data)
}

//出牌的消息
type DropMsgData struct {}
func NewDropMsg(owner *Player, data *DropMsgData) *Message{
	return newMsg(MsgDrop, owner, data)
}

//过牌的消息
type GuoMsgData struct {}
func NewGuoMsg(owner *Player, data *GuoMsgData) *Message{
	return newMsg(MsgGuo, owner, data)
}
