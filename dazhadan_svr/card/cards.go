package card

import (
	"sort"
	"dazhadan/dazhadan_svr/util"
	"dazhadan/dazhadan_svr/log"
)

type Cards struct {
	Data 	[]*Card			`json:"data"`
}

//创建一个Cards对象
func NewCards(allCard ...*Card) *Cards{
	cards := &Cards{
		Data :	make([]*Card, 0),
	}
	for _, card := range allCard {
		cards.AddAndSort(card)
	}
	return cards
}

func CreateNewCards(cardSlice []*Card) *Cards{
	newCardSlice := make([]*Card, 0)
	for _, new_card := range cardSlice {
		newCardSlice = append(newCardSlice, new_card)
	}
	return &Cards{
		Data: newCardSlice,
	}
}

func CopyCards(cardSlice []*Card) []*Card{
	newCardSlice := make([]*Card, 0)
	for _, card := range cardSlice {
		new_card := &Card{
			CardNo:card.CardNo,
			CardType:card.CardType,
			CardId:card.CardId,
			Weight:card.Weight,
		}
		newCardSlice = append(newCardSlice, new_card)
	}
	return newCardSlice
}

//获取cards的数据
func (cards *Cards) GetData() []*Card {
	return cards.Data
}

//获取第idx个牌
func (cards *Cards) At(idx int) *Card {
	if idx >= cards.Len() {
		return nil
	}
	return cards.Data[idx]
}

//cards的长度，牌数
func (cards *Cards) Len() int {
	return len(cards.Data)
}

//比较指定索引对应的两个牌的大小
func (cards *Cards) Less(i, j int) bool {
	cardI := cards.At(i)
	cardJ := cards.At(j)
	if cardI == nil || cardJ == nil{
		return false
	}

	if cardI.CardNo != cardJ.CardNo {
		return cardI.Weight > cardJ.Weight
	}

	if cardI.CardType > cardJ.CardType {
		return true
	}
	return false
}

//交换索引为，j的两个数据
func (cards *Cards) Swap(i, j int) {
	if i == j {
		return
	}
	length := cards.Len()
	if i >= length || j >= length {
		return
	}
	swap := cards.At(i)
	cards.Data[i] = cards.At(j)
	cards.Data[j] = swap
}

//追加一张牌
func (cards *Cards) AppendCard(card *Card) {
	if card == nil {
		return
	}
	cards.Data = append(cards.Data, card)
}

//增加一张牌并排序
func (cards *Cards) AddAndSort(card *Card){
	if card == nil {
		return
	}
	cards.AppendCard(card)
	cards.Sort()//default sort
}

//追加一个cards对象
func (cards *Cards) AppendCards(other *Cards) {
	cards.Data = append(cards.Data, other.Data...)
}

//取走一组指定的牌，并返回成功或者失败
func (cards *Cards) TakeWayGroup(drop_cards []*Card) bool {
	if drop_cards == nil || len(drop_cards) == 0 {
		return true
	}

	//查找这些牌是否都在cards中
	cards_len := len(drop_cards)
	same_num := 0
	for _, drop_card := range drop_cards {
		for _, card := range cards.Data {
			if card.SameAs(drop_card) {
				same_num ++
				break
			}
		}
	}
	//log.Debug("cards_len:", cards_len, ", same_num:", same_num)
	if cards_len != same_num {
		return false
	}

	//log.Debug(cards)
	//log.Debug(drop_cards)
	//删除相应的牌
	for _, drop_card := range drop_cards {
		//log.Debug("compare", drop_card)
		for idx, card := range cards.Data {
			if card.SameAs(drop_card) {
				//log.Debug("same", drop_card)
				cards.Data = append(cards.Data[0:idx], cards.Data[idx + 1:]...)
				break
			}
		}
	}
	log.Debug("left==", cards)
	return true
}

//取走一张指定的牌，并返回成功或者失败
func (cards *Cards) TakeWay(drop *Card) bool {
	if drop == nil {
		return true
	}
	for idx, card := range cards.Data {
		if card.SameAs(drop) {
			cards.Data = append(cards.Data[0:idx], cards.Data[idx+1:]...)
			return true
		}
	}
	return false
}

//取走第一张牌
func (cards *Cards) PopFront() *Card {
	if cards.Len() == 0 {
		return nil
	}
	card := cards.At(0)
	cards.Data = cards.Data[1:]
	return card
}

//取走最后一张牌
func (cards *Cards) Tail(num int) []*Card {
	cards_len := cards.Len()
	if cards_len < num {
		return nil
	}
	return cards.Data[cards_len - num:]
}

//随机取走一张牌
func (cards *Cards) RandomTakeWayOne() *Card {
	length := cards.Len()
	if length == 0 {
		return nil
	}
	idx := util.RandomN(length)
	card := cards.At(idx)
	cards.Data = append(cards.Data[0:idx], cards.Data[idx+1:]...)
	return card
}

//清空牌
func (cards *Cards) Clear() {
	cards.Data = cards.Data[0:0]
}

//排序
func (cards *Cards)Sort() {
	sort.Sort(cards)
}

//是否是一样的牌组
func (cards *Cards) SameAs(other *Cards) bool {
	if cards == nil || other == nil {
		return false
	}

	length := other.Len()
	if cards.Len() != length {
		return false
	}

	for idx := 0; idx < length; idx++ {
		if !cards.At(idx).SameAs(other.At(idx)) {
			return false
		}
	}
	return true
}

func (cards *Cards) String() string {
	str := ""
	for _, card := range cards.Data{
		str += card.String() + ","
	}
	return str
}

//检查是否存在子集subCards
func (cards *Cards) hasCards(subCards *Cards) bool {
	if subCards.Len() == 0 {
		return true
	}
	tmpCards := CreateNewCards(cards.GetData())
	for _, subCard := range subCards.Data {
		if !tmpCards.HasCard(subCard) {
			return false
		}
		tmpCards.TakeWay(subCard)
	}

	return true
}

func (cards *Cards) HasCard(card *Card) bool{
	for _, tmp := range cards.Data {
		if tmp.SameAs(card) {
			return true
		}
	}
	return false
}
