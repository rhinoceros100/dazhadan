package playing

import (
	"dazhadan/dazhadan_svr/card"
	"dazhadan/dazhadan_svr/log"
	"dazhadan/dazhadan_svr/util"
	"time"
	"fmt"
)

type RoomStatusType int
const (
	RoomStatusWaitAllPlayerEnter	RoomStatusType = iota	// 等待玩家进入房间
	RoomStatusWaitAllPlayerReady				// 等待玩家准备
	RoomStatusGameStart					// 发牌开始打
	RoomStatusPlayGame					// 正在进行游戏，结束后会进入RoomStatusShowCards
	RoomStatusEndPlayGame					// 游戏结束后会回到等待游戏开始状态，或者进入结束房间状态
	RoomStatusRoomEnd					// 房间结束状态
)

func (status RoomStatusType) String() string {
	switch status {
	case RoomStatusWaitAllPlayerEnter :
		return "RoomStatusWaitAllPlayerEnter"
	case RoomStatusWaitAllPlayerReady:
		return "RoomStatusWaitAllPlayerReady"
	case RoomStatusGameStart:
		return "RoomStatusGameStart"
	case RoomStatusPlayGame:
		return "RoomStatusPlayGame"
	case RoomStatusEndPlayGame:
		return "RoomStatusEndPlayGame"
	case RoomStatusRoomEnd:
		return "RoomStatusRoomEnd"
	}
	return "unknow RoomStatus"
}

type RoomObserver interface {
	OnRoomClosed(room *Room)
}

type Room struct {
	id			uint64					//房间id
	config 			*RoomConfig				//房间配置
	players 		[]*Player				//当前房间的玩家列表

	observers		[]RoomObserver				//房间观察者，需要实现OnRoomClose，房间close的时候会通知它
	roomStatus		RoomStatusType				//房间当前的状态
	playedGameCnt		int					//已经玩了的游戏的次数

	//begin playingGameData, reset when start playing game
	cardPool		*card.Pool				//洗牌池
	opMaster		*Player					//当前出过牌的玩家
	waitOperator		*Player					//等待出牌的玩家
	masterPlayer 		*Player					//庄
	assistPlayer 		*Player					//和庄一家
	daduPlayer 		*Player					//打独的玩家
	isDadu			bool					//是否打独
	fanpai			*card.Card				//翻出的牌
	tableScore		int32					//现在桌面上未被赢取的分数
	//end playingGameData, reset when start playing game

	roomOperateCh		chan *Operate
	daduCh		[]chan *Operate				//打独
	dropCardCh	[]chan *Operate				//出牌
	guoCh		[]chan *Operate				//过牌

	stop bool
}

func NewRoom(id uint64, config *RoomConfig) *Room {
	room := &Room{
		id:			id,
		config:			config,
		players:		make([]*Player, 0),
		cardPool:		card.NewPool(),
		observers:		make([]RoomObserver, 0),
		roomStatus:		RoomStatusWaitAllPlayerEnter,
		playedGameCnt:	0,

		roomOperateCh: make(chan *Operate, 1024),
		daduCh: make([]chan *Operate, config.NeedPlayerNum),
		dropCardCh: make([]chan *Operate, config.NeedPlayerNum),
		guoCh: make([]chan *Operate, config.NeedPlayerNum),
	}
	for idx := 0; idx < int(config.NeedPlayerNum); idx ++ {
		room.daduCh[idx] = make(chan *Operate, 1)
		room.dropCardCh[idx] = make(chan *Operate, 1)
		room.guoCh[idx] = make(chan *Operate, 1)
	}
	return room
}

func (room *Room) GetId() uint64 {
	return room.id
}

func (room *Room) PlayerOperate(op *Operate) {
	pos := op.Operator.position
	log.Debug(time.Now().Unix(), room, op.Operator, "PlayerOperate", op.Op, " pos:", pos)

	switch op.Op {
	case OperateEnterRoom, OperateLeaveRoom:
		room.roomOperateCh <- op
	case OperateReadyRoom:
		room.roomOperateCh <- op
	case OperateConfirmDadu:
		room.daduCh[pos] <- op
	case OperateDrop:
		room.dropCardCh[pos] <- op
	case OperateGuo:
		room.guoCh[pos] <- op
	}
}

