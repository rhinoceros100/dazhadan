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
	RoomStatusPlayGame					// 正在进行游戏，结束后会进入RoomStatusEndPlayGame
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
	creatorUid		uint64					//创建房间玩家的uid
	opMaster		*Player					//当前出过牌的玩家
	waitOperator		*Player					//等待出牌的玩家
	masterPlayer 		*Player					//庄
	nextOpMaster		*Player					//下一个出牌玩家
	assistPlayer 		*Player					//和庄一
	isPlayAlone		bool					//是否打独
	turnCard		*card.Card				//翻出的牌
	tableScore		int32					//现在桌面上未被赢取的分数
	endPlayingNum		int32					//本局已经出完牌的玩家数量

	cardsType 		int					//上一次出的牌型
	planeNum 		int					//飞机数量
	weight	 		int					//权重
	//end playingGameData, reset when start playing game

	roomOperateCh	chan *Operate
	playAloneCh	[]chan *Operate				//打独
	dropCardCh	[]chan *Operate				//出牌
	passCh		[]chan *Operate				//过牌
	roomReadyCh	[]chan *Operate

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
		creatorUid:		0,
		tableScore:		0,
		playedGameCnt:		0,
		endPlayingNum:		0,
		cardsType:		0,
		planeNum:		0,
		weight:			0,
		opMaster:		 nil,
		waitOperator:		 nil,
		masterPlayer:		 nil,
		nextOpMaster:		 nil,
		assistPlayer:		 nil,

		roomOperateCh: make(chan *Operate, 1024),
		playAloneCh: make([]chan *Operate, config.NeedPlayerNum),
		dropCardCh: make([]chan *Operate, config.NeedPlayerNum),
		passCh: make([]chan *Operate, config.NeedPlayerNum),
		roomReadyCh: make([]chan *Operate, config.NeedPlayerNum),
	}
	for idx := 0; idx < int(config.NeedPlayerNum); idx ++ {
		room.playAloneCh[idx] = make(chan *Operate, 1)
		room.dropCardCh[idx] = make(chan *Operate, 1)
		room.passCh[idx] = make(chan *Operate, 1)
		room.roomReadyCh[idx] = make(chan *Operate, 1)
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
		if room.roomStatus == RoomStatusWaitAllPlayerEnter {
			room.roomOperateCh <- op
		}else {
			room.roomReadyCh[pos] <- op
		}
	case OperateConfirmPlayAlone:
		room.playAloneCh[pos] <- op
	case OperateDrop:
		room.dropCardCh[pos] <- op
	case OperatePass:
		room.passCh[pos] <- op
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

func (room *Room) GetTableScore() int32 {
	return room.tableScore
}

func (room *Room) AddTableScore(score int32) {
	room.tableScore += score
}

func (room *Room) ResetTableScore() {
	room.tableScore = 0
}

func (room *Room) GetEndPlayingNum() int32 {
	return room.endPlayingNum
}

func (room *Room) IncEndPlayingNum() int32 {
	room.endPlayingNum += 1
	return room.endPlayingNum
}

func (room *Room) GetCardsType() int {
	return room.cardsType
}

func (room *Room) SetCardsType(cardsType int) {
	room.cardsType = cardsType
}

func (room *Room) GetPlaneNum() int {
	return room.planeNum
}

func (room *Room) SetPlaneNum(planeNum int) {
	room.planeNum = planeNum
}

func (room *Room) GetWeight() int {
	return room.weight
}

func (room *Room) SetWeight(weight int) {
	room.weight = weight
}

func (room *Room) close() {
	log.Debug(time.Now().Unix(), room, "Room.close")
	room.stop = true
	for _, observer := range room.observers {
		observer.OnRoomClosed(room)
	}

	msg := room.totalSummary()
	for _, player := range room.players {
		player.OnRoomClosed(msg)
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

func (room *Room) canCover(cardsType, planeNum, weight int) (canCover bool) {
	canCover = false
	if room.GetCardsType() == card.CardsType_NO {
		return cardsType != card.CardsType_NO
	}
	//已经出的牌型非炸弹牌型
	if room.GetCardsType() < 20{
		if cardsType > 20 {
			return true
		}
		//普通牌型打普通牌型必须为同一牌型，并且飞机数量必须相同
		if cardsType != room.GetCardsType(){
			return false
		}
		if cardsType == card.CardsType_STAIGHT || cardsType == card.CardsType_PAIRS || cardsType >= 11 {
			if planeNum != room.GetPlaneNum() {
				return false
			}
		}
		return weight > room.GetWeight()
	}

	//更大的炸弹可以管住
	if cardsType > room.GetCardsType(){
		return true
	}
	return weight > room.GetWeight()
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

func (room *Room) waitDropCard(player *Player, mustDrop bool, canDrop bool) bool{
	wait_time := room.config.WaitDropSec
	if !canDrop{
		wait_time = time.Duration(2)
	}
	for{
		select {
		case <- time.After(time.Second * wait_time):
			random := util.RandomN(2)
			log.Debug(time.Now().Unix(), player, "waitDropCard do PlayerOperate, random:", random)

			if mustDrop || random == 0 {
			//if mustDrop {
				room.nextOpMaster = nil
				num := room.getRandomDropNum(player)
				tailCards := player.GetTailCard(num)
				dropCards := card.CopyCards(tailCards)

				cards_num := player.playingCards.CardsInHand.Len()
				is_last_cards := false
				if cards_num == len(dropCards) {
					is_last_cards = true
				}
				drop_cards := card.CreateNewCards(dropCards)

				data := &OperateDropData{whatGroup:dropCards}
				data.cardsType, data.planeNum, data.weight = card.GetCardsType(drop_cards, is_last_cards, 0, 0)
				can_cover := room.canCover(data.cardsType, data.planeNum, data.weight)
				log.Debug("******can_cover:", can_cover)
				op := NewOperateDrop(player, data)
				room.dealPlayerOperate(op)
				return true
			}else{
				data := &OperatePassData{}
				op := NewOperatePass(player, data)
				room.dealPlayerOperate(op)
				return false
			}
		case op := <-room.dropCardCh[player.position]:
			log.Debug(time.Now().Unix(), player, "Player.waitDropCard:", op.Data)
			room.nextOpMaster = nil
			room.dealPlayerOperate(op)
			return true
		case op := <-room.passCh[player.position] :
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
		case op := <-room.roomReadyCh[player.GetPosition()]:
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
	timeout := int64(room.config.WaitPlayerOperateTimeout)
	start_time := time.Now().Unix()
	for  {
		//房间结束
		if room.roomStatus == RoomStatusRoomEnd {
			//如果此时房间已经结束，则直接返回，房间结束
			log.Debug(time.Now().Unix(), "waitAllPlayerReady room.roomStatus == RoomStatusRoomEnd")
			return
		}

		//所有人都已准备
		if room.isAllPlayerReady() {
			room.switchStatus(RoomStatusGameStart)
			return
		}

		//超时结束
		time_now := time.Now().Unix()
		if start_time + timeout < time_now {
			log.Debug(time.Now().Unix(), room, "waitAllPlayerReady timeout", timeout)
			room.switchStatus(RoomStatusRoomEnd) //超时发现没有足够的玩家都进入房间了，则结束
			return
		}

		time.Sleep(time.Millisecond * 2)
	}

	/*for {
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
	}*/
}

func (room *Room) gameStart() {
	log.Debug(time.Now().Unix(), room, "gameStart", room.playedGameCnt)

	//测试八炸数量
	/*bomb4, bomb5, bomb6, bomb7, bomb8, bomb_joker := 0, 0, 0, 0, 0, 0
	for i := 0; i < 100000; i++ {
		room.cardPool.ReGenerate()
		room.putCardsToPlayers(card.INIT_CARD_NUM, room.config.InitType, 5)
		for _, player := range room.players {
			num4, num5, num6, num7, num8, num_joker := player.GetBomb8Num()
			bomb4 += num4
			bomb5 += num5
			bomb6 += num6
			bomb7 += num7
			bomb8 += num8
			bomb_joker += num_joker
			player.Reset()
		}
		//log.Debug("i:", i, " bomb4:", bomb4, " bomb5:", bomb5, " bomb6:", bomb6, " bomb7:", bomb7, " bomb8:", bomb8, " bomb_joker:", bomb_joker)
		//log.Debug("====================")
	}
	log.Debug(time.Now().Unix(), "bomb4:", bomb4, " bomb5:", bomb5, " bomb6:", bomb6, " bomb7:", bomb7, " bomb8:", bomb8, " bomb_joker:", bomb_joker)
	time.Sleep(time.Second * 1000)*/

	// 重置牌池, 洗牌
	room.Reset()
	room.cardPool.ReGenerate()

	//发牌
	fanpai_seq := util.RandomN(card.TOTAL_CARD_NUM)
	room.masterPlayer, room.assistPlayer, room.turnCard = room.putCardsToPlayers(card.INIT_CARD_NUM, room.config.InitType, fanpai_seq)
	//room.curOperator = room.masterPlayer
	room.switchOpMaster(room.masterPlayer, true, true, false)
	log.Debug(time.Now().Unix(), "master", room.masterPlayer)
	log.Debug(time.Now().Unix(), "assist", room.assistPlayer)
	log.Debug(time.Now().Unix(), "turnCard", room.turnCard)

	//通知所有玩家手上的牌
	for _, player := range room.players {
		player.OnGetInitCards()
	}

	//等待玩家打独
	tmpPlayer := room.opMaster
	room.isPlayAlone = false
	is_play_alone := false
	for {
		log.Debug("******")
		wait_data := &WaitPlayAloneMsgData{
			WaitPlayAlonePlayer:tmpPlayer,
			LeftSec:int32(room.config.WaitPlayAloneSec),
		}
		msg := NewWaitPlayAloneMsg(tmpPlayer, wait_data)
		for _, player := range room.players {
			player.OnWaitPlayAlone(msg)
		}

		room.switchWaitPlayer(tmpPlayer, false, false, false)
		is_play_alone = room.waitPlayerPlayAlone(tmpPlayer)
		if is_play_alone {
			room.isPlayAlone = true
			room.masterPlayer = tmpPlayer
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
		IsPlayAlone:room.isPlayAlone,
		Assist:room.assistPlayer,
		Master:room.masterPlayer,
	}
	msg := NewStartPlayMsg(nil, sp_data)
	for _, player := range room.players {
		player.OnStartPlay(msg)
	}
}

func (room *Room) playGame() {
	log.Debug(time.Now().Unix(), room, "Room.playGame", room.playedGameCnt)

	is_round_end := false
	if room.opMaster.GetNeedDrop() {
		room.switchWaitPlayer(room.opMaster, true, true, true)
		room.waitDropCard(room.opMaster, true, true)

		//查看玩家是否出完手牌
		if room.isAllCardsDropped(room.opMaster) {
			room.opMaster.SetIsEndPlaying(true)
			rank := room.IncEndPlayingNum()
			room.opMaster.SetRank(rank)

			is_round_end = room.isRoundEnd(room.opMaster)
			if !is_round_end {
				room.nextOpMaster = room.oppositePlayer(room.opMaster)
			}else{
				//结束时赢取牌桌上的分数
				table_score := room.GetTableScore()
				room.opMaster.AddScore(table_score)
				room.ResetTableScore()
			}
		}else {
			room.opMaster.SetNeedDrop(false)
		}
	}

	if is_round_end {
		room.switchStatus(RoomStatusEndPlayGame)
		//通知开始出牌
		msg := room.summary()
		for _, player := range room.players {
			player.OnSummary(msg)
		}
		return
	}

	tmpPlayer := room.opMaster
	for {
		tmpPlayer = room.nextPlayer(tmpPlayer)
		if tmpPlayer == room.opMaster{
			//将桌上分数收入囊中
			table_score := room.GetTableScore()
			tmpPlayer.AddScore(table_score)
			room.ResetTableScore()
			//重置已出牌型
			room.SetCardsType(card.CardsType_NO)
			room.SetPlaneNum(0)
			room.SetWeight(0)

			if nil != room.nextOpMaster {
				room.switchOpMaster(room.nextOpMaster, true, true, false)
				room.nextOpMaster = nil
			}else {
				room.opMaster.SetNeedDrop(true)
			}
			break
		}

		canDrop := tmpPlayer.GetCanDrop()
		room.switchWaitPlayer(tmpPlayer, false, canDrop, true)
		is_drop := room.waitDropCard(tmpPlayer, false, canDrop)
		if is_drop {
			//查看玩家是否出完手牌
			if room.isAllCardsDropped(tmpPlayer) {
				tmpPlayer.SetIsEndPlaying(true)
				rank := room.IncEndPlayingNum()
				tmpPlayer.SetRank(rank)

				is_round_end = room.isRoundEnd(tmpPlayer)
				if !is_round_end {
					room.nextOpMaster = room.oppositePlayer(tmpPlayer)
				}else{
					//结束时赢取牌桌上的分数
					table_score := room.GetTableScore()
					tmpPlayer.AddScore(table_score)
					room.ResetTableScore()
				}
			}

			room.switchOpMaster(tmpPlayer, false, true, false)
			break
		}
	}

	if is_round_end {
		room.switchStatus(RoomStatusEndPlayGame)
		//通知开始出牌
		msg := room.summary()
		for _, player := range room.players {
			player.OnSummary(msg)
		}
	}
}

func (room *Room) isAllCardsDropped(player * Player) bool{
	return player.GetLeftCardNum() == 0
}

//在一个玩家出完手牌时判断此局是否已经结束
func (room *Room) isRoundEnd(endPlayingPlayer * Player) bool{
	if room.isPlayAlone {
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

func (room *Room) waitPlayerPlayAlone(player *Player) bool {
	for{
		select {
		case <- time.After(time.Second * room.config.WaitPlayAloneSec):
			data := &OperateConfirmPlayAloneData{IsPlayAlone:false}
			op := NewOperateConfirmPlayAlone(player, data)
			log.Debug(time.Now().Unix(), player, "waitPlayerPlayAlone do PlayerOperate")
			room.PlayerOperate(op)
			continue
		case op := <-room.playAloneCh[player.position]:
			if pa_data, ok := op.Data.(*OperateConfirmPlayAloneData); ok {
				log.Debug(time.Now().Unix(), player, "Player.waitPlayerPlayAlone:", op.Data)
				room.dealPlayerOperate(op)
				return pa_data.IsPlayAlone
			}
		}
	}

	log.Debug(time.Now().Unix(), player, "Player.waitPlayerPlayAlone fasle")
	return false
}

func (room *Room) switchPosition() {
	log.Debug(time.Now().Unix(), "switchPosition")
	//打独不需要交换位置
	if room.isPlayAlone{
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


func (room *Room) switchOpMaster(player *Player, mustDrop bool, canDrop bool, needNotify bool) {
	log.Debug(time.Now().Unix(), room, "switchOperator", room.opMaster, "=>", player)
	room.opMaster = player
	player.SetNeedDrop(mustDrop)

	if needNotify {
		op := room.makeSwitchOperatorOperate(player, mustDrop, canDrop)
		for _, player := range room.players {
			player.OnPlayerSuccessOperated(op)
		}
	}
}

func (room *Room) switchWaitPlayer(player *Player, mustDrop bool, canDrop bool, needNotify bool) {
	log.Debug(time.Now().Unix(), room, "switchWaitPlayer", room.waitOperator, "=>", player)
	room.waitOperator = player

	if needNotify {
		op := room.makeSwitchOperatorOperate(player, mustDrop, canDrop)
		for _, player := range room.players {
			player.OnPlayerSuccessOperated(op)
		}
	}
}

func (room *Room) makeSwitchOperatorOperate(operator *Player, mustDrop bool, canDrop bool) *Operate {
	return NewSwitchOperator(operator, &OperateSwitchOperatorData{
		MustDrop:mustDrop,
		CanDrop:canDrop,
	})
}

/*func (room *Room) switchOperator(player *Player, mustDrop bool) {
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
}*/

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

func (room *Room) summary() *Message {
	info_type := card.InfoType_Normal
	master_player := room.masterPlayer
	assist_player := room.assistPlayer
	//进行胜负结算
	if room.isPlayAlone {
		//打独成功
		if master_player.GetRank() == 1 {
			info_type = card.InfoType_PlayAloneSucc
			master_player.AddCoin(room.config.PlayAloneCoin * 3)
			master_player.SetIsWin(true)
			for _, player := range room.players {
				if player != master_player {
					player.AddCoin(0-room.config.PlayAloneCoin)
					player.SetIsWin(false)
				}
			}
		}else{ //打独失败
			info_type = card.InfoType_PlayAloneFail
			master_player.AddCoin(0-room.config.PlayAloneCoin * 3)
			master_player.SetIsWin(false)
			for _, player := range room.players {
				if player != master_player {
					player.AddCoin(room.config.PlayAloneCoin)
					player.SetIsWin(true)
				}
			}
		}
	}else{
		//双基
		if master_player.GetRank() + assist_player.GetRank() == 3 {
			info_type = card.InfoType_Shuangji
			master_player.AddCoin(room.config.ShuangjiCoin)
			master_player.SetIsWin(true)
			assist_player.AddCoin(room.config.ShuangjiCoin)
			assist_player.SetIsWin(true)
			for _, player := range room.players {
				if player != master_player && player != assist_player {
					player.AddCoin(0-room.config.ShuangjiCoin)
					player.SetIsWin(false)
				}
			}
		}else if master_player.GetRank() + assist_player.GetRank() == 0 {
			info_type = card.InfoType_Shuangji
			master_player.AddCoin(0-room.config.ShuangjiCoin)
			master_player.SetIsWin(false)
			assist_player.AddCoin(0-room.config.ShuangjiCoin)
			assist_player.SetIsWin(false)
			for _, player := range room.players {
				if player != master_player && player != assist_player {
					player.AddCoin(room.config.ShuangjiCoin)
					player.SetIsWin(true)
				}
			}
		}else{//普通打完，通过分数计算输赢
			for _, player := range room.players {
				if player.GetRank() == 0 {
					player.AddScore(-50)
				}
				if player.GetRank() == 1 {
					player.AddScore(50)
				}
			}
			master_side_scores := master_player.GetScore() + assist_player.GetScore()
			common_side_scores := int32(0)
			for _, player := range room.players {
				if player != master_player && player != assist_player {
					common_side_scores += player.GetScore()
				}
			}
			if master_side_scores > common_side_scores {
				master_player.AddCoin(room.config.WinCoin)
				master_player.SetIsWin(true)
				assist_player.AddCoin(room.config.WinCoin)
				assist_player.SetIsWin(true)
				for _, player := range room.players {
					if player != master_player && player != assist_player {
						player.AddCoin(0-room.config.WinCoin)
						player.SetIsWin(false)
					}
				}
			}else if master_side_scores < common_side_scores {
				master_player.AddCoin(0-room.config.WinCoin)
				master_player.SetIsWin(false)
				assist_player.AddCoin(0-room.config.WinCoin)
				assist_player.SetIsWin(false)
				for _, player := range room.players {
					if player != master_player && player != assist_player {
						player.AddCoin(room.config.WinCoin)
						player.SetIsWin(true)
					}
				}
			}
		}
	}

	//遍历查看玩家总奖数量
	var total_prize_num, prize_multiple int32 = 0, 1
	for _, prize_player := range room.players {
		total_prize_num += prize_player.GetPrize()
	}
	if total_prize_num == 1 && room.config.HaveDujiangDouble{
		prize_multiple = 2
		log.Debug("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@prize_multiple = 2")
	}
	//算奖金
	for _, prize_player := range room.players {
		prize_num := prize_player.GetPrize()
		if prize_num > 0 {
			prize_player.AddCoin(room.config.PrizeCoin * prize_multiple * 3 * prize_num)
			prize_player.AddPrizeCoin(room.config.PrizeCoin * prize_multiple * 3 * prize_num)

			for _, player := range room.players {
				if player != prize_player{
					player.AddCoin((0 - room.config.PrizeCoin * prize_multiple) * prize_num)
					player.AddPrizeCoin((0 - room.config.PrizeCoin * prize_multiple) * prize_num)
				}
			}
		}
	}

	data := &SummaryMsgData{}
	data.Scores = make([]*PlayerSummaryData, 0)
	for _, player := range room.players {
		player_summary_data := &PlayerSummaryData{
			P:player,
			Rank:player.GetRank(),
			Coin:player.GetCoin(),
			Score:player.GetScore(),
			Prize:player.GetPrize(),
			TotalCoin:player.GetTotalCoin(),
			PrizeCoin:player.GetPrizeCoin(),
			IsWin:player.GetIsWin(),
		}
		data.Scores = append(data.Scores, player_summary_data)

		player.AddTotalPrize(player.GetPrize())
		if player.GetIsWin() {
			player.IncWinNum()
			if info_type == card.InfoType_Shuangji{
				player.IncShuangjiNum()
			}
			if info_type == card.InfoType_PlayAloneSucc || info_type == card.InfoType_PlayAloneFail{
				player.IncPaSuccNum()
			}
		}
	}
	data.InfoType = info_type
	return NewSummaryMsg(nil, data)
}

func (room *Room) totalSummary() *Message {
	var max_win, max_lose int32 = 0, 0
	for _, player := range room.players {
		total_coin := player.GetTotalCoin()
		if total_coin > max_win {
			max_win = total_coin
		}
		if total_coin < max_lose {
			max_lose = total_coin
		}
	}

	data := &RoomClosedMsgData{}
	data.Summaries = make([]*TotalSummaryData, 0)
	for _, player := range room.players {
		summary_data := &TotalSummaryData{
			P:player,
			WinNum:player.GetWinNum(),
			ShuangjiNum:player.GetShuangjiNum(),
			PaSuccNum:player.GetPaSuccNum(),
			TotalPrize:player.GetTotalPrize(),
			TotalCoin:player.GetTotalCoin(),
			IsWinner:false,
			IsMostWinner:false,
			IsMostLoser:false,
			IsCreator:false,
		}
		if summary_data.TotalCoin > 0 {
			summary_data.IsWinner = true
		}
		if player.GetId() == room.creatorUid {
			summary_data.IsCreator = true
		}
		if summary_data.TotalCoin == max_lose {
			summary_data.IsMostLoser = true
		}
		if summary_data.TotalCoin == max_win {
			summary_data.IsMostWinner = true
		}
		data.Summaries = append(data.Summaries, summary_data)
	}
	return NewRoomClosedMsg(nil, data)
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
	for i := int32(1); i <= need_player_num; i++ {
		next_pos := (pos + i) % need_player_num
		for _, room_player := range room.players {
			if room_player.GetPosition() == next_pos{
				if room_player == room.opMaster {
					//log.Debug(time.Now().Unix(), "nextPlayer", "pos:", pos, "next_pos:", next_pos)
					return room_player
				}
				if !room_player.GetIsEndPlaying() {
					//log.Debug(time.Now().Unix(), "nextPlayer", "pos:", pos, "next_pos:", next_pos)
					return room_player
				}
			}
		}
	}

	return room.players[0]
}

//取指定玩家的对家
func (room *Room) oppositePlayer(player *Player) *Player {
	if nil == player{
		return nil
	}

	pos := player.GetPosition()
	opp_pos := (pos + 2) % int32(room.config.NeedPlayerNum)
	for _, room_player := range room.players {
		if room_player.GetPosition() == opp_pos && !room_player.GetIsEndPlaying(){
			return room_player
		}
	}
	return nil
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

	case OperateConfirmPlayAlone:
		if pa_data, ok := op.Data.(*OperateConfirmPlayAloneData); ok {
			//log.Debug(log_time, room, "Room.dealPlayerOperate player play_alone :", op.Operator)
			op.Operator.PlayAlone(pa_data.IsPlayAlone)
			op.ResultCh <- true
			room.broadcastPlayerSuccessOperated(op)
			return true
		}

	case OperateDrop:
		if drop_data, ok := op.Data.(*OperateDropData); ok {
			if op.Operator.Drop(drop_data.whatGroup) {
				//出牌，计算分数和奖
				dropped_score := room.getCardsScores(drop_data.whatGroup)
				room.AddTableScore(dropped_score)
				prize := room.getPrizeByCardsType(drop_data.cardsType)
				op.Operator.AddPrize(prize)
				if drop_data.cardsType != card.CardsType_NO{
					room.SetCardsType(drop_data.cardsType)
					room.SetPlaneNum(drop_data.planeNum)
					room.SetWeight(drop_data.weight)
				}

				log.Debug(time.Now().Unix(), room, "Room.dealPlayerOperate player drop :", op.Operator)
				op.ResultCh <- true
				room.broadcastPlayerSuccessOperated(op)
				return true
			}
		}

	case OperatePass:
		if _, ok := op.Data.(*OperatePassData); ok {
			log.Debug(time.Now().Unix(), room, "Room.dealPlayerOperate player pass :", op.Operator)
			op.ResultCh <- true
			room.broadcastPlayerSuccessOperated(op)
			return true
		}

	}
	op.ResultCh <- false
	return false
}

func (room *Room) getCardsScores(drop_cards []*card.Card) (int32)  {
	score := int32(0)
	for _, drop_card := range drop_cards {
		if drop_card.CardNo == 5 {
			score += 5
		}else if drop_card.CardNo == 10 {
			score += 10
		}else if drop_card.CardNo == 13 {
			score += 10
		}
	}

	return score
}

func (room *Room) getPrizeByCardsType(cards_type int) (int32) {
	prize := int32(0)
	switch cards_type {
	case 24:
		prize = 1
	case 25:
		prize = 2
	case 26:
		prize = 3
	case 27, 28:
		prize = 5
	}
	return prize
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
	//log.Debug(time.Now().Unix(), room, "Room.putCardsToPlayers, init_type:", init_type)
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
				//log.Debug("put_card:", put_card)
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

func (room *Room) Reset() {
	room.opMaster = nil
	room.waitOperator = nil
	room.masterPlayer = nil
	room.nextOpMaster = nil
	room.assistPlayer = nil
	room.isPlayAlone = false
	room.turnCard = nil
	room.tableScore = 0
	room.endPlayingNum = 0
	room.cardsType = 0
	room.planeNum = 0
	room.weight = 0
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
		case op := <-room.playAloneCh[idx]:
			op.ResultCh <- false
		default:
		}

		select {
		case op := <-room.dropCardCh[idx]:
			op.ResultCh <- false
		default:
		}

		select {
		case op := <-room.passCh[idx]:
			op.ResultCh <- false
		default:
		}

		select {
		case op := <-room.roomReadyCh[idx]:
			op.ResultCh <- false
		default:
		}
	}
}
