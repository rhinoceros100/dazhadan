package card

import "fmt"

const INIT_CARD_NUM int = 27      	//开局发牌数量
const TYPE2_ROUND_CARD_NUM int = 7      //一次发多张每次发牌数量
const TYPE2_ROUND_TIMES int = 3      	//一次发多张发牌轮数
const TYPE2_LAST_ROUND_CARD_NUM int = 6 //一次最后一轮发牌数量
const TOTAL_CARD_NUM int = 108 		//牌的总数量

type PlayingCards struct {
	CardsInHand			*Cards		//手上的牌
}

func NewPlayingCards() *PlayingCards {
	return  &PlayingCards{
		CardsInHand: NewCards(),
	}
}

func (playingCards *PlayingCards) Reset() {
	playingCards.CardsInHand.Clear()
}

func (playingCards *PlayingCards) AddCards(cards *Cards) {
	playingCards.CardsInHand.AppendCards(cards)
	playingCards.CardsInHand.Sort()
}

//增加一张牌
func (playingCards *PlayingCards) AddCard(card *Card) {
	playingCards.CardsInHand.AddAndSort(card)
}

func (playingCards *PlayingCards) String() string{
	return fmt.Sprintf(
		"{%v}",
		playingCards.CardsInHand,
	)
}

func (playingCards *PlayingCards) Tail(num int) []*Card {
	return playingCards.CardsInHand.Tail(num)
}

//丢弃一张牌
func (playingCards *PlayingCards) DropCards(cards []*Card) bool {
	return playingCards.CardsInHand.TakeAwayGroup(cards)
}