func (room *Room) addObserver(observer RoomObserver) {
	room.observers = append(room.observers, observer)
}

func (room *Room) Start() {
	go func() {
		start_time := time.Now().Unix()
		for  {
			if !room.stop {
				room.checkStatus()
				time.Sleep(time.Microsecond * 10)
			}else{
				break
			}
		}
		end_time := time.Now().Unix()
		log.Debug(end_time - start_time, "over^^")
	}()
}

func (room *Room) checkStatus() {
	switch room.roomStatus {
	case RoomStatusWaitAllPlayerEnter:
		room.waitAllPlayerEnter()
	case RoomStatusWaitAllPlayerReady:
		room.waitAllPlayerReady()
	case RoomStatusGameStart:
		room.gameStart()
	case RoomStatusPlayGame:
		room.playGame()
	case RoomStatusEndPlayGame:
		room.endPlayGame()
	case RoomStatusRoomEnd:
		room.close()
	}
}

func (room *Room) GetPlayerNum() int32 {
	return int32(len(room.players))
}

func (room *Room) isRoomEnd() bool {
	return room.playedGameCnt >= room.config.MaxPlayGameCnt
}

func (room *Room) close() {
	log.Debug(time.Now().Unix(), room, "Room.close")
	room.stop = true
	for _, observer := range room.observers {
		observer.OnRoomClosed(room)
	}

	for _, player := range room.players {
		player.OnRoomClosed()
	}
}

func (room *Room) isEnterPlayerEnough() bool {
	length := room.GetPlayerNum()
	log.Debug(time.Now().Unix(), room, "Room.isEnterPlayerEnough, player num :", length, ", need :", room.config.NeedPlayerNum)
	return length >= room.config.NeedPlayerNum
}

func (room *Room) switchStatus(status RoomStatusType) {
	log.Debug(time.Now().Unix(), room, "room status switch,", room.roomStatus, " =>", status)
	room.roomStatus = status
	log.Debug("---------------------------------------")
}

//等待游戏开局
func (room *Room) waitAllPlayerEnter() {
	log.Debug(time.Now().Unix(), room, "waitAllPlayerEnter......")
	breakTimerTime := time.Duration(0)
	timeout := time.Duration(room.config.WaitPlayerEnterRoomTimeout) * time.Second
	for {
		timer := timeout - breakTimerTime
		select {
		case <-time.After(timer):
			log.Debug(time.Now().Unix(), room, "waitAllPlayerEnter timeout", timeout)
			room.switchStatus(RoomStatusRoomEnd) //超时发现没有足够的玩家都进入房间了，则结束
			return
		case op := <-room.roomOperateCh:
			if op.Op == OperateEnterRoom || op.Op == OperateLeaveRoom || op.Op == OperateReadyRoom {
				log.Debug(time.Now().Unix(), room, "waitAllPlayerEnter catch operate:", op)
				room.dealPlayerOperate(op)
				if room.isAllPlayerEnter() {
					room.switchStatus(RoomStatusWaitAllPlayerReady)
					return
				}
			}
		}
	}
}

func (room *Room) isAllPlayerEnter() bool {
	length := len(room.players)
	log.Debug(room, "Room.isAllPlayerEnter, num:", length, "need:", room.config.NeedPlayerNum)
	if length < int(room.config.NeedPlayerNum) {
		return false
	}
	for _, player := range room.players{
		if !player.GetIsReady() {
			return false
		}
	}

	return true
}

func (room *Room) waitDropCard(player *Player, mustDrop bool) bool{
	for{
		select {
		case <- time.After(time.Second * room.config.WaitDropSec):
			random := util.RandomN(2)
			log.Debug(time.Now().Unix(), player, "waitDropCard do PlayerOperate, random:", random)

			if mustDrop || random == 0 {
			//if mustDrop {
				num := room.getRandomDropNum(player)
				tailCards := player.GetTailCard(num)
				dropCards := card.CopyCards(tailCards)
				data := &OperateDropData{dropCards}
				op := NewOperateDrop(player, data)
				room.dealPlayerOperate(op)
				return true
			}else{
				data := &OperateGuoData{}
				op := NewOperateGuo(player, data)
				room.dealPlayerOperate(op)
				return false
			}
		case op := <-room.dropCardCh[player.position]:
			log.Debug(time.Now().Unix(), player, "Player.waitDropCard:", op.Data)
			room.dealPlayerOperate(op)
			return true
		case op := <-room.guoCh[player.position] :
			log.Debug(room, "Room.waitDropCard operate :", op)
			room.dealPlayerOperate(op)
			return false
		}
	}

	log.Debug(time.Now().Unix(), player, "Player.waitBet fasle")
	return false
}

