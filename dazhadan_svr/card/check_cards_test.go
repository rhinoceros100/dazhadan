package card

import (
	"testing"
	"github.com/bmizerany/assert"
)

func TestGetCardsType(t *testing.T) {
	card1 := &Card{CardType: CardType_Diamond,CardNo: 1, Weight:14,}	//方片1
	card2 := &Card{CardType: CardType_Club,CardNo: 1, Weight:14,}		//梅花1
	card3 := &Card{CardType: CardType_Diamond,CardNo: 1, Weight:14,}	//方片1
	card4 := &Card{CardType: CardType_Club,CardNo: 2, Weight:15,}		//梅花2
	card5 := &Card{CardType: CardType_Heart,CardNo: 2, Weight:15,}		//红桃2
	card6 := &Card{CardType: CardType_Spade,CardNo: 2, Weight:15,}		//黑桃2
	card7 := &Card{CardType: CardType_Spade,CardNo: 2, Weight:15,}		//黑桃2
	card8 := &Card{CardType: CardType_Spade,CardNo: 2, Weight:15,}		//黑桃2
	card9 := &Card{CardType: CardType_Diamond,CardNo: 2, Weight:15,}	//方片2

	card10 := &Card{CardType: CardType_BlackJoker,CardNo: 14, Weight:16}	//小王
	card11 := &Card{CardType: CardType_BlackJoker,CardNo: 14, Weight:16}	//小王
	card12 := &Card{CardType: CardType_RedJoker,CardNo: 14, Weight:17}	//大王
	card13 := &Card{CardType: CardType_RedJoker,CardNo: 14, Weight:17}	//大王

	card14 := &Card{CardType: CardType_Diamond,CardNo: 3, Weight:3,}	//方片3
	card15 := &Card{CardType: CardType_Diamond,CardNo: 4, Weight:4,}	//方片4
	card16 := &Card{CardType: CardType_Diamond,CardNo: 5, Weight:5,}	//方片5
	card17 := &Card{CardType: CardType_Spade,CardNo: 6, Weight:6,}		//方片6
	card18 := &Card{CardType: CardType_Spade,CardNo: 7, Weight:7,}		//方片7
	card19 := &Card{CardType: CardType_Diamond,CardNo: 3, Weight:3,}	//方片3
	card20 := &Card{CardType: CardType_Diamond,CardNo: 4, Weight:4,}	//方片4
	card21 := &Card{CardType: CardType_Diamond,CardNo: 5, Weight:5,}	//方片5
	card22 := &Card{CardType: CardType_Spade,CardNo: 6, Weight:6,}		//方片6
	card23 := &Card{CardType: CardType_Spade,CardNo: 7, Weight:7,}		//方片7
	card24 := &Card{CardType: CardType_Spade,CardNo: 3, Weight:3,}		//黑桃3
	card25 := &Card{CardType: CardType_Spade,CardNo: 4, Weight:4,}		//黑桃4
	card26 := &Card{CardType: CardType_Spade,CardNo: 5, Weight:5,}		//黑桃5
	card27 := &Card{CardType: CardType_Spade,CardNo: 6, Weight:6,}		//黑桃6
	card28 := &Card{CardType: CardType_Spade,CardNo: 7, Weight:7,}		//黑桃7

	card29 := &Card{CardType: CardType_Spade,CardNo: 13, Weight:13,}		//黑桃K
	card30 := &Card{CardType: CardType_Diamond,CardNo: 13, Weight:13,}	//方片K
	card31 := &Card{CardType: CardType_Heart,CardNo: 13, Weight:13,}		//红桃K

	cards1 := make([]*Card, 0)
	cards1 = append(cards1, card1)
	cards1 = append(cards1, card2)
	cards1 = append(cards1, card3)
	cards1 = append(cards1, card4)
	drop_cards1 := CreateNewCards(cards1)

	cards2 := make([]*Card, 0)
	cards2 = append(cards2, card4)
	cards2 = append(cards2, card5)
	cards2 = append(cards2, card6)
	cards2 = append(cards2, card9)
	drop_cards2 := CreateNewCards(cards2)

	cards3 := make([]*Card, 0)
	cards3 = append(cards3, card10)
	cards3 = append(cards3, card11)
	cards3 = append(cards3, card12)
	cards3 = append(cards3, card13)
	drop_cards3 := CreateNewCards(cards3)

	cards4 := make([]*Card, 0)
	cards4 = append(cards4, card10)
	cards4 = append(cards4, card11)
	cards4 = append(cards4, card12)
	drop_cards4 := CreateNewCards(cards4)

	cards5 := make([]*Card, 0)
	//cards5 = append(cards5, card1)
	//cards5 = append(cards5, card4)
	cards5 = append(cards5, card14)
	cards5 = append(cards5, card15)
	cards5 = append(cards5, card16)
	cards5 = append(cards5, card17)
	cards5 = append(cards5, card18)
	drop_cards5 := CreateNewCards(cards5)

	cards6 := make([]*Card, 0)
	cards6 = append(cards6, card1)
	cards6 = append(cards6, card2)
	cards6 = append(cards6, card5)
	cards6 = append(cards6, card6)
	drop_cards6 := CreateNewCards(cards6)

	cards7 := make([]*Card, 0)
	cards7 = append(cards7, card4)
	cards7 = append(cards7, card5)
	cards7 = append(cards7, card6)
	cards7 = append(cards7, card7)
	cards7 = append(cards7, card8)
	cards7 = append(cards7, card9)
	cards7 = append(cards7, card10)
	cards7 = append(cards7, card11)
	cards7 = append(cards7, card12)
	cards7 = append(cards7, card13)
	drop_cards7 := CreateNewCards(cards7)

	cards8 := make([]*Card, 0)
	cards8 = append(cards8, card1)
	cards8 = append(cards8, card2)
	cards8 = append(cards8, card3)
	cards8 = append(cards8, card4)
	cards8 = append(cards8, card5)
	cards8 = append(cards8, card6)
	cards8 = append(cards8, card14)
	cards8 = append(cards8, card15)
	cards8 = append(cards8, card16)
	cards8 = append(cards8, card17)
	drop_cards8 := CreateNewCards(cards8)

	cards9 := make([]*Card, 0)
	cards9 = append(cards9, card14)
	cards9 = append(cards9, card15)
	cards9 = append(cards9, card16)
	cards9 = append(cards9, card19)
	cards9 = append(cards9, card20)
	cards9 = append(cards9, card21)
	drop_cards9 := CreateNewCards(cards9)

	cards10 := make([]*Card, 0)
	cards10 = append(cards10, card14)
	cards10 = append(cards10, card15)
	cards10 = append(cards10, card16)
	cards10 = append(cards10, card17)
	cards10 = append(cards10, card19)
	cards10 = append(cards10, card20)
	cards10 = append(cards10, card21)
	cards10 = append(cards10, card24)
	cards10 = append(cards10, card25)
	cards10 = append(cards10, card26)
	drop_cards10 := CreateNewCards(cards10)

	cards11 := make([]*Card, 0)
	cards11 = append(cards11, card14)
	cards11 = append(cards11, card15)
	cards11 = append(cards11, card16)
	cards11 = append(cards11, card17)
	cards11 = append(cards11, card18)
	cards11 = append(cards11, card19)
	cards11 = append(cards11, card20)
	cards11 = append(cards11, card21)
	cards11 = append(cards11, card22)
	cards11 = append(cards11, card23)
	cards11 = append(cards11, card24)
	cards11 = append(cards11, card25)
	cards11 = append(cards11, card26)
	cards11 = append(cards11, card27)
	cards11 = append(cards11, card28)
	drop_cards11 := CreateNewCards(cards11)

	cards12 := make([]*Card, 0)
	cards12 = append(cards12, card1)
	cards12 = append(cards12, card2)
	cards12 = append(cards12, card3)
	cards12 = append(cards12, card4)
	cards12 = append(cards12, card5)
	cards12 = append(cards12, card10)
	cards12 = append(cards12, card11)
	cards12 = append(cards12, card29)
	cards12 = append(cards12, card30)
	cards12 = append(cards12, card31)
	drop_cards12 := CreateNewCards(cards12)

	cards13 := make([]*Card, 0)
	cards13 = append(cards13, card1)
	cards13 = append(cards13, card2)
	cards13 = append(cards13, card3)
	cards13 = append(cards13, card4)
	cards13 = append(cards13, card5)
	cards13 = append(cards13, card6)
	cards13 = append(cards13, card7)
	cards13 = append(cards13, card10)
	cards13 = append(cards13, card11)
	cards13 = append(cards13, card12)
	drop_cards13 := CreateNewCards(cards13)

	t.Log(GetCardsType(drop_cards1, true, 0, 0))		//三带二 5
	t.Log(GetCardsType(drop_cards2, false, 0, 0))		//2的四炸 24
	t.Log(GetCardsType(drop_cards3, false, 0, 0))		//王炸 28
	t.Log(GetCardsType(drop_cards4, false, 0, 0))		//无牌型 0
	t.Log(GetCardsType(drop_cards5, false, 0, 0))		//顺子 2
	t.Log(GetCardsType(drop_cards6, false, 0, 0))		//无牌型 0
	t.Log(GetCardsType(drop_cards7, false, 0, 0))		//无牌型 0
	t.Log(GetCardsType(drop_cards7, true, 0, 0))		//六带五 8
	t.Log(GetCardsType(drop_cards8, false, 0, 0))		//三带二飞机 11
	t.Log(GetCardsType(drop_cards9, false, 0, 0))		//连对 4
	t.Log(GetCardsType(drop_cards10, false, 0, 0))		//三带二飞机 11
	t.Log(GetCardsType(drop_cards11, false, 0, 0))		//三带二飞机 11
	t.Log(GetCardsType(drop_cards11, true, 0, 0))		//三带二飞机 11
	t.Log(GetCardsType(drop_cards12, false, 0, 0))		//三带二飞机 11
	t.Log(GetCardsType(drop_cards13, false, 0, 0))		//三带二飞机 11
}

