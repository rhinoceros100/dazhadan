package card

func GetCardsType(drop_cards []*Card, is_last_cards bool) (cards_type int, plane_num int, weight int) {
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

	return
}

//获取一组牌中数字相同的牌的数量
func GetSameCardsNum(drop_cards []*Card) (max_same_card int, same_card_nums []int) {
	arr := [15]int{}
	for _, drop_card := range drop_cards {
		arr[drop_card.CardNo] ++
	}

	same_card_nums = make([]int, 0)
	max_same_card = 0
	for i, num := range arr {
		if num > max_same_card {
			same_card_nums = same_card_nums[0:0]
			max_same_card = num
			same_card_nums = append(same_card_nums, i)
		}else if num == max_same_card {
			same_card_nums = append(same_card_nums, i)
		}
	}
	return
}