func (room *Room) getRandomDropNum(player *Player) int{
	num := room.config.RandomDropNum
	hand_cards_len := player.GetPlayingCards().CardsInHand.Len()
	if hand_cards_len < num{
		num = hand_cards_len
	}
	return num
}

func (room *Room) waitInitPlayerReady(player *Player) {
	time.Sleep(time.Second * room.config.WaitReadySec)
	if (room.roomStatus == RoomStatusWaitAllPlayerEnter || room.roomStatus == RoomStatusWaitAllPlayerReady) && !player.GetIsReady() {
		data := &OperateReadyRoomData{}
		op := NewOperateReadyRoom(player, data)
		log.Debug(player, "waitInitPlayerReady do PlayerOperate")
		room.PlayerOperate(op)
	}
}

//TODO 等待开始时有玩家需要等待两次的bug
func (room *Room) waitPlayerReady(player *Player) bool {
	log.Debug(time.Now().Unix(), player, "waitPlayerReady")
	for{
		select {
		case <- time.After(time.Second * room.config.WaitReadySec):
			data := &OperateReadyRoomData{}
			op := NewOperateReadyRoom(player, data)
			log.Debug("******")
			log.Debug(time.Now().Unix(), player, "waitPlayerReady do PlayerOperate")
			room.PlayerOperate(op)
			continue
		case op := <-room.roomOperateCh:
			log.Debug(time.Now().Unix(), player, "Player.waitPlayerReady")
			room.dealPlayerOperate(op)
			return true
		}
	}

	log.Debug(time.Now().Unix(), player, "Player.waitPlayerReady fasle")
	return false
}

func (room *Room) waitAllPlayerReady() {
	log.Debug(time.Now().Unix(), room, room.playedGameCnt, "Room.waitAllPlayerReady......")
	if room.playedGameCnt == 0 {
		room.switchStatus(RoomStatusGameStart)
		return
	}

	//等待所有玩家准备
	for _, player := range room.players {
		go room.waitPlayerReady(player)
	}
	breakTimerTime := time.Duration(0)
	timeout := time.Duration(room.config.WaitPlayerEnterRoomTimeout) * time.Second
	for {
		timer := timeout - breakTimerTime
		select {
		case <-time.After(timer):
			log.Debug(time.Now().Unix(), room, "waitAllPlayerReady timeout", timeout)
			room.switchStatus(RoomStatusRoomEnd) //超时发现没有足够的玩家都进入房间了，则结束
			return
		case op := <-room.roomOperateCh:
			if room.roomStatus == RoomStatusRoomEnd {
				//如果此时房间已经结束，则直接返回，房间结束
				log.Debug(time.Now().Unix(), "waitAllPlayerReady room.roomStatus == RoomStatusRoomEnd")
				return
			}
			if op.Op == OperateReadyRoom || op.Op == OperateLeaveRoom{
				log.Debug(time.Now().Unix(), room, "Room.waitAllPlayerReady catch operate:", op)
				room.dealPlayerOperate(op)
				if room.isAllPlayerReady() {
					room.switchStatus(RoomStatusGameStart)
					return
				}
			}
		}
	}
}

