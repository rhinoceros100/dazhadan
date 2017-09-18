package card

func GetCardsType(the_cards *Cards, is_last_cards bool, check_cards_type, check_plane_num int) (cards_type int, plane_num int, weight int) {
	the_cards.Sort()
	drop_cards := the_cards.GetData()
	cards_type = CardsType_NO
	plane_num = 0
	weight = 0
	cards_len := len(drop_cards)

	if cards_len == 0 {
		return
	}
	if cards_len == 1 {
		cards_type = CardsType_SINGLE
		weight = drop_cards[0].Weight
		return
	}

	most, sames := GetSameCardsNum(drop_cards)
	if cards_len == 2 {
		if most == 2 {
			//大王、小王不能成对
			if drop_cards[0].Weight != drop_cards[1].Weight {
				return CardsType_NO, plane_num, weight
			}
			return CardsType_PAIR, plane_num, drop_cards[0].Weight
		}
		return
	}
	//单牌和双牌需要区分大小王的权重，其他不需要
	weight = GetStraightWeight(sames)
	if cards_len == 3 {
		if most == 3 {
			if is_last_cards && drop_cards[0].CardNo != 14{
				return CardsType_32, plane_num, weight
			}else {
				return CardsType_NO, plane_num, weight
			}
		}
		_, cards_type := Is510K(drop_cards)
		return cards_type, plane_num, weight
	}
	if cards_len == 4 {
		//判断是否为王炸
		is_bomb_joker := true
		for _, drop_card := range drop_cards {
			if drop_card.CardNo != 14 {
				is_bomb_joker = false
			}
		}
		if is_bomb_joker {
			cards_type = CardsType_JOKER
			return cards_type, plane_num, weight
		}
		//王炸或者四炸
		if most == 4 {
			if weight == 15 {
				cards_type = CardsType_BOMB5
				weight = 2
			}else {
				cards_type = CardsType_BOMB4
			}
			return cards_type, plane_num, weight
		}
		//判断是否为三带二
		if most == 3 {
			if is_last_cards {
				return CardsType_32, plane_num, weight
			}
		}
		return CardsType_NO, plane_num, weight
	}

	//如果要判定飞机，优先判定
	if check_cards_type == CardsType_PLANE32 {
		return Check3Plane(drop_cards, is_last_cards, check_cards_type, check_plane_num)
	}

	if cards_len>= 5 && cards_len <= 8 {
		//炸弹
		if most == cards_len {
			switch most {
			case 5:
				return CardsType_BOMB5, plane_num, weight
			case 6:
				return CardsType_BOMB6, plane_num, weight
			case 7:
				return CardsType_BOMB7, plane_num, weight
			case 8:
				return CardsType_BOMB8, plane_num, weight
			}
		}

		//判断m带m-1
		if most >= 4 && most <= 7 {
			if is_last_cards{
				switch most {
				case 7:
					cards_type = CardsType_76
				case 6:
					cards_type = CardsType_65
				case 5:
					cards_type = CardsType_54
				case 4:
					cards_type = CardsType_43
				}
			}
			if most == 4 && cards_len == 7 {
				cards_type = CardsType_43
			}
			if most == 4 && len(sames) == 2 && is_last_cards {
				if IsStraight(sames) {
					plane_num = 2
					cards_type = CardsType_PLANE43
				}
			}
			return
		}
		if most == 3 {
			if cards_len == 5 {
				cards_type = CardsType_32
				return
			}
			if len(sames) == 2 {
				if is_last_cards {
					if IsStraight(sames) {
						plane_num = 2
						cards_type = CardsType_PLANE32
					}
				}
			}
			if cards_type == CardsType_NO {
				return Check3Plane(drop_cards, is_last_cards, check_cards_type, check_plane_num)
			}else {
				return
			}
		}
		if most == 2 {
			same_len := len(sames)
			/*println("same_len:", same_len)
			println("cards_len:", cards_len)
			for _, same := range sames {
				println(same)
			}*/
			if cards_len == same_len * 2 {
				if IsStraight(sames) {
					plane_num = same_len
					cards_type = CardsType_PAIRS
				}
			}
			return
		}
		if IsStraight(sames) {
			plane_num = cards_len
			cards_type = CardsType_STAIGHT
		}
		return
	}

	if cards_len > 8 {
		same_len := len(sames)
		if most >= 3 && most <= 8 {
			take := 2 * most - 1
			if (cards_len == take) || (cards_len < take && is_last_cards) {
				switch most {
				case 8:
					cards_type = CardsType_87
				case 7:
					cards_type = CardsType_76
				case 6:
					cards_type = CardsType_65
				case 5:
					cards_type = CardsType_54
				case 4:
					cards_type = CardsType_43
				case 3:
					cards_type = CardsType_32
				}
			}
			if (cards_len == take * same_len) || (cards_len < take * same_len && is_last_cards) {
				if IsStraight(sames) {
					plane_num = same_len
					switch most {
					case 8:
						cards_type = CardsType_PLANE87
					case 7:
						cards_type = CardsType_PLANE76
					case 6:
						cards_type = CardsType_PLANE65
					case 5:
						cards_type = CardsType_PLANE54
					case 4:
						cards_type = CardsType_PLANE43
					case 3:
						cards_type = CardsType_PLANE32
					}
				}
			}
			//类似于3334447779的情况
			if (cards_type == CardsType_NO && cards_len < take * same_len) {
				for i := 3; i <= 9; i++ {//i个most，其中i-1个可以组成顺子
					if same_len == i && cards_len <= take * (i-1) {
						if (cards_len == take * (i-1)) || (cards_len < take * (i-1) && is_last_cards) {
							if IsStraight(sames[0:(i-1)]) || IsStraight(sames[1:]){
								plane_num = i - 1
								switch most {
								case 8:
									cards_type = CardsType_PLANE87
								case 7:
									cards_type = CardsType_PLANE76
								case 6:
									cards_type = CardsType_PLANE65
								case 5:
									cards_type = CardsType_PLANE54
								case 4:
									cards_type = CardsType_PLANE43
								case 3:
									cards_type = CardsType_PLANE32
								}
							}
							if IsStraight(sames[0:(i-1)]) {
								weight = sames[0]
							}else if IsStraight(sames[1:]) {
								weight = sames[1]
							}
						}
					}
				}
				if cards_type != CardsType_NO {
					return
				}
				for i := 5; i <= 9; i++ {//i个most，其中i-2个可以组成顺子
					if same_len == i && cards_len <= take * (i-2) {
						if (cards_len == take * (i-2)) || (cards_len < take * (i-2) && is_last_cards) {
							if IsStraight(sames[0:(i-2)]) || IsStraight(sames[1:(i-1)]) || IsStraight(sames[2:]){
								plane_num = i - 2
								switch most {
								case 5:
									cards_type = CardsType_PLANE54
								case 4:
									cards_type = CardsType_PLANE43
								case 3:
									cards_type = CardsType_PLANE32
								}
							}
							if IsStraight(sames[0:(i-2)]) {
								weight = sames[0]
							}else if IsStraight(sames[1:(i-1)]) {
								weight = sames[1]
							}else if IsStraight(sames[2:]) {
								weight = sames[2]
							}
						}
					}
				}
				if cards_type != CardsType_NO {
					return
				}
				for i := 8; i <= 9; i++ {//i个most，其中i-3个可以组成顺子
					if same_len == i && cards_len <= take * (i-3) {
						if (cards_len == take * (i-3)) || (cards_len < take * (i-3) && is_last_cards) {
							if IsStraight(sames[0:(i-3)]) || IsStraight(sames[1:(i-2)]) || IsStraight(sames[2:(i-1)]) || IsStraight(sames[3:]){
								plane_num = i - 3
								switch most {
								case 3:
									cards_type = CardsType_PLANE32
								}
							}
							if IsStraight(sames[0:(i-3)]) {
								weight = sames[0]
							}else if IsStraight(sames[1:(i-2)]) {
								weight = sames[1]
							}else if IsStraight(sames[2:(i-1)]) {
								weight = sames[2]
							}else if IsStraight(sames[3:]) {
								weight = sames[3]
							}
						}
					}
				}
			}
			if cards_type == CardsType_NO {
				return Check3Plane(drop_cards, is_last_cards, check_cards_type, check_plane_num)
			}else {
				return
			}
		}

		if most == 2 {
			if cards_len == same_len * 2 {
				if IsStraight(sames) {
					plane_num = same_len
					cards_type = CardsType_PAIRS
				}
			}
			return
		}
		if IsStraight(sames) {
			plane_num = cards_len
			cards_type = CardsType_STAIGHT
		}
		return
	}

	return CardsType_NO, plane_num, weight
}