func TestGetSameCardsNum(t *testing.T) {
	card1 := &Card{CardType: CardType_Diamond,CardNo: 1, Weight:14,}	//方片1
	card2 := &Card{CardType: CardType_Club,CardNo: 1, Weight:14,}		//梅花1
	card3 := &Card{CardType: CardType_Diamond,CardNo: 1, Weight:14,}	//方片1
	card4 := &Card{CardType: CardType_Club,CardNo: 2, Weight:15,}		//梅花2
	card5 := &Card{CardType: CardType_Heart,CardNo: 2, Weight:15,}		//红桃2
	card6 := &Card{CardType: CardType_Spade,CardNo: 2, Weight:15,}		//黑桃2
	card7 := &Card{CardType: CardType_Spade,CardNo: 5, Weight:5,}		//黑桃5
	card8 := &Card{CardType: CardType_Spade,CardNo: 8, Weight:8,}		//黑桃8
	//card9 := &Card{CardType: CardType_Diamond,CardNo: 2,}	//方片2

	cards1 := make([]*Card, 0)
	cards1 = append(cards1, card1)
	cards1 = append(cards1, card2)
	cards1 = append(cards1, card3)
	cards1 = append(cards1, card4)

	cards2 := make([]*Card, 0)
	cards2 = append(cards2, card1)
	cards2 = append(cards2, card2)
	cards2 = append(cards2, card3)
	cards2 = append(cards2, card4)
	cards2 = append(cards2, card5)
	cards2 = append(cards2, card6)
	cards2 = append(cards2, card7)
	cards2 = append(cards2, card8)
	//cards2 = append(cards2, card9)

	t.Log(GetSameCardsNum(cards1))
	t.Log(GetSameCardsNum(cards2))
}