func (room *Room) gameStart() {
	log.Debug(time.Now().Unix(), room, "gameStart", room.playedGameCnt)

	// 重置牌池, 洗牌
	room.Reset()
	room.cardPool.ReGenerate()

	//发牌
	fanpai_seq := util.RandomN(card.TOTAL_CARD_NUM)
	room.masterPlayer, room.assistPlayer, room.fanpai = room.putCardsToPlayers(card.INIT_CARD_NUM, room.config.InitType, fanpai_seq)
	//room.curOperator = room.masterPlayer
	room.switchOperator(room.masterPlayer, true)
	log.Debug(time.Now().Unix(), "master", room.masterPlayer)
	log.Debug(time.Now().Unix(), "assist", room.assistPlayer)
	log.Debug(time.Now().Unix(), "fanpai", room.fanpai)

	//通知所有玩家手上的牌
	for _, player := range room.players {
		player.OnGetInitCards()
	}

	//等待玩家打独
	tmpPlayer := room.opMaster
	room.isDadu = false
	is_dadu := false
	for {
		log.Debug("******")
		wait_data := &WaitDaduMsgData{
			WaitDaduPlayer:tmpPlayer,
			LeftSec:int32(room.config.WaitDaduSec),
		}
		msg := NewWaitDaduMsg(tmpPlayer, wait_data)
		for _, player := range room.players {
			player.OnWaitWaitDadu(msg)
		}

		is_dadu = room.waitPlayerDadu(tmpPlayer)
		if is_dadu {
			room.isDadu = true
			room.daduPlayer = tmpPlayer
			break
		}

		tmpPlayer = room.nextPlayer(tmpPlayer)
		if tmpPlayer == room.opMaster{
			break
		}
	}

	//交换位置
	room.switchPosition()
	time.Sleep(time.Second * room.config.AfterSwitchPositionSleep)

	//切换状态，开始打牌
	room.switchStatus(RoomStatusPlayGame)
	log.Debug(time.Now().Unix(), room, "Room.playGame", room.playedGameCnt)

	//通知开始出牌
	sp_data := &StartPlayMsgData{
		IsDadu:room.isDadu,
		Assist:room.assistPlayer,
		Master:room.masterPlayer,
	}
	msg := NewStartPlayMsg(nil, sp_data)
	for _, player := range room.players {
		player.OnStartPlay(msg)
	}
}

//TODO 玩家打完牌之后切换另外一个玩家逻辑
func (room *Room) playGame() {
	//log.Debug(time.Now().Unix(), room, "Room.playGame", room.playedGameCnt)

	//room.switchOperator(room.masterPlayer, true)
	if room.opMaster.GetNeedDrop() {
		log.Debug("wait drop:", room.opMaster)
		room.waitDropCard(room.opMaster, true)
	}
	//room.opMaster.SetNeedDrop(false)
	tmpPlayer := room.opMaster
	is_round_end := false
	for {
		tmpPlayer = room.nextPlayer(tmpPlayer)
		if tmpPlayer == room.opMaster{
			//room.opMaster.SetNeedDrop(true)
			break
		}

		room.waitOperator = tmpPlayer
		is_drop := room.waitDropCard(tmpPlayer, false)
		if is_drop {
			//查看玩家是否出完手牌
			if room.isAllCardsDropped(tmpPlayer) {
				tmpPlayer.SetIsIsEndPlaying(true)
				is_round_end = room.isRoundEnd(tmpPlayer)
			}

			//room.curOperator = tmpPlayer
			room.switchOperator(tmpPlayer, false)
			break
		}
	}

	if is_round_end {
		room.switchStatus(RoomStatusEndPlayGame)
		//通知开始出牌
		jiesuan_data := &JiesuanMsgData{}
		msg := NewJiesuanMsg(nil, jiesuan_data)
		for _, player := range room.players {
			player.OnJiesuan(msg)
		}
	}
}

func (room *Room) isAllCardsDropped(player * Player) bool{
	return player.GetLeftCardNum() == 0
}

//在一个玩家出完手牌时判断此局是否已经结束
func (room *Room) isRoundEnd(endPlayingPlayer * Player) bool{
	if room.isDadu {
		return true
	}

	//只要任意一方有两个人打完就可以结束
	master := room.masterPlayer
	assist := room.assistPlayer
	if endPlayingPlayer == master {
		return assist.GetIsEndPlaying()
	}else if endPlayingPlayer == assist {
		return master.GetIsEndPlaying()
	}else {
		for _, room_player := range room.players {
			if room_player != master && room_player != assist && room_player != endPlayingPlayer{
				return room_player.GetIsEndPlaying()
			}
		}
	}
	return false
}

