package poker

import (
	"crypto/rand"
	"log"
	"math/big"
)

type Player struct {
	C     *big.Int
	D     *big.Int
	Cards []*big.Int
}

func NewPlayer(p *big.Int) *Player {
	pMinus1 := new(big.Int).Sub(p, big.NewInt(1))
	gcd := new(big.Int)
	for {
		c, err := rand.Prime(rand.Reader, 127)
		if err != nil {
			log.Fatalf("Something went wrong:%s", err.Error())
		}
		gcd.GCD(nil, nil, c, pMinus1)
		if gcd.Cmp(big.NewInt(1)) == 0 {
			d := new(big.Int).ModInverse(c, pMinus1)
			return &Player{C: c, D: d}
		}
	}
}

func (pl Player) EncryptDeck(d *Deck, p *big.Int) *Deck {
	deckToEncrypt := CopyDeck(d)
	deckToEncrypt.Encrypt(pl.C, p)
	deckToEncrypt.Shuffle()
	return deckToEncrypt
}

func (pl *Player) AddCard(c *big.Int) {
	pl.Cards = append(pl.Cards, c)
}

func (pl Player) DecryptCard(c, p *big.Int) *big.Int {
	return new(big.Int).Exp(c, pl.D, p)
}

func (pl Player) DecryptCards(cards []*big.Int, p *big.Int) []*big.Int {
	newCards := make([]*big.Int, len(cards))
	for i := range cards {
		newCards[i] = pl.DecryptCard(cards[i], p)
	}
	return newCards
}