func TestIs510K(t *testing.T) {
	card1 := &Card{CardType: CardType_Diamond,CardNo: 5,}	//方片5
	card2 := &Card{CardType: CardType_Club,CardNo: 10,}	//梅花10
	card3 := &Card{CardType: CardType_Diamond,CardNo: 13,}	//方片K
	card4 := &Card{CardType: CardType_Diamond,CardNo: 10,}	//方片10
	card5 := &Card{CardType: CardType_Heart,CardNo: 2,}	//红桃2
	card9 := &Card{CardType: CardType_Spade,CardNo: 2,}	//黑桃2

	cards1 := make([]*Card, 0)
	cards1 = append(cards1, card1)
	cards1 = append(cards1, card2)
	cards1 = append(cards1, card3)
	t.Log(Is510K(cards1))

	cards2 := make([]*Card, 0)
	cards2 = append(cards2, card1)
	cards2 = append(cards2, card4)
	cards2 = append(cards2, card3)
	cards2 = append(cards2, card9)
	t.Log(Is510K(cards2))

	cards3 := make([]*Card, 0)
	cards3 = append(cards3, card1)
	cards3 = append(cards3, card4)
	cards3 = append(cards3, card3)
	t.Log(Is510K(cards3))

	cards4 := make([]*Card, 0)
	cards4 = append(cards4, card1)
	cards4 = append(cards4, card5)
	cards4 = append(cards4, card3)
	t.Log(Is510K(cards4))
}