func (room *Room) waitPlayerDadu(player *Player) bool {
	for{
		select {
		case <- time.After(time.Second * room.config.WaitDaduSec):
			data := &OperateConfirmDaduData{IsDadu:false}
			op := NewOperateConfirmDadu(player, data)
			log.Debug(time.Now().Unix(), player, "waitPlayerDadu do PlayerOperate")
			room.PlayerOperate(op)
			continue
		case op := <-room.daduCh[player.position]:
			if dadu_data, ok := op.Data.(*OperateConfirmDaduData); ok {
				log.Debug(time.Now().Unix(), player, "Player.waitPlayerDadu:", op.Data)
				room.dealPlayerOperate(op)
				return dadu_data.IsDadu
			}
		}
	}

	log.Debug(time.Now().Unix(), player, "Player.waitPlayerDadu fasle")
	return false
}

func (room *Room) switchPosition() {
	log.Debug(time.Now().Unix(), "switchPosition")
	//打独不需要交换位置
	if room.isDadu{
		return
	}

	master := room.masterPlayer
	assist := room.assistPlayer
	master_pos := master.GetPosition()
	assist_pos := assist.GetPosition()
	//庄家同时摸到两张牌，则与对家一起打
	if master == assist {
		//log.Debug(time.Now().Unix(), "switchPosition master == assist")
		assist_pos := (master_pos + 2) % room.config.NeedPlayerNum
		room.assistPlayer = room.getPlayerByPos(assist_pos)
		return
	}

	//如果两人已经是对家，则不需要交换位置
	if (master_pos + assist_pos) % 2 == 0 {
		//log.Debug(time.Now().Unix(), "switchPosition already opposite")
		return
	}

	//对家与assist交换位置
	//log.Debug(time.Now().Unix(), "switchPosition switch")
	sw_data := &SwitchPositionMsgData{
		AssistUid:assist.GetId(),
		AssistPos:assist.GetPosition(),
	}
	opp_pos := (master_pos + 2) % room.config.NeedPlayerNum
	opp_player := room.getPlayerByPos(opp_pos)
	sw_data.OppUid = opp_player.GetId()
	sw_data.OppPos = opp_player.GetPosition()
	msg := NewSwitchPositionMsg(assist, sw_data)
	opp_player.position = assist.GetPosition()
	assist.position = sw_data.OppPos

	for _, player := range room.players {
		player.OnSwitchPosition(msg)
	}
}

func (room *Room) switchOperator(player *Player, mustDrop bool) {
	log.Debug(time.Now().Unix(), room, "switchOperator", room.opMaster, "=>", player)
	room.opMaster = player
	player.SetNeedDrop(mustDrop)

	op := room.makeSwitchOperatorOperate(player, mustDrop)
	for _, player := range room.players {
		player.OnPlayerSuccessOperated(op)
	}
}


func (room *Room) makeSwitchOperatorOperate(operator *Player, mustDrop bool) *Operate {
	return NewSwitchOperator(operator, &OperateSwitchOperatorData{MustDrop:mustDrop})
}

func (room *Room) endPlayGame() {
	room.playedGameCnt++
	log.Debug(time.Now().Unix(), room, "Room.endPlayGame cnt :", room.playedGameCnt)
	if room.isRoomEnd() {
		//log.Debug(time.Now().Unix(), room, "Room.endPlayGame room end")
		room.switchStatus(RoomStatusRoomEnd)
	} else {
		for _, player := range room.players {
			player.OnEndPlayGame()
		}
		//log.Debug(time.Now().Unix(), room, "Room.endPlayGame restart play game")
		room.switchStatus(RoomStatusWaitAllPlayerReady)
		log.Debug("============================================================================")
	}
}

func (room *Room) jiesuan() *Message {
	master_player := room.masterPlayer
	master_jiesuan_data := &PlayerJiesuanData{
		P:master_player,
		Score:0,
	}

	data := &JiesuanMsgData{}
	data.Scores = make([]*PlayerJiesuanData, 0)
	for _, player := range room.players {
		if player != master_player {
			player_jiesuan_data := &PlayerJiesuanData{
				P:player,
				Score:0,
			}
			data.Scores = append(data.Scores, player_jiesuan_data)
		}
	}
	data.Scores = append(data.Scores, master_jiesuan_data)
	return NewJiesuanMsg(nil, data)
}

//取指定玩家的下一个玩家
func (room *Room) getPlayerByPos(position int32) *Player {
	for _, room_player := range room.players {
		if room_player.GetPosition() == position {
			return room_player
		}
	}
	if room.GetPlayerNum() > 0 {
		return room.players[0]
	}
	return nil
}