//检查牌型是否为三带二飞机的牌型  如3334444556
func Check3Plane(drop_cards []*Card, is_last_cards bool, check_cards_type, check_plane_num int) (cards_type int, plane_num int, weight int) {
	cards_type = CardsType_NO
	plane_num = 0
	weight = 0
	take := 2 * 3 -1
	cards_len := len(drop_cards)
	sames := GetMoreThanNCardsNum(drop_cards, 3)
	same_len := len(sames)

	//check_plane_num>0时为指定牌型，=0时为查找是否符合飞机牌型
	if check_plane_num > 0 {
		if same_len < check_plane_num {
			return
		}
		if cards_len == take * check_plane_num || (is_last_cards && (cards_len < take * check_plane_num)) {
			diff := same_len - check_plane_num
			for i := 0; i <= diff; i++ {
				if IsStraight(sames[diff-i:same_len-i]) {
					plane_num = check_plane_num
					weight = sames[diff-i]
					cards_type = CardsType_PLANE32
					return
				}
			}
		}
		return
	}else {
		if same_len < 2 {
			return
		}

		for j := same_len; j >= 2; j-- {
			diff := same_len - j
			if cards_len == take * j || (is_last_cards && (cards_len < take * j)) {
				for i := 0; i <= diff; i++ {
					if IsStraight(sames[diff-i:same_len-i]) {
						plane_num = j
						weight = sames[diff-i]
						cards_type = CardsType_PLANE32
						return
					}
				}
			}
		}
	}
	return
}