func TestIsSameCardType(t *testing.T) {
	card1 := &Card{CardType: CardType_Diamond,CardNo: 1,}	//方片1
	card2 := &Card{CardType: CardType_Club,CardNo: 1,}	//梅花1
	card3 := &Card{CardType: CardType_Diamond,CardNo: 1,}	//方片1
	//card4 := &Card{CardType: CardType_Club,CardNo: 2,}	//梅花2
	card5 := &Card{CardType: CardType_Heart,CardNo: 2,}	//红桃2
	card6 := &Card{CardType: CardType_Spade,CardNo: 2,}	//黑桃2
	card7 := &Card{CardType: CardType_Spade,CardNo: 5,}	//黑桃5
	card8 := &Card{CardType: CardType_Spade,CardNo: 8,}	//黑桃8
	card9 := &Card{CardType: CardType_Spade,CardNo: 2,}	//黑桃2

	cards1 := make([]*Card, 0)
	cards1 = append(cards1, card1)
	cards1 = append(cards1, card2)
	cards1 = append(cards1, card3)
	assert.Equal(t, IsSameCardType(cards1), false)

	cards2 := make([]*Card, 0)
	cards2 = append(cards2, card6)
	cards2 = append(cards2, card7)
	cards2 = append(cards2, card8)
	cards2 = append(cards2, card9)
	assert.Equal(t, IsSameCardType(cards2), true)

	cards3 := make([]*Card, 0)
	cards3 = append(cards3, card5)
	assert.Equal(t, IsSameCardType(cards3), true)

	cards4 := make([]*Card, 0)
	assert.Equal(t, IsSameCardType(cards4), false)
}

func TestIsStraight(t *testing.T) {
	nums1 := []int{2,3,4}
	assert.Equal(t, IsStraight(nums1), false)

	nums2 := []int{3,4}
	assert.Equal(t, IsStraight(nums2), true)

	nums3 := []int{1,2,13,12}
	assert.Equal(t, IsStraight(nums3), true)
}

func TestGetStraightWeight(t *testing.T) {
	nums1 := []int{2,3,4}
	t.Log(GetStraightWeight(nums1))

	nums2 := []int{3,4}
	t.Log(GetStraightWeight(nums2))

	nums3 := []int{1,2,13,12}
	t.Log(GetStraightWeight(nums3))
}

type roon struct {
	cardsType int
	planeNum int
	weight int
}

func canCover(cardsType, planeNum, weight int, roo *roon) (canCover bool) {
	canCover = false
	if roo.cardsType == CardsType_NO {
		return cardsType != CardsType_NO
	}
	//已经出的牌型非炸弹牌型
	if roo.cardsType < 20{
		if cardsType > 20 {
			return true
		}
		//普通牌型打普通牌型必须为同一牌型，并且飞机数量必须相同
		if cardsType != roo.cardsType{
			return false
		}
		if cardsType == CardsType_STAIGHT || cardsType == CardsType_PAIRS || cardsType >= 11 {
			if planeNum != roo.planeNum {
				return false
			}
		}
		return weight > roo.weight
	}

	//更大的炸弹可以管住
	if cardsType > roo.cardsType{
		return true
	}
	return weight > roo.weight
}

