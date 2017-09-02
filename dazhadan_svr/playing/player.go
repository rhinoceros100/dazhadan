package playing

import (
	"dazhadan/dazhadan_svr/card"
	"dazhadan/dazhadan_svr/log"
	"time"
	"fmt"
)

type PlayerObserver interface {
	OnMsg(player *Player, msg *Message)
}

type Player struct {
	id			uint64			//玩家id
	position		int32			//玩家在房间的位置
	room			*Room			//玩家所在的房间
	isReady			bool

	isDadu			bool
	isEndPlaying		bool
	needDrop		bool
	totalCoin	        int32                   //总金币
	rank		        int32                   //排名
	score		        int32                   //一轮得分
	prize		        int32                   //获得奖励次数

	playingCards 	*card.PlayingCards	//玩家手上的牌
	niuCards         []*card.Card
	observers	 []PlayerObserver
}

func NewPlayer(id uint64) *Player {
	player :=  &Player{
		id:		id,
		position:       10,
		isReady:        false,
		isDadu:     	false,
		isEndPlaying:   false,
		needDrop:     	false,

		rank:     0,
		score:     0,
		prize:     0,
		totalCoin:     0,
		playingCards:	card.NewPlayingCards(),
		observers:	make([]PlayerObserver, 0),
		niuCards:       make([]*card.Card, 0),
	}
	return player
}

func (player *Player) IsMaster() bool {
	return player == player.room.masterPlayer
}

func (player *Player) GetId() uint64 {
	return player.id
}

func (player *Player) GetPosition() int32 {
	return player.position
}

func (player *Player) GetTotalCoin() int32 {
	return player.totalCoin
}

func (player *Player) AddTotalCoin(add int32) int32 {
	player.totalCoin += add
	return player.totalCoin
}

func (player *Player) GetScore() int32 {
	return player.score
}

func (player *Player) AddScore(score int32) {
	player.score += score
}

func (player *Player) ResetScore() {
	player.score = 0
}

func (player *Player) GetRank() int32 {
	return player.rank
}

func (player *Player) SetRank(rank int32) {
	player.rank = rank
}

func (player *Player) GetPrize() int32 {
	return player.score
}

func (player *Player) AddPrize(prize int32) {
	player.prize += prize
}

func (player *Player) ResetPrize() {
	player.prize = 0
}

func (player *Player) GetIsDadu() bool {
	return player.isDadu
}

func (player *Player) SetIsDadu(is_dadu bool) {
	player.isDadu = is_dadu
}

func (player *Player) GetIsEndPlaying() bool {
	return player.isEndPlaying
}

func (player *Player) SetIsIsEndPlaying(is_end_playing bool) {
	player.isEndPlaying = is_end_playing
}

func (player *Player) GetNeedDrop() bool {
	return player.needDrop
}

func (player *Player) SetNeedDrop(need_drop bool) {
	player.needDrop = need_drop
}

func (player *Player) GetNiuCards() []*card.Card {
	return player.niuCards
}

func (player *Player) SetNiuCards(niu_cards []*card.Card) {
	player.niuCards = niu_cards
}

func (player *Player) Reset() {
	//log.Debug(time.Now().Unix(), player,"Player.Reset")
	player.playingCards.Reset()
	player.SetIsReady(false)
	player.SetIsDadu(false)
	player.SetIsIsEndPlaying(false)
	player.SetNeedDrop(false)
	player.SetRank(0)
	player.ResetPrize()
	player.ResetScore()
}

func (player *Player) AddObserver(ob PlayerObserver) {
	player.observers = append(player.observers, ob)
}

func (player *Player) AddCard(card *card.Card) {
	//log.Debug(time.Now().Unix(), player, "Player.AddCard :", card)
	player.playingCards.AddCard(card)
}

func (player *Player) OperateEnterRoom(room *Room) bool{
	//log.Debug(time.Now().Unix(), player, "OperateEnterRoom room :", room)
	for _, room_player := range room.players{
		if room_player == player{
			log.Error("Player already in room:", player)
			return false
		}
	}

	data := &OperateEnterRoomData{}
	op := NewOperateEnterRoom(player, data)
	room.PlayerOperate(op)
	return player.waitResult(op.ResultCh)
}

