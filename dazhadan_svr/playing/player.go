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

	isPlayAlone		bool
	isEndPlaying		bool
	needDrop		bool
	isWin			bool
	totalCoin	        int32                   //总金币
	prizeCoin	        int32                   //奖金币
	coin		        int32                   //本场金币
	rank		        int32                   //排名
	score		        int32                   //一轮得分
	prize		        int32                   //获得奖励次数

	winNum		        int32                   //总共赢的次数
	shuangjiNum		int32                   //总共双基的次数
	paSuccNum		int32                   //打独成功的次数
	totalPrize		int32                   //获得的总奖金数

	playingCards 	*card.PlayingCards	//玩家手上的牌
	observers	 []PlayerObserver
}

func NewPlayer(id uint64) *Player {
	player :=  &Player{
		id:		id,
		position:       10,
		isReady:        false,
		isPlayAlone:    false,
		isEndPlaying:   false,
		needDrop:     	false,
		isWin:     	false,

		rank:      	0,
		score:     	0,
		prize:     	0,
		totalCoin: 	0,
		prizeCoin: 	0,
		coin:      	0,
		winNum:      	0,
		shuangjiNum:    0,
		paSuccNum:      0,
		totalPrize:     0,

		playingCards:	card.NewPlayingCards(),
		observers:	make([]PlayerObserver, 0),
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

func (player *Player) GetCoin() int32 {
	return player.coin
}

func (player *Player) AddCoin(add int32) {
	player.coin += add
}

func (player *Player) ResetCoin() {
	player.coin = 0
}

func (player *Player) GetPrizeCoin() int32 {
	return player.prizeCoin
}

func (player *Player) AddPrizeCoin(add int32) {
	player.prizeCoin += add
}

func (player *Player) ResetPrizeCoin() {
	player.prizeCoin = 0
}

func (player *Player) GetIsWin() bool {
	return player.isWin
}

func (player *Player) SetIsWin(is_win bool) {
	player.isWin = is_win
}

func (player *Player) GetScore() int32 {
	return player.score
}

func (player *Player) AddScore(add int32) {
	player.score += add
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
	return player.prize
}

func (player *Player) AddPrize(add int32) {
	player.prize += add
}

func (player *Player) ResetPrize() {
	player.prize = 0
}

func (player *Player) GetWinNum() int32 {
	return player.winNum
}

func (player *Player) IncWinNum() int32 {
	player.winNum++
	return player.winNum
}

func (player *Player) GetShuangjiNum() int32 {
	return player.shuangjiNum
}

func (player *Player) IncShuangjiNum() int32 {
	player.shuangjiNum++
	return player.shuangjiNum
}

func (player *Player) GetPaSuccNum() int32 {
	return player.paSuccNum
}

func (player *Player) IncPaSuccNum() int32 {
	player.paSuccNum++
	return player.paSuccNum
}

func (player *Player) GetTotalPrize() int32 {
	return player.totalPrize
}

func (player *Player) AddTotalPrize(add int32) int32 {
	player.totalPrize += add
	return player.totalPrize
}

func (player *Player) GetIsPlayAlone() bool {
	return player.isPlayAlone
}

func (player *Player) SetIsPlayAlone(is_play_alone bool) {
	player.isPlayAlone = is_play_alone
}

func (player *Player) GetIsEndPlaying() bool {
	return player.isEndPlaying
}

func (player *Player) SetIsEndPlaying(is_end_playing bool) {
	player.isEndPlaying = is_end_playing
}

func (player *Player) GetNeedDrop() bool {
	return player.needDrop
}

func (player *Player) SetNeedDrop(need_drop bool) {
	player.needDrop = need_drop
}


func (player *Player) Reset() {
	//log.Debug(time.Now().Unix(), player,"Player.Reset")
	player.playingCards.Reset()
	player.SetIsReady(false)
	player.SetIsPlayAlone(false)
	player.SetIsEndPlaying(false)
	player.SetNeedDrop(false)
	player.SetRank(0)
	player.SetIsWin(false)
	player.ResetPrize()
	player.ResetScore()

	coin := player.GetCoin()
	player.AddTotalCoin(coin)
	player.ResetCoin()
	player.ResetPrizeCoin()
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

func (player *Player) OperateConfirmPlayAlone(is_play_alone bool) bool {
	log.Debug(player, "OperateConfirmPlayAlone:", is_play_alone)
	data := &OperateConfirmPlayAloneData{is_play_alone}
	op := NewOperateConfirmPlayAlone(player, data)
	player.room.PlayerOperate(op)
	return player.waitResult(op.ResultCh)
}

func (player *Player) OperateDropCard(cards []*card.Card) bool {
	log.Debug(player, "OperateDrop cards :", cards)
	data := &OperateDropData{
		whatGroup: cards,
	}

	cards_num := player.playingCards.CardsInHand.Len()
	is_last_cards := false
	if cards_num == len(cards) {
		is_last_cards = true
	}
	drop_cards := card.CreateNewCards(cards)
	check_cards_type, check_plane_num := 0, 0
	if player.room.GetCardsType() == card.CardsType_PLANE32 {
		check_cards_type = card.CardsType_PLANE32
		check_plane_num = player.room.GetPlaneNum()
	}
	data.cardsType, data.planeNum, data.weight = card.GetCardsType(drop_cards, is_last_cards, check_cards_type, check_plane_num)
	can_cover := player.room.canCover(data.cardsType, data.planeNum, data.weight)
	log.Debug("******can_cover:", can_cover)
	op := NewOperateDrop(player, data)
	player.room.PlayerOperate(op)
	return player.waitResult(op.ResultCh)
}

func (player *Player) OperatePass() bool {
	log.Debug(player, "OperatePass")
	data := &OperatePassData{}
	op := NewOperatePass(player, data)
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

func (player *Player) PlayAlone(is_play_alone bool) {
	//log.Debug(time.Now().Unix(), player, "PlayAlone", player.room)
	player.SetIsPlayAlone(is_play_alone)
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
	case OperateConfirmPlayAlone:
		player.OnPlayerPlayAlone(op)
	case OperateSwitchOperator:
		player.onSwitchOperator(op)
	case OperateDrop:
		player.OnDrop(op)
	case OperatePass:
		player.OnPass(op)
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

func (player *Player) OnPlayerPlayAlone(op *Operate) {
	//log.Debug(time.Now().Unix(), player, "OnPlayerPlayAlone")
	if pa_data, ok := op.Data.(*OperateConfirmPlayAloneData); ok {
		data := &ConfirmPlayAloneMsgData{
			IsPlayAlone:pa_data.IsPlayAlone,
			PlayAlonePlayer:op.Operator,
		}
		player.notifyObserver(NewConfirmPlayAloneMsg(player, data))
	}
}

func (player *Player) onSwitchOperator(op *Operate) {
	if so_data, ok := op.Data.(*OperateSwitchOperatorData); ok {
		/*if op.Operator == player {
			return
		}*/
		msgData := &SwitchOperatorMsgData{
			NeedDropCard:so_data.MustDrop,
			CanDrop:so_data.CanDrop,
			SwitchedPlayer:op.Operator,
		}
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
			TableScore:player.room.GetTableScore(),
			CardsType:drop_data.cardsType,
			PlaneNum:drop_data.planeNum,
			Weight:drop_data.weight,
		}
		player.notifyObserver(NewDropMsg(op.Operator, msgData))
	}
}

func (player *Player) OnPass(op *Operate) {
	if _, ok := op.Data.(*OperatePassData); ok {
		/*if op.Operator == player {
			return
		}*/
		msgData := &PassMsgData{}
		player.notifyObserver(NewPassMsg(op.Operator, msgData))
	}
}

func (player *Player) OnWaitPlayAlone(msg *Message) {
	player.notifyObserver(msg)
}

func (player *Player) OnSwitchPosition(msg *Message) {
	player.notifyObserver(msg)
}

func (player *Player) OnStartPlay(msg *Message) {
	player.notifyObserver(msg)
}

func (player *Player) OnSummary(msg *Message) {
	//log.Debug(time.Now().Unix(), player, "OnSummary")

	player.notifyObserver(msg)
}

func (player *Player) OnGetInitCards() {
	//log.Debug(time.Now().Unix(), player, "OnGetInitCards", player.playingCards)

	data := &GetInitCardsMsgData{
		PlayingCards: player.playingCards,
	}
	player.notifyObserver(NewGetInitCardsMsg(player, data))
}

func (player *Player) OnRoomClosed(msg *Message) {
	//log.Debug(time.Now().Unix(), player, "OnRoomClosed")
	player.room = nil
	//player.Reset()

	player.notifyObserver(msg)
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

func (player *Player) GetBomb8Num() (bomb4_num, bomb5_num, bomb6_num, bomb7_num, bomb8_num, bomb_joker int) {
	//log.Debug("GetBomb8Num")
	nums := [15]int{}
	for _, hand_card := range player.playingCards.CardsInHand.GetData() {
		nums[hand_card.CardNo]++
	}

	bomb4_num, bomb5_num, bomb6_num, bomb7_num, bomb8_num, bomb_joker = 0, 0, 0, 0, 0, 0
	for i, num := range nums {
		//log.Debug("i:", i, "num:", num)
		if num >= 4 {
			bomb4_num ++
		}
		if num >= 5 {
			bomb5_num ++
		}
		if num >= 6 {
			bomb6_num ++
		}
		if num >= 7 {
			bomb7_num ++
		}
		if num >= 8 {
			bomb8_num ++
		}
		if i == 14 && num == 4 {
			bomb_joker ++
		}
	}
	return
}

func (player *Player) GetCanDrop() bool{
	canDrop := true
	cardsType := player.room.GetCardsType()
	//planeNum := player.room.GetPlaneNum()
	drop_weight := player.room.GetWeight()

	//计算玩家此时的牌
	var normal510k, true510k = player.GetHave510K(player.GetPlayingCards().CardsInHand.GetData())
	var bomb4, bomb5, bomb6, bomb7, bomb8 = 0, 0, 0, 0, 0
	var c1, c2, c3, straight, pairs, triples = 0, 0, 0, 0, 0, 0
	var straight_start, straight_end, pairs_start, pairs_end, triples_start, triples_end = 0, 0, 0, 0, 0, 0
	var bombType = 0
	if normal510k{
		bombType = card.CardsType_510K
	}
	if true510k{
		bombType = card.CardsType_TRUE510K
	}

	//计算每个点数的牌的数量
	arr := [18]int{}
	for _, hand_card := range player.GetPlayingCards().CardsInHand.GetData() {
		arr[hand_card.Weight] ++
	}
	//计算玩家可能出现的牌型
	for weight, num := range arr{
		if num > 0 {
			c1 = weight
			if weight < 16 {
				if straight_start == 0 {
					straight_start = weight
					straight_end = weight
				}else{
					if straight_end != weight - 1 {
						straight_start = 0
						straight_end = 0
					}else{
						straight_end = weight
					}
				}
				if straight_end - straight_start >= 4{
					straight = straight_start
				}
			}
		}
		if num > 1 {
			c2 = weight
			if weight < 16 {
				if pairs_start == 0 {
					pairs_start = weight
					pairs_end = weight
				}else{
					if pairs_end != weight - 1 {
						pairs_start = 0
						pairs_end = 0
					}else{
						pairs_end = weight
					}
				}
				if pairs_end - pairs_start >= 2{
					pairs = pairs_start
				}
			}
		}
		if num > 2 {
			c3 = weight
			if weight < 16 {
				if triples_start == 0 {
					triples_start = weight
					triples_end = weight
				}else{
					if triples_end != weight - 1 {
						triples_start = 0
						triples_end = 0
					}else{
						triples_end = weight
					}
				}
				if triples_end - triples_start >= 1{
					triples = triples_start
				}
			}
		}
		if num > 3 {
			bomb4 = weight
			if bombType < card.CardsType_BOMB4 {
				bombType = card.CardsType_BOMB4
			}
		}
		if num > 4 {
			bomb5 = weight
			if bombType < card.CardsType_BOMB5 {
				bombType = card.CardsType_BOMB5
			}
		}
		if num > 5 {
			bomb6 = weight
			if bombType < card.CardsType_BOMB6 {
				bombType = card.CardsType_BOMB6
			}
		}
		if num > 6 {
			bomb7 = weight
			if bombType < card.CardsType_BOMB7 {
				bombType = card.CardsType_BOMB7
			}
		}
		if num > 7 {
			bomb8 = weight
			if bombType < card.CardsType_BOMB8 {
				bombType = card.CardsType_BOMB8
			}
		}
	}
	if arr[16] == 2 && arr[17] == 2 {
		bombType = card.CardsType_JOKER
		//bombjoker = 16
	}

	//判断玩家是否能够要的起牌
	if cardsType > 20 {//如果出的牌是炸弹，查看手上的牌是否可以拼出更大的炸弹
		if bombType < cardsType {
			return false
		}
		if bombType > cardsType {
			return true
		}
		if bombType == card.CardsType_510K || bombType == card.CardsType_TRUE510K {
			return false
		}
		if bombType == card.CardsType_BOMB4 {
			return bomb4 > drop_weight
		}
		if bombType == card.CardsType_BOMB5 {
			return bomb5 > drop_weight
		}
		if bombType == card.CardsType_BOMB6 {
			return bomb6 > drop_weight
		}
		if bombType == card.CardsType_BOMB7 {
			return bomb7 > drop_weight
		}
		if bombType == card.CardsType_BOMB8 {
			return bomb8 > drop_weight
		}
	}else {//出的牌是普通牌型
		if bombType > 20 {
			return true
		}
		if cardsType == card.CardsType_SINGLE {
			return c1 > drop_weight
		}
		if cardsType == card.CardsType_PAIR {
			return c2 > drop_weight
		}
		if cardsType == card.CardsType_32 {
			return c3 > drop_weight
		}
		//顺子，连对，飞机只考虑了权重，未考虑张数，如果玩家出的连牌较多可能实际管不上却返回true，情况较少，可忽略
		if cardsType == card.CardsType_STAIGHT {
			return straight > drop_weight
		}
		if cardsType == card.CardsType_PAIRS {
			return pairs > drop_weight
		}
		if cardsType == card.CardsType_PLANE32 {
			return triples > drop_weight
		}
	}

	return canDrop
}

func (player *Player) GetHave510K(cards []*card.Card) (have510k, haveTrue510k bool){
	var d5, c5, h5, s5, d10, c10, h10, s10, d13, c13, h13, s13 = false, false, false, false, false, false, false, false, false, false, false, false
	var n5, n10, n13 = 0, 0, 0
	for _, hand_card := range cards {
		if hand_card.CardNo == 5 {
			n5++
			switch hand_card.CardType {
			case card.CardType_Diamond:
				d5 = true
			case card.CardType_Club:
				c5 = true
			case card.CardType_Heart:
				h5 = true
			case card.CardType_Spade:
				s5 = true
			}
		}
		if hand_card.CardNo == 10 {
			n10++
			switch hand_card.CardType {
			case card.CardType_Diamond:
				d10 = true
			case card.CardType_Club:
				c10 = true
			case card.CardType_Heart:
				h10 = true
			case card.CardType_Spade:
				s10 = true
			}
		}
		if hand_card.CardNo == 13 {
			n13++
			switch hand_card.CardType {
			case card.CardType_Diamond:
				d13 = true
			case card.CardType_Club:
				c13 = true
			case card.CardType_Heart:
				h13 = true
			case card.CardType_Spade:
				s13 = true
			}
		}
	}
	if n5 > 0 && n10 > 0 && n13 > 0 {
		have510k = true
	}
	if (d5 && d10 && d13) || (c5 && c10 && c13) || (h5 && h10 && h13) || (s5 && s10 && s13){
		haveTrue510k = true
	}
	return
}