func TestCanCover(t *testing.T) {
	roo1 := &roon{cardsType:0, planeNum:0, weight:3}
	assert.Equal(t, canCover(3, 3, 2, roo1), true)
	assert.Equal(t, canCover(0, 3, 2, roo1), false)

	roo2 := &roon{cardsType:CardsType_32, planeNum:0, weight:6}
	assert.Equal(t, canCover(CardsType_32, 3, 3, roo2), false)
	assert.Equal(t, canCover(CardsType_32, 3, 9, roo2), true)
	assert.Equal(t, canCover(CardsType_510K, 3, 9, roo2), true)

	roo3 := &roon{cardsType:CardsType_PLANE43, planeNum:3, weight:5}
	assert.Equal(t, canCover(CardsType_32, 3, 6, roo3), false)
	assert.Equal(t, canCover(CardsType_PLANE43, 2, 9, roo3), false)
	assert.Equal(t, canCover(CardsType_PLANE43, 3, 5, roo3), false)
	assert.Equal(t, canCover(CardsType_PLANE43, 3, 8, roo3), true)
	assert.Equal(t, canCover(CardsType_BOMB4, 2, 8, roo3), true)

	roo4 := &roon{cardsType:CardsType_TRUE510K, planeNum:0, weight:0}
	assert.Equal(t, canCover(CardsType_TRUE510K, 2, 0, roo4), false)
	assert.Equal(t, canCover(CardsType_BOMB4, 2, 8, roo4), true)

	roo5 := &roon{cardsType:CardsType_BOMB4, planeNum:0, weight:5}
	assert.Equal(t, canCover(CardsType_BOMB4, 2, 3, roo5), false)
	assert.Equal(t, canCover(CardsType_BOMB4, 2, 8, roo5), true)
	assert.Equal(t, canCover(CardsType_BOMB7, 2, 2, roo5), true)
}

func TestCheck3Plane(t *testing.T) {
	card1 := &Card{CardType: CardType_Diamond,CardNo: 1, Weight:14,}	//方片1
	card2 := &Card{CardType: CardType_Club,CardNo: 1, Weight:14,}		//梅花1
	card3 := &Card{CardType: CardType_Diamond,CardNo: 1, Weight:14,}	//方片1
	card4 := &Card{CardType: CardType_Club,CardNo: 2, Weight:15,}		//梅花2
	card5 := &Card{CardType: CardType_Heart,CardNo: 2, Weight:15,}		//红桃2
	card6 := &Card{CardType: CardType_Spade,CardNo: 2, Weight:15,}		//黑桃2
	card7 := &Card{CardType: CardType_Spade,CardNo: 2, Weight:15,}		//黑桃2
	//card8 := &Card{CardType: CardType_Spade,CardNo: 2, Weight:15,}		//黑桃2
	//card9 := &Card{CardType: CardType_Diamond,CardNo: 2, Weight:15,}	//方片2

	card14 := &Card{CardType: CardType_Diamond,CardNo: 3, Weight:3,}	//方片3
	card15 := &Card{CardType: CardType_Diamond,CardNo: 4, Weight:4,}	//方片4
	card16 := &Card{CardType: CardType_Diamond,CardNo: 5, Weight:5,}	//方片5
	//card17 := &Card{CardType: CardType_Spade,CardNo: 6, Weight:6,}		//方片6
	//card18 := &Card{CardType: CardType_Spade,CardNo: 7, Weight:7,}		//方片7
	//card19 := &Card{CardType: CardType_Diamond,CardNo: 3, Weight:3,}	//方片3
	//card20 := &Card{CardType: CardType_Diamond,CardNo: 4, Weight:4,}	//方片4
	//card21 := &Card{CardType: CardType_Diamond,CardNo: 5, Weight:5,}	//方片5
	//card22 := &Card{CardType: CardType_Spade,CardNo: 6, Weight:6,}		//方片6
	//card23 := &Card{CardType: CardType_Spade,CardNo: 7, Weight:7,}		//方片7
	//card24 := &Card{CardType: CardType_Spade,CardNo: 3, Weight:3,}		//黑桃3
	//card25 := &Card{CardType: CardType_Spade,CardNo: 4, Weight:4,}		//黑桃4
	//card26 := &Card{CardType: CardType_Spade,CardNo: 5, Weight:5,}		//黑桃5
	//card27 := &Card{CardType: CardType_Spade,CardNo: 6, Weight:6,}		//黑桃6
	//card28 := &Card{CardType: CardType_Spade,CardNo: 7, Weight:7,}		//黑桃7

	cards1 := make([]*Card, 0)
	cards1 = append(cards1, card1)
	cards1 = append(cards1, card2)
	cards1 = append(cards1, card3)
	cards1 = append(cards1, card4)
	cards1 = append(cards1, card5)
	cards1 = append(cards1, card6)
	cards1 = append(cards1, card7)
	cards1 = append(cards1, card14)
	cards1 = append(cards1, card15)
	cards1 = append(cards1, card16)

	t.Log(Check3Plane(cards1, false, CardsType_PLANE32, 2))		//三带二飞机 11
	t.Log(Check3Plane(cards1, false, CardsType_PLANE32, 0))		//三带二飞机 11
}
