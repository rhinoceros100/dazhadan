package card

import "dazhadan/dazhadan_svr/util"

const CARD_PACKS int = 2
const SHUFFLE_TIMES int = 2

type Pool struct {
	cards *Cards
}

func NewPool() *Pool {
	pool := &Pool{
		cards:	NewCards(),
	}
	return pool
}

func (pool *Pool) generate() {
	for cardNo := 1 ; cardNo <= 13; cardNo ++ {
		for num := 0; num < CARD_PACKS; num ++ {
			card_diamond := &Card{
				CardType_Diamond,
				cardNo,
				0,
				0,
			}
			card_diamond.MakeIDWeight(num)

			card_club := &Card{
				CardType_Club,
				cardNo,
				0,
				0,
			}
			card_club.MakeIDWeight(num)

			card_heart := &Card{
				CardType_Heart,
				cardNo,
				0,
				0,
			}
			card_heart.MakeIDWeight(num)

			card_spade := &Card{
				CardType_Spade,
				cardNo,
				0,
				0,
			}
			card_spade.MakeIDWeight(num)

			pool.cards.AppendCard(card_diamond)
			pool.cards.AppendCard(card_club)
			pool.cards.AppendCard(card_heart)
			pool.cards.AppendCard(card_spade)
		}
	}

	for num := 0; num < CARD_PACKS; num ++ {
		card_blackJoker := &Card{
			CardType_BlackJoker,
			14,
			0,
			0,
		}
		card_blackJoker.MakeIDWeight(num)

		card_redjoker := &Card{
			CardType_RedJoker,
			14,
			0,
			0,
		}
		card_redjoker.MakeIDWeight(num)

		pool.cards.AppendCard(card_blackJoker)
		pool.cards.AppendCard(card_redjoker)
	}
}

func (pool *Pool) ReGenerate() {
	pool.cards.Clear()
	pool.generate()
	pool.shuffle()
}

//洗牌，打乱牌
func (pool *Pool) shuffle() {
	length := pool.cards.Len()
	for cnt := 0; cnt < length * SHUFFLE_TIMES; cnt++ {
		i := util.RandomN(length)
		j := util.RandomN(length)
		pool.cards.Swap(i, j)
	}
}

func (pool *Pool) PopFront() *Card {
	return pool.cards.PopFront()
}

func (pool *Pool) At(idx int) *Card {
	return pool.cards.At(idx)
}

func (pool *Pool) GetCardNum() int {
	return pool.cards.Len()
}
