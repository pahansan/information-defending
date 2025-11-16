package poker

import (
	"math/big"
	"math/rand"
	"slices"
)

const (
	TwoSpades int = iota
	TwoClubs
	TwoDiamonds
	TwoHearts
	ThreeSpades
	ThreeClubs
	ThreeDiamonds
	ThreeHearts
	FourSpades
	FourClubs
	FourDiamonds
	FourHearts
	FiveSpades
	FiveClubs
	FiveDiamonds
	FiveHearts
	SixSpades
	SixClubs
	SixDiamonds
	SixHearts
	SevenSpades
	SevenClubs
	SevenDiamonds
	SevenHearts
	EightSpades
	EightClubs
	EightDiamonds
	EightHearts
	NineSpades
	NineClubs
	NineDiamonds
	NineHearts
	TenSpades
	TenClubs
	TenDiamonds
	TenHearts
	JackSpades
	JackClubs
	JackDiamonds
	JackHearts
	QueenSpades
	QueenClubs
	QueenDiamonds
	QueenHearts
	KingSpades
	KingClubs
	KingDiamonds
	KingHearts
	AceSpades
	AceClubs
	AceDiamonds
	AceHearts
)

type Deck struct {
	Cards []*big.Int
}

func NewDeck() *Deck {
	d := &Deck{Cards: make([]*big.Int, 52)}
	for i := range d.Cards {
		for {
			newCard := big.NewInt(int64(rand.Int()) + 2)
			if d.FindCard(newCard) == -1 {
				d.Cards[i] = newCard
				break
			}
		}
	}
	return d
}

func CopyDeck(other *Deck) *Deck {
	d := NewDeck()
	for i := range other.Cards {
		d.Cards[i] = new(big.Int).Set(other.Cards[i])
	}
	return d
}

func (d *Deck) Shuffle() {
	rand.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

func (d *Deck) Encrypt(c, p *big.Int) {
	for i := range d.Cards {
		d.Cards[i] = d.Cards[i].Exp(d.Cards[i], c, p)
	}
}

func (d Deck) FindCard(num *big.Int) int {
	return slices.IndexFunc(d.Cards, func(n *big.Int) bool { return n.Cmp(num) == 0 })
}