//获取一组牌中数量最多的数字相同的牌的数量
func GetSameCardsNum(drop_cards []*Card) (most int, same_card_nums []int) {
	arr := [18]int{}
	for _, drop_card := range drop_cards {
		arr[drop_card.Weight] ++
	}

	same_card_nums = make([]int, 0)
	most = 0
	for i, num := range arr {
		if num > most {
			same_card_nums = same_card_nums[0:0]
			most = num
			same_card_nums = append(same_card_nums, i)
		}else if num == most {
			same_card_nums = append(same_card_nums, i)
		}
	}
	return
}

//获取一组牌中多于N张的牌
func GetMoreThanNCardsNum(drop_cards []*Card, min_num int) (card_nums []int) {
	arr := [18]int{}
	for _, drop_card := range drop_cards {
		arr[drop_card.Weight] ++
	}

	card_nums = make([]int, 0)
	for i, num := range arr {
		if num >= min_num {
			card_nums = append(card_nums, i)
		}
	}
	return
}

func Is510K(drop_cards []*Card) (is_510k bool, cards_type int) {
	cards_len := len(drop_cards)
	is_510k = false
	cards_type = CardsType_NO

	if cards_len != 3 {
		return
	}

	if drop_cards[0].CardNo == 5 && drop_cards[1].CardNo == 10 && drop_cards[2].CardNo == 13 {
		is_510k = true
		cards_type = CardsType_510K
		if IsSameCardType(drop_cards) {
			cards_type = CardsType_TRUE510K
		}
	}
	return
}

func IsSameCardType(drop_cards []*Card) (bool) {
	if len(drop_cards) == 0 {
		return false
	}

	card_type := drop_cards[0].CardType
	for _, drop_card := range drop_cards {
		if drop_card.CardType != card_type {
			return false
		}
	}
	return true
}

func IsStraight(nums []int) (is_straight bool) {
	num_len := len(nums)
	is_straight = false

	if num_len < 2 {
		return
	}

	//获取牌的权重3-15
	weight_arr := [16]bool{}
	for _, num := range nums {
		weight := num
		if weight == 16 {
			return false
		}
		weight_arr[weight] = true
	}

	weights := make([]int, 0)
	for weight, b := range weight_arr{
		if b {
			weights = append(weights, weight)
		}
	}

	//判断是否连续
	if len(weights) < 2 {
		return
	}
	prev_weight := weights[0] - 1
	for _, weight := range weights {
		if weight - prev_weight != 1 {
			return
		}
		prev_weight = weight
	}
	is_straight = true
	return
}

//获取A-K，王的权重
func GetWeightByCardNum(num int) (weight int) {
	if num < 3 {
		return num + 13
	}else if num <= 13 {
		return num
	}else {
		return 16
	}
	return 0
}

func GetStraightWeight(nums []int) (weight int) {
	weight = 18
	for _, num := range nums {
		w := num
		if w < weight {
			weight = w
		}
	}
	return
}