//取指定玩家的下一个玩家
func (room *Room) nextPlayer(player *Player) *Player {
	pos := player.GetPosition()

	need_player_num := int32(room.config.NeedPlayerNum)
	for next_pos := pos + 1; next_pos < need_player_num; next_pos++ {
		for _, room_player := range room.players {
			if room_player.GetPosition() == next_pos && !room_player.GetIsEndPlaying(){
				//log.Debug(time.Now().Unix(), "nextPlayer", "pos:", pos, "next_pos:", next_pos)
				return room_player
			}
		}
	}

	for next_pos := int32(0); next_pos < pos; next_pos++ {
		for _, room_player := range room.players {
			if room_player.GetPosition() == next_pos && !room_player.GetIsEndPlaying(){
				//log.Debug(time.Now().Unix(), "nextPlayer", "pos:", pos, "next_pos:", next_pos)
				return room_player
			}
		}
	}

	//log.Debug(time.Now().Unix(), "nextPlayer", "pos:", pos, "next_pos:", 0)
	return room.players[0]
}

func (room *Room) isAllPlayerReady() bool{
	for _, player := range room.players {
		if !player.isReady {
			return false
		}
	}
	return true
}

//处理玩家操作
func (room *Room) dealPlayerOperate(op *Operate) bool{
	//log_time := time.Now().Unix()
	//log.Debug(log_time, room, "Room.dealPlayerOperate :", op)
	switch op.Op {
	case OperateEnterRoom:
		if _, ok := op.Data.(*OperateEnterRoomData); ok {
			if room.addPlayer(op.Operator) {
				//玩家进入成功
				player_pos := room.getMinUsablePosition()
				op.Operator.EnterRoom(room, player_pos)
				//log.Debug(log_time, room, "Room.dealPlayerOperate player enter :", op.Operator)
				op.ResultCh <- true
				room.broadcastPlayerSuccessOperated(op)
				return true
			}
		}

	case OperateReadyRoom:
		if _, ok := op.Data.(*OperateReadyRoomData); ok {
			if room.readyPlayer(op.Operator) { //	玩家确认开始游戏
				op.Operator.ReadyRoom(room)
				//log.Debug(log_time, room, "Room.dealPlayerOperate player ready :", op.Operator)
				op.ResultCh <- true
				room.broadcastPlayerSuccessOperated(op)
				return true
			}
		}

	case OperateLeaveRoom:
		if _, ok := op.Data.(*OperateLeaveRoomData); ok {
			//log.Debug(log_time, room, "Room.dealPlayerOperate player leave :", op.Operator)
			room.delPlayer(op.Operator)
			op.Operator.LeaveRoom()
			op.ResultCh <- true
			room.broadcastPlayerSuccessOperated(op)
			return true
		}

	case OperateConfirmDadu:
		if dadu_data, ok := op.Data.(*OperateConfirmDaduData); ok {
			//log.Debug(log_time, room, "Room.dealPlayerOperate player dadu :", op.Operator)
			op.Operator.Dadu(dadu_data.IsDadu)
			op.ResultCh <- true
			room.broadcastPlayerSuccessOperated(op)
			return true
		}

	case OperateDrop:
		if drop_data, ok := op.Data.(*OperateDropData); ok {
			if op.Operator.Drop(drop_data.whatGroup) {
				//出牌
				log.Debug(time.Now().Unix(), room, "Room.dealPlayerOperate player drop :", op.Operator)
				op.ResultCh <- true
				room.broadcastPlayerSuccessOperated(op)
				return true
			}
		}

	case OperateGuo:
		if _, ok := op.Data.(*OperateGuoData); ok {
			log.Debug(time.Now().Unix(), room, "Room.dealPlayerOperate player guo :", op.Operator)
			op.ResultCh <- true
			room.broadcastPlayerSuccessOperated(op)
			return true
		}

	}
	op.ResultCh <- false
	return false
}

