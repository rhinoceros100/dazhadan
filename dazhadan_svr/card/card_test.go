package card

import (
	"testing"
	"github.com/bmizerany/assert"
)

func TestCard(t *testing.T) {
	card1 := &Card{
		CardType: CardType_Diamond,
		CardNo: 1,
	}
	card2 := &Card{
		CardType: CardType_Club,
		CardNo: 1,
	}
	card3 := &Card{
		CardType: CardType_Diamond,
		CardNo: 1,
	}

	card4 := &Card{
		CardType: CardType_Club,
		CardNo: 2,
	}
	assert.Equal(t, card1.SameAs(card3), true)
	assert.Equal(t, card2.Next(), card4)
	assert.Equal(t, card2.SameCardNoAs(card3), true)
	assert.Equal(t, card2.SameCardTypeAs(card4), true)
}
