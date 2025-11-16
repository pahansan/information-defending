package poker

import (
	"crypto/rand"
	"math/big"
)

type Table struct {
	CardsDeck *Deck
	Board     *Player
	Players   []*Player
	P         *big.Int
}

func GenerateP() (*big.Int, error) {
	return rand.Prime(rand.Reader, 128)
}

func NewTable(n int) (*Table, error) {
	p, err := GenerateP()
	if err != nil {
		return nil, err
	}
	players := make([]*Player, n)
	for i := range players {
		players[i] = NewPlayer(p)
	}
	board := NewPlayer(p)
	deck := NewDeck()
	return &Table{CardsDeck: deck, Board: board, Players: players, P: p}, nil
}

func (t Table) Encrypt() *Deck {
	deck := t.CardsDeck
	for i := range t.Players {
		deck = t.Players[i].EncryptDeck(deck, t.P)
	}
	deck = t.Board.EncryptDeck(deck, t.P)
	return deck
}

func (t *Table) Deal(deck *Deck) {
	for i := range 5 {
		t.Board.AddCard(deck.Cards[i])
	}
	for i := range t.Players {
		for j := range 2 {
			t.Players[i].AddCard(deck.Cards[i*j+5])
		}
	}
}

func (t Table) DecryptPlayerN(n int) []*big.Int {
	cards := t.Board.DecryptCards(t.Players[n].Cards, t.P)
	for i := range t.Players {
		if i == n {
			continue
		}
		cards = t.Players[i].DecryptCards(cards, t.P)
	}
	return t.Players[n].DecryptCards(t.Players[n].Cards, t.P)
}

func (t Table) FindCards(cards []*big.Int) []int {
	idxs := []int{}
	for i := range cards {
		idxs = append(idxs, t.CardsDeck.FindCard(cards[i]))
	}
	return idxs
}

func (t Table) FindCardsPlayerN(n int) []int {
	return t.FindCards(t.DecryptPlayerN(n))
}

func (t Table) DecryptBoard() []*big.Int {
	cards := t.Players[0].DecryptCards(t.Board.Cards, t.P)
	for i := 1; i < len(t.Players); i++ {
		cards = t.Players[i].DecryptCards(cards, t.P)
	}
	return t.Board.DecryptCards(t.Board.Cards, t.P)
}

func (t Table) FindBoardCards() []int {
	return t.FindCards(t.DecryptBoard())
}
