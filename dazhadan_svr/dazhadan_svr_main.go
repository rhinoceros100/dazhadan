package main

import (
	"bufio"
	"os"
	"strings"
	"strconv"
	"time"
	"dazhadan/dazhadan_svr/playing"
	"dazhadan/dazhadan_svr/util"
	"dazhadan/dazhadan_svr/log"
)

func help() {
	log.Debug("-----------------help---------------------")
	log.Debug("h")
	log.Debug("exit")
	log.Debug("mycards")
	log.Debug(playing.OperateEnterRoom, int(playing.OperateEnterRoom))
	log.Debug(playing.OperateReadyRoom, int(playing.OperateReadyRoom))
	log.Debug(playing.OperateLeaveRoom, int(playing.OperateLeaveRoom))
	log.Debug("-----------------help---------------------")
}

type PlayerObserver struct {}
func (ob *PlayerObserver) OnMsg(player *playing.Player, msg *playing.Message) {
	log_time := time.Now().Unix()
	switch msg.Type {
	case playing.MsgEnterRoom:
		if enter_data, ok := msg.Data.(*playing.EnterRoomMsgData); ok {
			log.Debug(log_time, player, "OnMsg MsgEnterRoom, EnterPlayer:", enter_data.EnterPlayer)
		}
	case playing.MsgReadyRoom:
		if enter_data, ok := msg.Data.(*playing.ReadyRoomMsgData); ok {
			log.Debug(log_time, player, "OnMsg MsgReadyRoom, ReadyPlayer:", enter_data.ReadyPlayer)
		}
	case playing.MsgLeaveRoom:
		if enter_data, ok := msg.Data.(*playing.LeaveRoomMsgData); ok {
			log.Debug(log_time, player, "OnMsg MsgLeaveRoom, LeavePlayer:", enter_data.LeavePlayer)
		}
	case playing.MsgGameEnd:
		if _, ok := msg.Data.(*playing.GameEndMsgData); ok {
			log.Debug(log_time, player, "OnMsg MsgGameEnd")
		}
	case playing.MsgRoomClosed:
		if _, ok := msg.Data.(*playing.RoomClosedMsgData); ok {
			log.Debug(log_time, player, "OnMsg MsgRoomClosed")
		}
	case playing.MsgGetInitCards:
		if init_data, ok := msg.Data.(*playing.GetInitCardsMsgData); ok {
			log.Debug(log_time, player, "OnMsg MsgGetInitCards, PlayingCards:", init_data.PlayingCards)
		}
	case playing.MsgWaitDadu:
		if dadu_data, ok := msg.Data.(*playing.WaitDaduMsgData); ok {
			log.Debug(log_time, player, "OnMsg MsgWaitDadu left_sec:", dadu_data.LeftSec, "wait_player", dadu_data.WaitDaduPlayer)
		}
	case playing.MsgConfirmDadu:
		if dadu_data, ok := msg.Data.(*playing.ConfirmDaduMsgData); ok {
			log.Debug(log_time, player, "OnMsg MsgConfirmDadu, IsDadu:", dadu_data.IsDadu, "dadu_player", dadu_data.DaduPlayer)
		}
	case playing.MsgSwitchPosition:
		if sw_data, ok := msg.Data.(*playing.SwitchPositionMsgData); ok {
			log.Debug(log_time, player, "OnMsg MsgSwitchPosition, uid1:", sw_data.OppUid, "pos1:", sw_data.OppPos, "uid2:", sw_data.AssistUid, "pos2:", sw_data.AssistPos)
		}
	case playing.MsgStartPlay:
		if sp_data, ok := msg.Data.(*playing.StartPlayMsgData); ok {
			log.Debug(log_time, player, "OnMsg MsgStartPlay, is_dadu:", sp_data.IsDadu, "assist:", sp_data.Assist, "master:", sp_data.Master)
		}
	case playing.MsgSwitchOperator:
		if _, ok := msg.Data.(*playing.SwitchOperatorMsgData); ok {
			log.Debug(log_time, player, "OnMsg MsgSwitchOperator", msg.Owner)
		}
	case playing.MsgDrop:
		if _, ok := msg.Data.(*playing.DropMsgData); ok {
			log.Debug(log_time, player, "OnMsg MsgDrop", msg.Owner)
		}
	case playing.MsgGuo:
		if _, ok := msg.Data.(*playing.GuoMsgData); ok {
			log.Debug(log_time, player, "OnMsg MsgGuo", msg.Owner)
		}
	case playing.MsgJiesuan:
		if jiesuan_data, ok := msg.Data.(*playing.JiesuanMsgData); ok {
			log.Debug(log_time, player, "OnMsg MsgJiesuan, jiesuan_data:")
			for _, score_data := range jiesuan_data.Scores	{
				log.Debug(score_data.P, score_data.P.IsMaster(), "Score:", score_data.Score, score_data.Paixing)
			}
		}
	}
}

func main() {
	running := true

	//init room
	conf := playing.NewRoomConfig()
	conf.Init(1, 2, 1)
	room := playing.NewRoom(util.UniqueId(), conf)
	room.Start()

	robots := []*playing.Player{
		playing.NewPlayer(1),
		playing.NewPlayer(2),
		playing.NewPlayer(3),
	}

	for _, robot := range robots {
		robot.OperateEnterRoom(room)
		robot.AddObserver(&PlayerObserver{})
	}

	curPlayer := playing.NewPlayer(4)
	curPlayer.AddObserver(&PlayerObserver{})

	go func() {
		time.Sleep(time.Second * 1)
		robots[0].OperateDoReady()
		time.Sleep(time.Second * 2)
		robots[1].OperateDoReady()
		time.Sleep(time.Second * 5)
		robots[2].OperateDoReady()
		//curPlayer.OperateDoReady()
	}()

	reader := bufio.NewReader(os.Stdin)
	for running {
		data, _, _ := reader.ReadLine()
		cmd := string(data)
		if cmd == "h" {
			help()
		} else if cmd == "exit" {
			return
		} else if cmd == "mycards" {
			log.Debug(curPlayer.GetPlayingCards())
		}
		splits := strings.Split(cmd, " ")
		c, _ := strconv.Atoi(splits[0])
		switch playing.OperateType(c) {
		case playing.OperateEnterRoom:
			curPlayer.OperateEnterRoom(room)
		case playing.OperateReadyRoom:
			curPlayer.OperateDoReady()
		case playing.OperateLeaveRoom:
			curPlayer.OperateLeaveRoom()
		}
	}
}
