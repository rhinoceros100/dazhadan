package playing

import (
	"fmt"
	"dazhadan/dazhadan_svr/card"
)

type MsgType	int
const  (
	MsgGetInitCards	MsgType = iota + 1
	MsgWaitPlayAlone
	MsgConfirmPlayAlone
	MsgSwitchPosition
	MsgStartPlay

	MsgSwitchOperator
	MsgDrop
	MsgPass
	MsgSummary

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
	case MsgWaitPlayAlone:
		return "MsgWaitPlayAlone"
	case MsgConfirmPlayAlone:
		return "MsgConfirmPlayAlone"
	case MsgSwitchPosition:
		return "MsgSwitchPosition"
	case MsgStartPlay:
		return "MsgStartPlay"
	case MsgSwitchOperator:
		return "MsgSwitchOperator"
	case MsgDrop:
		return "MsgDrop"
	case MsgPass:
		return "MsgPass"
	case MsgSummary:
		return "MsgSummary"
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
type WaitPlayAloneMsgData struct {
	WaitPlayAlonePlayer *Player
	LeftSec int32
}
func NewWaitPlayAloneMsg(owner *Player, data *WaitPlayAloneMsgData) *Message {
	return newMsg(MsgWaitPlayAlone, owner, data)
}

//玩家确认打独的消息
type ConfirmPlayAloneMsgData struct {
	IsPlayAlone bool
	PlayAlonePlayer *Player
}
func NewConfirmPlayAloneMsg(owner *Player, data *ConfirmPlayAloneMsgData) *Message {
	return newMsg(MsgConfirmPlayAlone, owner, data)
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
	IsPlayAlone bool
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

type PlayerSummaryData struct {
	P *Player
	Rank int32
	Coin int32
	Score int32
	Prize int32
	TotalCoin int32
	PrizeCoin int32
	IsWin bool
}

//结算消息
type SummaryMsgData struct {
	Scores []*PlayerSummaryData
	InfoType int32
}
func NewSummaryMsg(owner *Player, data *SummaryMsgData) *Message {
	return newMsg(MsgSummary, owner, data)
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
type TotalSummaryData struct {
	P *Player
	WinNum int32
	ShuangjiNum int32
	PaSuccNum int32
	TotalPrize int32
	TotalCoin int32

	IsWinner bool
	IsMostWinner bool
	IsMostLoser bool
	IsCreator bool
}
type RoomClosedMsgData struct {
	Summaries []*TotalSummaryData
}
func NewRoomClosedMsg(owner *Player, data *RoomClosedMsgData) *Message{
	return newMsg(MsgRoomClosed, owner, data)
}

//出牌的消息
type DropMsgData struct {
	WhatGroup []*card.Card
	TableScore int32
	CardsType int
	PlaneNum int
	Weight int
}
func NewDropMsg(owner *Player, data *DropMsgData) *Message{
	return newMsg(MsgDrop, owner, data)
}

//过牌的消息
type PassMsgData struct {}
func NewPassMsg(owner *Player, data *PassMsgData) *Message{
	return newMsg(MsgPass, owner, data)
}
