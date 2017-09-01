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
	isScramble		bool
	Paixing		        int                     //牌型
	maxid		        int                     //手牌最大的id
	roundScore              int32                   //本轮得分
	totalCoin	        int32                   //总金币

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

		isScramble:     false,
		maxid:   1,
		roundScore:     0,
		totalCoin:     0,
		Paixing:   	card.DouniuType_Meiniu,
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

func (player *Player) GetRoundScore() int32 {
	return player.roundScore
}

func (player *Player) SetRoundScore(round_score int32) {
	player.roundScore = round_score
}

func (player *Player) GetIsDadu() bool {
	return player.isDadu
}

func (player *Player) SetIsDadu(is_dadu bool) {
	player.isDadu = is_dadu
}

func (player *Player) GetIsScramble() bool {
	return player.isScramble
}

func (player *Player) SetIsScramble(is_scramble bool) {
	player.isScramble = is_scramble
}

func (player *Player) GetPaixing() int {
	return player.Paixing
}

func (player *Player) SetPaixing(paixing int) {
	player.Paixing = paixing
}

func (player *Player) GetMaxid() int {
	return player.maxid
}

func (player *Player) SetMaxid(maxid int) {
	player.maxid = maxid
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
	player.SetIsScramble(false)
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

func (player *Player) OperateScramble(scramble_multiple int32) bool{
	log.Debug(time.Now().Unix(), player, "OperateScramble", player.room)
	if player.room == nil || player.GetIsScramble(){
		return false
	}

	if scramble_multiple < 0 || scramble_multiple > 4 {
		log.Error("Player is not playing", player)
		return false
	}

	data := &OperateScrambleData{ScrambleMultiple:scramble_multiple}
	op := NewOperateScramble(player, data)
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

func (player *Player) Scramble(multiple int32) {
	log.Debug(time.Now().Unix(), player, "Scramble", player.room)
	player.SetIsScramble(true)
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
	if _, ok := op.Data.(*OperateDropData); ok {
		/*if op.Operator == player {
			return
		}*/
		msgData := &DropMsgData{}
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