func (player *Player) OperateLeaveRoom() bool{
	//log.Debug(time.Now().Unix(), player, "OperateLeaveRoom", player.room)
	if player.room == nil {
		return true
	}
	room_status := player.room.roomStatus
	if room_status > RoomStatusWaitAllPlayerEnter {
		log.Error("Wrong room status:", room_status)
		return false
	}

	data := &OperateLeaveRoomData{}
	op := NewOperateLeaveRoom(player, data)
	player.room.PlayerOperate(op)
	return player.waitResult(op.ResultCh)
}

func (player *Player) OperateDoReady() bool{
	//log.Debug(time.Now().Unix(), player, "OperateDoReady", player.room)
	if player.room == nil || player.GetIsReady(){
		return false
	}
	room_status := player.room.roomStatus
	if room_status != RoomStatusWaitAllPlayerEnter && room_status != RoomStatusWaitAllPlayerReady {
		log.Error("Wrong room status:", room_status)
		return false
	}

	data := &OperateReadyRoomData{}
	op := NewOperateReadyRoom(player, data)
	player.room.PlayerOperate(op)
	return player.waitResult(op.ResultCh)
}

func (player *Player) GetIsReady() bool {
	return player.isReady
}

func (player *Player) SetIsReady(is_ready bool) {
	player.isReady = is_ready
}

func (player *Player) GetPlayingCards() *card.PlayingCards {
	return player.playingCards
}

func (player *Player) waitResult(resultCh chan bool) bool{
	log_time := time.Now().Unix()
	select {
	case <- time.After(time.Second * 10):
		close(resultCh)
		log.Debug(time.Now().Unix(), player, "Player.waitResult timeout")
		return false
	case result := <- resultCh:
		log.Debug(log_time, player, "Player.waitResult result :", result)
		return result
	}
	log.Debug(log_time, player, "Player.waitResult fasle")
	return false
}

func (player *Player) EnterRoom(room *Room, position int32) {
	log.Debug(time.Now().Unix(), player, "enter", room)
	player.room = room
	player.position = position
}

func (player *Player) ReadyRoom(room *Room) {
	log.Debug(time.Now().Unix(), player, "ready", room)
}

func (player *Player) LeaveRoom() {
	log.Debug(time.Now().Unix(), player, "leave", player.room)
	player.room = nil
	player.position = -1
}

func (player *Player) Dadu(is_dadu bool) {
	//log.Debug(time.Now().Unix(), player, "Dadu", player.room)
	player.SetIsDadu(is_dadu)
}

func (player *Player) String() string{
	if player == nil {
		return "{player=nil}"
	}
	return fmt.Sprintf("{player=%v, pos=%v}", player.id, player.position)
}

//玩家成功操作的通知
func (player *Player) OnPlayerSuccessOperated(op *Operate) {
	//log.Debug(time.Now().Unix(), player, "OnPlayerSuccessOperated", op)
	switch op.Op {
	case OperateEnterRoom:
		player.onPlayerEnterRoom(op)
	case OperateLeaveRoom:
		player.onPlayerLeaveRoom(op)
	case OperateReadyRoom:
		player.onPlayerReadyRoom(op)
	case OperateConfirmDadu:
		player.OnPlayerDadu(op)
	case OperateSwitchOperator:
		player.onSwithOperator(op)
	case OperateDrop:
		player.OnDrop(op)
	case OperateGuo:
		player.OnGuo(op)
	}
}

func (player *Player) notifyObserver(msg *Message) {
	//log.Debug(time.Now().Unix(), player, "notifyObserverMsg", msg)
	for _, ob := range player.observers {
		ob.OnMsg(player, msg)
	}
}

//begin player operate event handler

func (player *Player) onPlayerEnterRoom(op *Operate) {
	if _, ok := op.Data.(*OperateEnterRoomData); ok {
		if player.room == nil {
			return
		}
		msgData := &EnterRoomMsgData{
			EnterPlayer : op.Operator,
			AllPlayer: player.room.players,
		}
		player.notifyObserver(NewEnterRoomMsg(player, msgData))
	}
}