//查找房间中未被占用的最新的position
func (room *Room) getMinUsablePosition() (int32)  {
	//log.Debug(time.Now().Unix(), room, "getMinUsablePosition")
	//获取所有已经被占用的position
	player_positions := make([]int32, 0)
	for _, room_player := range room.players {
		player_positions = append(player_positions, room_player.GetPosition())
	}

	//查找未被占用的position中最小的
	room_max_position := int32(room.config.NeedPlayerNum - 1)
	for i := int32(0); i <= room_max_position; i++ {
		is_occupied := false
		for _, occupied_pos := range player_positions{
			if occupied_pos == i {
				is_occupied = true
				break
			}
		}
		if !is_occupied {
			return i
		}
	}
	return room_max_position
}

//给所有玩家发牌
func (room *Room) putCardsToPlayers(init_num int, init_type int32, fanpai_seq int) (master, assist *Player, fanpai *card.Card) {
	log.Debug(time.Now().Unix(), room, "Room.putCardsToPlayers, init_type:", init_type)
	master, assist = nil, nil
	fanpai = nil

	pool_card_num := room.cardPool.GetCardNum()
	if fanpai_seq >= pool_card_num {
		return
	}
	fanpai = room.cardPool.At(fanpai_seq)

	if init_type == 1 {
		for num := 0; num < init_num; num++ {
			for _, player := range room.players {
				put_card := room.putCardToPlayer(player)
				if fanpai.SameAs(put_card) {
					master = player
				}else if fanpai.SameCardTypeNoAs(put_card) {
					assist = player
				}
			}
		}
	}else {
		//一次发多张牌
		for round := 0; round < card.TYPE2_ROUND_TIMES; round++ {
			for _, player := range room.players {
				for round := 0; round < card.TYPE2_ROUND_CARD_NUM; round++ {
					put_card := room.putCardToPlayer(player)
					if fanpai.SameAs(put_card) {
						master = player
					}else if fanpai.SameCardTypeNoAs(put_card) {
						assist = player
					}
				}
			}
		}
		for _, player := range room.players {
			for round := 0; round < card.TYPE2_LAST_ROUND_CARD_NUM; round++ {
				put_card := room.putCardToPlayer(player)
				if fanpai.SameAs(put_card) {
					master = player
				}else if fanpai.SameCardTypeNoAs(put_card) {
					assist = player
				}
			}
		}
	}
	return
}

//添加玩家
func (room *Room) addPlayer(player *Player) bool {
	/*if room.roomStatus != RoomStatusWaitAllPlayerEnter {
		return false
	}*/
	if room.GetPlayerNum() >= room.config.NeedPlayerNum {
		return false
	}
	room.players = append(room.players, player)
	return true
}

func (room *Room) readyPlayer(player *Player) bool {
	if room.roomStatus != RoomStatusWaitAllPlayerEnter && room.roomStatus != RoomStatusWaitAllPlayerReady{
		return false
	}
	player.SetIsReady(true)
	return true
}

func (room *Room) delPlayer(player *Player)  {
	for idx, p := range room.players {
		if p == player {
			room.players = append(room.players[0:idx], room.players[idx+1:]...)
			return
		}
	}
}

func (room *Room) broadcastPlayerSuccessOperated(op *Operate) {
	//log.Debug(time.Now().Unix(), room, "Room.broadcastPlayerSucOp :", op)
	for _, player := range room.players {
		player.OnPlayerSuccessOperated(op)
	}
}

//发牌给指定玩家
func (room *Room) putCardToPlayer(player *Player) *card.Card {
	card := room.cardPool.PopFront()
	if card == nil {
		return nil
	}
	player.AddCard(card)
	return card
}

func (room *Room) Reset()  {
	room.opMaster = nil
	room.waitOperator = nil
	room.masterPlayer = nil
	room.assistPlayer = nil
	room.daduPlayer = nil
	room.isDadu = false
	room.fanpai = nil
	room.tableScore = 0
}

func (room *Room) String() string {
	if room == nil {
		return "{room=nil}"
	}
	return fmt.Sprintf("{room=%v}", room.GetId())
}

func (room *Room) clearChannel() {
	for idx := 0 ; idx < int(room.config.NeedPlayerNum); idx ++ {
		select {
		case op := <-room.daduCh[idx]:
			op.ResultCh <- false
		default:
		}

		select {
		case op := <-room.dropCardCh[idx]:
			op.ResultCh <- false
		default:
		}

		select {
		case op := <-room.guoCh[idx]:
			op.ResultCh <- false
		default:
		}
	}
}
