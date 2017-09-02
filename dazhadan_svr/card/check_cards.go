package card

func GetPaixing(drop_cards []*Card) (paixing int, plane_num int, weight int) {
	paixing = CardsType_NO
	plane_num = 0
	weight = 0
	cards_len := len(drop_cards)

	if cards_len == 1 {
		paixing = CardsType_SINGLE
		return
	}

	return
}