func (player *Player) onPlayerReadyRoom(op *Operate) {
	//log.Debug(time.Now().Unix(), player, "onPlayerReadyRoom")

	data := &ReadyRoomMsgData{
		ReadyPlayer:op.Operator,
	}
	player.notifyObserver(NewReadyRoomMsg(player, data))
}

func (player *Player) onPlayerLeaveRoom(op *Operate) {
	if _, ok := op.Data.(*OperateLeaveRoomData); ok {
		if op.Operator == player {
			return
		}
		if player.room == nil {
			return
		}
		msgData := &LeaveRoomMsgData{
			LeavePlayer : op.Operator,
			AllPlayer: player.room.players,
		}
		player.notifyObserver(NewLeaveRoomMsg(player, msgData))
	}
}

func (player *Player) OnPlayerDadu(op *Operate) {
	//log.Debug(time.Now().Unix(), player, "OnPlayerDadu")
	if dadu_data, ok := op.Data.(*OperateConfirmDaduData); ok {
		data := &ConfirmDaduMsgData{
			IsDadu:dadu_data.IsDadu,
			DaduPlayer:op.Operator,
		}
		player.notifyObserver(NewConfirmDaduMsg(player, data))
	}
}


func (player *Player) onSwithOperator(op *Operate) {
	if _, ok := op.Data.(*OperateSwitchOperatorData); ok {
		/*if op.Operator == player {
			return
		}*/
		msgData := &SwitchOperatorMsgData{}
		player.notifyObserver(NewSwitchOperatorMsg(op.Operator, msgData))
	}
}

func (player *Player) OnDrop(op *Operate) {
	if drop_data, ok := op.Data.(*OperateDropData); ok {
		/*if op.Operator == player {
			return
		}*/
		msgData := &DropMsgData{
			WhatGroup:drop_data.whatGroup,
		}
		player.notifyObserver(NewDropMsg(op.Operator, msgData))
	}
}

func (player *Player) OnGuo(op *Operate) {
	if _, ok := op.Data.(*OperateGuoData); ok {
		/*if op.Operator == player {
			return
		}*/
		msgData := &GuoMsgData{}
		player.notifyObserver(NewGuoMsg(op.Operator, msgData))
	}
}

func (player *Player) OnWaitWaitDadu(msg *Message) {
	player.notifyObserver(msg)
}

func (player *Player) OnSwitchPosition(msg *Message) {
	player.notifyObserver(msg)
}

func (player *Player) OnStartPlay(msg *Message) {
	player.notifyObserver(msg)
}

func (player *Player) OnJiesuan(msg *Message) {
	//log.Debug(time.Now().Unix(), player, "OnJiesuan")

	player.notifyObserver(msg)
}

func (player *Player) OnGetInitCards() {
	//log.Debug(time.Now().Unix(), player, "OnGetInitCards", player.playingCards)

	data := &GetInitCardsMsgData{
		PlayingCards: player.playingCards,
	}
	player.notifyObserver(NewGetInitCardsMsg(player, data))
}

func (player *Player) OnRoomClosed() {
	//log.Debug(time.Now().Unix(), player, "OnRoomClosed")
	player.room = nil
	//player.Reset()

	data := &RoomClosedMsgData{}
	player.notifyObserver(NewRoomClosedMsg(player, data))
}

func (player *Player) OnEndPlayGame() {
	//log.Debug(time.Now().Unix(), player, "OnPlayingGameEnd")
	player.Reset()
	data := &GameEndMsgData{}
	player.notifyObserver(NewGameEndMsg(player, data))
}

func (player *Player) GetTailCard(num int) []*card.Card {
	//log.Debug(time.Now().Unix(), player, "GetTailCard", num)
	return player.playingCards.Tail(num)
}

func (player *Player) GetLeftCardNum() (int) {
	return player.playingCards.CardsInHand.Len()
}

func (player *Player) Drop(cards []*card.Card) bool {
	log.Debug(time.Now().Unix(), player, "Drop card :", cards)
	return player.playingCards.DropCards(cards)
}
