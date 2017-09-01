package playing

import (
	"time"
)

type RoomConfig struct {
	NeedPlayerNum			int32        `json:"need_player_num"`
	WinCoin				int32        `json:"win_coin"`
	ShuangjiCoin			int32        `json:"shuangji_coin"`
	DaduCoin			int32        `json:"dadu_coin"`
	PrizeCoin			int32        `json:"prize_coin"`
	InitType			int32        `json:"init_type"`
	MaxPlayGameCnt			int        `json:"max_play_game_cnt"`	//最大的游戏局数

	WaitPlayerEnterRoomTimeout	int        `json:"wait_player_enter_room_timeout"`
	WaitPlayerOperateTimeout	int        `json:"wait_player_operate_timeout"`
	WaitDaduSec                	time.Duration      `json:"wait_dadu_sec"`	//等待打独时长
	WaitDropSec                	time.Duration      `json:"wait_drop_sec"`	//等待出牌时长
	AfterSwitchPositionSleep        time.Duration      `json:"after_switch_position_sleep"`	//交换位置后sleep时长

	/************************************************/
	WaitScrambleSec                	time.Duration      `json:"wait_scramble_sec"`	//等待抢庄时长
	WaitBetSec                 	time.Duration      `json:"wait_bet_sec"`	//等待下注时长
	WaitShowCardsSec              	time.Duration      `json:"wait_show_cards_sec"`	//等待亮牌时长
	WaitReadySec              	time.Duration      `json:"wait_ready_sec"`	//等待准备时长
	AfterBetSleep                   time.Duration      `json:"after_bet_sleep"`	//下注后sleep时长
	AfterShowCardsSleep             time.Duration      `json:"after_show_cards_sleep"`	//亮牌后sleep时长
}

func NewRoomConfig() *RoomConfig {
	return &RoomConfig{}
}

func (config *RoomConfig) Init(score_type, prize_type, init_type int32) {
	if score_type == 1 {
		config.WinCoin = 10
		config.ShuangjiCoin = 15
		config.DaduCoin = 20
	}else {
		config.WinCoin = 10
		config.ShuangjiCoin = 20
		config.DaduCoin = 30
	}

	if score_type == 1 {
		config.PrizeCoin = 3
	}else {
		config.PrizeCoin = 5
	}

	config.InitType = init_type
	config.NeedPlayerNum = 4
	config.MaxPlayGameCnt = 3
	config.WaitPlayerEnterRoomTimeout = 300
	config.WaitPlayerOperateTimeout = 300
	config.WaitDaduSec = 10
	config.AfterSwitchPositionSleep = 1
	config.WaitDropSec = 3

	/************************************************/
	config.WaitScrambleSec = 10
	config.WaitBetSec = 15
	config.WaitShowCardsSec = 15
	config.WaitReadySec = 15

	config.AfterBetSleep = 4
	config.AfterShowCardsSleep = 5
}