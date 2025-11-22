package main

import (
	"crypto/rand"
	"fmt"
	"image/color"
	"log"
	"math/big"
	"slices"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// ==========================================
// РАЗДЕЛ: КОНСТАНТЫ И ТИПЫ КАРТ
// ==========================================

const (
	SuitSpades int = iota
	SuitClubs
	SuitDiamonds
	SuitHearts
)

// Вспомогательные константы рангов для оценки
const (
	Rank2 = iota
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
	Rank9
	Rank10
	RankJ
	RankQ
	RankK
	RankA
)

// ==========================================
// РАЗДЕЛ: КРИПТОГРАФИЯ
// ==========================================

type PlayerCrypto struct {
	ID    int
	C     *big.Int
	D     *big.Int
	Cards []*big.Int // Зашифрованные карты на руках
}

func NewPlayerCrypto(id int, p *big.Int) *PlayerCrypto {
	pMinus1 := new(big.Int).Sub(p, big.NewInt(1))
	gcd := new(big.Int)
	for {
		c, err := rand.Prime(rand.Reader, 127)
		if err != nil {
			log.Fatalf("Err: %s", err.Error())
		}
		gcd.GCD(nil, nil, c, pMinus1)
		if gcd.Cmp(big.NewInt(1)) == 0 {
			d := new(big.Int).ModInverse(c, pMinus1)
			return &PlayerCrypto{ID: id, C: c, D: d}
		}
	}
}

func (pl PlayerCrypto) DecryptCard(c, p *big.Int) *big.Int {
	return new(big.Int).Exp(c, pl.D, p)
}

func (pl *PlayerCrypto) AddCard(c *big.Int) {
	pl.Cards = append(pl.Cards, c)
}

type Deck struct {
	Cards []*big.Int
}

func NewDeck() *Deck {
	d := &Deck{Cards: make([]*big.Int, 52)}
	for i := range d.Cards {
		// Используем большие числа для ID карт
		newCard := big.NewInt(int64(50000000) + int64(i)*1000)
		d.Cards[i] = newCard
	}
	return d
}

func CopyDeck(other *Deck) *Deck {
	d := &Deck{Cards: make([]*big.Int, len(other.Cards))}
	for i := range other.Cards {
		d.Cards[i] = new(big.Int).Set(other.Cards[i])
	}
	return d
}

func (d *Deck) Shuffle() {
	for i := len(d.Cards) - 1; i > 0; i-- {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		idx := j.Int64()
		d.Cards[i], d.Cards[idx] = d.Cards[idx], d.Cards[i]
	}
}

func (d *Deck) Encrypt(c, p *big.Int) {
	for i := range d.Cards {
		d.Cards[i] = new(big.Int).Exp(d.Cards[i], c, p)
	}
}

func (d Deck) FindCard(num *big.Int) int {
	return slices.IndexFunc(d.Cards, func(n *big.Int) bool { return n.Cmp(num) == 0 })
}

type Table struct {
	ReferenceDeck *Deck
	CurrentDeck   *Deck
	BoardCards    []*big.Int
	PlayersCrypto []*PlayerCrypto
	P             *big.Int
}

func NewTable(n int) (*Table, error) {
	p, _ := rand.Prime(rand.Reader, 128)
	refDeck := NewDeck()
	players := make([]*PlayerCrypto, n)
	for i := range players {
		players[i] = NewPlayerCrypto(i, p)
	}
	currentDeck := CopyDeck(refDeck)
	return &Table{
		ReferenceDeck: refDeck,
		CurrentDeck:   currentDeck,
		PlayersCrypto: players,
		P:             p,
		BoardCards:    []*big.Int{},
	}, nil
}

func (t *Table) FullShuffle() {
	deck := t.CurrentDeck
	for _, p := range t.PlayersCrypto {
		// Игрок шифрует и мешает
		deckToEncrypt := CopyDeck(deck)
		deckToEncrypt.Encrypt(p.C, t.P)
		deckToEncrypt.Shuffle()
		deck = deckToEncrypt
	}
	t.CurrentDeck = deck
}

func (t *Table) DealCardToPlayer(playerIndex int) {
	if len(t.CurrentDeck.Cards) == 0 {
		return
	}
	card := t.CurrentDeck.Cards[0]
	t.CurrentDeck.Cards = t.CurrentDeck.Cards[1:]

	decryptedLayer := card
	for i, p := range t.PlayersCrypto {
		if i != playerIndex {
			decryptedLayer = p.DecryptCard(decryptedLayer, t.P)
		}
	}
	t.PlayersCrypto[playerIndex].AddCard(decryptedLayer)
}

func (t *Table) DealBoardCard() {
	if len(t.CurrentDeck.Cards) == 0 {
		return
	}
	card := t.CurrentDeck.Cards[0]
	t.CurrentDeck.Cards = t.CurrentDeck.Cards[1:]

	val := card
	for _, p := range t.PlayersCrypto {
		val = p.DecryptCard(val, t.P)
	}
	t.BoardCards = append(t.BoardCards, val)
}

func (t *Table) ResolvePlayerCard(playerIndex int, cardIdx int) int {
	if cardIdx >= len(t.PlayersCrypto[playerIndex].Cards) {
		return -1
	}
	encryptedCard := t.PlayersCrypto[playerIndex].Cards[cardIdx]
	plainCard := t.PlayersCrypto[playerIndex].DecryptCard(encryptedCard, t.P)
	return t.ReferenceDeck.FindCard(plainCard)
}

// ==========================================
// РАЗДЕЛ: ИГРОВАЯ ЛОГИКА И EVALUATOR
// ==========================================

type GamePhase int

const (
	PhasePreFlop GamePhase = iota
	PhaseFlop
	PhaseTurn
	PhaseRiver
	PhaseShowdown
)

type PlayerState struct {
	Chips    int
	Bet      int
	Folded   bool
	AllIn    bool
	IsWinner bool
	HandDesc string
}

type Game struct {
	Table        *Table
	Phase        GamePhase
	ActiveView   int // Чьими глазами мы смотрим (для демо)
	InfoText     string
	
	// Betting State
	PlayerStates []PlayerState
	Pot          int
	CurrentWager int // Текущая ставка, которую нужно уравнять
	ActionPlayer int // Индекс игрока, который должен ходить
	LastRaiser   int // Кто последний повысил (чтобы знать, когда круг закончен)
	GameOver     bool
}

func NewGame() *Game {
	nPlayers := 4
	tbl, _ := NewTable(nPlayers)
	tbl.FullShuffle()

	// Инициализация фишек
	pStates := make([]PlayerState, nPlayers)
	for i := range pStates {
		pStates[i] = PlayerState{Chips: 1000, Folded: false}
	}

	// Раздача Pre-flop
	for i := 0; i < 2; i++ {
		for pIdx := range tbl.PlayersCrypto {
			tbl.DealCardToPlayer(pIdx)
		}
	}

	g := &Game{
		Table:        tbl,
		Phase:        PhasePreFlop,
		ActiveView:   0,
		PlayerStates: pStates,
		ActionPlayer: 0, // Начинает игрок 0
		LastRaiser:   0,
		InfoText:     "Pre-Flop: Bets are open.",
	}
	return g
}

func (g *Game) Update() error {
	if g.GameOver {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			// Кнопка Restart
			if x > ScreenWidth-120 && x < ScreenWidth-10 && y > ScreenHeight-50 && y < ScreenHeight-10 {
				*g = *NewGame()
			}
		}
		return nil
	}

	// Обработка кликов
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		// Кнопка Switch View (для отладки/демо)
		if x > 10 && x < 150 && y > ScreenHeight-50 && y < ScreenHeight-10 {
			g.ActiveView = (g.ActiveView + 1) % len(g.Table.PlayersCrypto)
		}

		// Логика кнопок действий (только если сейчас ход ActiveView)
		if g.Phase != PhaseShowdown && g.ActionPlayer == g.ActiveView {
			// Fold (300, 550)
			if inRect(x, y, 300, 550, 80, 40) {
				g.PerformAction("fold")
			}
			// Call/Check (390, 550)
			if inRect(x, y, 390, 550, 80, 40) {
				g.PerformAction("call")
			}
			// Raise (480, 550)
			if inRect(x, y, 480, 550, 80, 40) {
				g.PerformAction("raise")
			}
		}
	}
	return nil
}

func inRect(x, y, rx, ry, rw, rh int) bool {
	return x >= rx && x <= rx+rw && y >= ry && y <= ry+rh
}

func (g *Game) PerformAction(action string) {
	idx := g.ActionPlayer
	player := &g.PlayerStates[idx]

	switch action {
	case "fold":
		player.Folded = true
		g.InfoText = fmt.Sprintf("Player %d Folded", idx)
	
	case "call":
		needed := g.CurrentWager - player.Bet
		if needed > player.Chips {
			needed = player.Chips // All in (упрощенно)
		}
		player.Chips -= needed
		player.Bet += needed
		g.Pot += needed
		g.InfoText = fmt.Sprintf("Player %d Called %d", idx, needed)

	case "raise":
		// Упрощенный рейз: всегда +50 к текущей ставке или all-in
		raiseAmt := 50
		needed := (g.CurrentWager + raiseAmt) - player.Bet
		if needed > player.Chips {
			// Если фишек не хватает на рейз, делаем просто колл/олл-ин (упрощение)
			needed = player.Chips
			player.Chips -= needed
			player.Bet += needed
			g.Pot += needed
			if player.Bet > g.CurrentWager {
				g.CurrentWager = player.Bet
				g.LastRaiser = idx // Круг ставок обновляется
			}
		} else {
			player.Chips -= needed
			player.Bet += needed
			g.Pot += needed
			g.CurrentWager = player.Bet
			g.LastRaiser = idx
		}
		g.InfoText = fmt.Sprintf("Player %d Raised", idx)
	}

	g.NextTurn()
}

func (g *Game) NextTurn() {
	// Находим следующего игрока, который не сфолдил
	nextIdx := (g.ActionPlayer + 1) % len(g.Table.PlayersCrypto)
	count := 0
	for g.PlayerStates[nextIdx].Folded {
		nextIdx = (nextIdx + 1) % len(g.Table.PlayersCrypto)
		count++
		if count > len(g.Table.PlayersCrypto) { break } // Все сфолдили?
	}

	// Проверка на завершение круга ставок
	// Круг завершен, если мы вернулись к LastRaiser и ставки уравнены
	isRoundComplete := false
	if nextIdx == g.LastRaiser {
		// Проверяем, все ли активные игроки уравняли ставку
		allMatched := true
		for i, p := range g.PlayerStates {
			if !p.Folded && p.Chips > 0 && p.Bet != g.CurrentWager {
				allMatched = false
				// Особый случай: если LastRaiser - это единственный активный, то все ок
				if i == g.LastRaiser {
					// Но если он повысил, другие должны ответить. 
					// Если мы пришли к нему снова, значит другие либо ответили, либо сфолдили.
				}
			}
		}
		if allMatched {
			isRoundComplete = true
		}
	}
	
	// Проверка: остался только один игрок?
	activeCount := 0
	winnerIdx := -1
	for i, p := range g.PlayerStates {
		if !p.Folded {
			activeCount++
			winnerIdx = i
		}
	}
	
	if activeCount == 1 {
		g.PlayerStates[winnerIdx].IsWinner = true
		g.PlayerStates[winnerIdx].Chips += g.Pot
		g.Pot = 0
		g.InfoText = fmt.Sprintf("Player %d wins by default!", winnerIdx)
		g.GameOver = true
		g.Phase = PhaseShowdown
		return
	}

	if isRoundComplete {
		g.NextPhase()
	} else {
		g.ActionPlayer = nextIdx
	}
}

func (g *Game) NextPhase() {
	// Сброс ставок для новой фазы
	for i := range g.PlayerStates {
		g.PlayerStates[i].Bet = 0
	}
	g.CurrentWager = 0
	// Начинает игрок после дилера (0) -> 1 (упрощенно)
	// Или первый активный
	startP := 0
	for g.PlayerStates[startP].Folded {
		startP = (startP + 1) % 4
	}
	g.ActionPlayer = startP
	g.LastRaiser = startP

	switch g.Phase {
	case PhasePreFlop:
		g.Table.DealBoardCard()
		g.Table.DealBoardCard()
		g.Table.DealBoardCard()
		g.Phase = PhaseFlop
		g.InfoText = "Flop"
	case PhaseFlop:
		g.Table.DealBoardCard()
		g.Phase = PhaseTurn
		g.InfoText = "Turn"
	case PhaseTurn:
		g.Table.DealBoardCard()
		g.Phase = PhaseRiver
		g.InfoText = "River"
	case PhaseRiver:
		g.EvaluateWinner()
		g.Phase = PhaseShowdown
		g.GameOver = true
		g.InfoText = "Showdown!"
	}
}

// --- Логика оценки рук ---

type HandRank int

const (
	HighCard HandRank = iota
	Pair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
)

func (g *Game) EvaluateWinner() {
	bestRank := -1
	var winners []int
	
	boardIDs := make([]int, len(g.Table.BoardCards))
	for i, bc := range g.Table.BoardCards {
		boardIDs[i] = g.Table.ReferenceDeck.FindCard(bc)
	}

	for i := range g.PlayerStates {
		if g.PlayerStates[i].Folded {
			continue
		}
		
		// Для оценки нам нужно расшифровать карты игрока
		// В ментальном покере в фазе Showdown все обмениваются ключами.
		// Тут мы используем ResolvePlayerCard (который имеет доступ ко всем ключам в памяти Table)
		c1 := g.Table.ResolvePlayerCard(i, 0)
		c2 := g.Table.ResolvePlayerCard(i, 1)
		
		allCards := append([]int{c1, c2}, boardIDs...)
		rank, desc := EvaluateHand(allCards)
		g.PlayerStates[i].HandDesc = desc

		if int(rank) > bestRank {
			bestRank = int(rank)
			winners = []int{i}
		} else if int(rank) == bestRank {
			// Тут нужно сравнить кикеры, но для демо упростим: ничья
			winners = append(winners, i)
		}
	}

	splitPot := g.Pot / len(winners)
	winText := "Winner: "
	for _, wIdx := range winners {
		g.PlayerStates[wIdx].IsWinner = true
		g.PlayerStates[wIdx].Chips += splitPot
		winText += fmt.Sprintf("P%d ", wIdx)
	}
	g.Pot = 0
	g.InfoText = winText
}

func EvaluateHand(cardIDs []int) (HandRank, string) {
	// Преобразуем ID в ранги и масти
	type Card struct { Rank, Suit int }
	cards := make([]Card, len(cardIDs))
	for i, id := range cardIDs {
		cards[i] = Card{Rank: id / 4, Suit: id % 4}
	}
	
	// Сортировка по рангу (сверху вниз)
	sort.Slice(cards, func(i, j int) bool {
		return cards[i].Rank > cards[j].Rank
	})

	// 1. Flush
	countsSuit := make([]int, 4)
	for _, c := range cards { countsSuit[c.Suit]++ }
	isFlush := false
	for _, c := range countsSuit { if c >= 5 { isFlush = true } }

	// 2. Straight
	// Убираем дубликаты рангов для проверки стрита
	uniqueRanks := []int{}
	seen := map[int]bool{}
	for _, c := range cards {
		if !seen[c.Rank] {
			uniqueRanks = append(uniqueRanks, c.Rank)
			seen[c.Rank] = true
		}
	}
	// Спец случай туз: A, 5, 4, 3, 2
	isStraight := false
	consecutive := 0
	for i := 0; i < len(uniqueRanks)-1; i++ {
		if uniqueRanks[i] - uniqueRanks[i+1] == 1 {
			consecutive++
		} else {
			consecutive = 0
		}
		if consecutive >= 4 { isStraight = true }
	}
	// Проверка A-2-3-4-5 (A=12, 2=0)
	// (Упрощенно опустим для краткости кода, но обычный стрит ловит)

	if isFlush && isStraight { return StraightFlush, "Straight Flush" }
	
	// 3. Matches (Pairs, Trips, Quads)
	countsRank := make(map[int]int)
	for _, c := range cards { countsRank[c.Rank]++ }
	
	pairs := 0
	trips := 0
	quads := 0
	for _, count := range countsRank {
		if count == 2 { pairs++ }
		if count == 3 { trips++ }
		if count == 4 { quads++ }
	}

	if quads > 0 { return FourOfAKind, "Four of a Kind" }
	if trips > 0 && pairs > 0 { return FullHouse, "Full House" }
	if trips > 1 { return FullHouse, "Full House" } // 2 тройки это тоже фулл
	if isFlush { return Flush, "Flush" }
	if isStraight { return Straight, "Straight" }
	if trips > 0 { return ThreeOfAKind, "Three of a Kind" }
	if pairs >= 2 { return TwoPair, "Two Pair" }
	if pairs == 1 { return Pair, "Pair" }

	return HighCard, "High Card"
}

// ==========================================
// РАЗДЕЛ: ГРАФИЧЕСКИЙ ИНТЕРФЕЙС
// ==========================================

const (
	ScreenWidth  = 800
	ScreenHeight = 600
	CardWidth    = 60
	CardHeight   = 90
)

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{34, 100, 34, 255}) // Темно-зеленый стол

	// Отрисовка карт на столе
	boardStartX := ScreenWidth/2 - (5*70)/2
	for i, encryptedVal := range g.Table.BoardCards {
		cardID := g.Table.ReferenceDeck.FindCard(encryptedVal)
		drawCard(screen, float64(boardStartX+i*70), ScreenHeight/2-45, cardID, true)
	}

	// Пот
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("POT: %d", g.Pot), ScreenWidth/2-30, ScreenHeight/2+60)

	// Отрисовка игроков
	centers := []struct{ x, y float64 }{
		{ScreenWidth / 2, ScreenHeight - 130}, // P0 (Вы)
		{100, ScreenHeight / 2},               // P1 (сдвинут внутрь)
		{ScreenWidth / 2, 60},                // P2
		{ScreenWidth - 100, ScreenHeight / 2}, // P3 (сдвинут внутрь)
	}

	for i, pCrypto := range g.Table.PlayersCrypto {
		pState := g.PlayerStates[i]
		cx, cy := centers[i].x, centers[i].y

		// Статус (Fold, Winner)
		label := fmt.Sprintf("P%d: $%d", i, pState.Chips)
		if pState.Folded { label += " (Fold)" }
		if pState.IsWinner { label += " WIN!" }
		if i == g.ActiveView { label += " [YOU]" }
		if i == g.ActionPlayer && !g.GameOver { label += " <ACT>" }
		
		// Карты
		// Вычисляем ширину руки, чтобы расположить карты рядом по центру
		// Ширина карты 60 + 10px отступ = 70px на слот (кроме последней, но для симметрии считаем так)
		handSlotWidth := float64(CardWidth + 10)
		totalHandWidth := float64(len(pCrypto.Cards)) * handSlotWidth
		// Смещаем начало влево на половину общей ширины, плюс поправка на половину слота
		startX := cx - (totalHandWidth / 2) + 5 

		for cIdx := range pCrypto.Cards {
			offsetX := float64(cIdx) * handSlotWidth
			// Видим карты, если это активный вид или конец игры
			isVisible := (i == g.ActiveView) || (g.Phase == PhaseShowdown)
			if pState.Folded && !isVisible { 
				// Сфолженные карты можно скрыть или затемнить, тут просто рисуем рубашку
				isVisible = false 
			}

			cardVal := -1
			if isVisible {
				cardVal = g.Table.ResolvePlayerCard(i, cIdx)
			}
			drawCard(screen, startX+offsetX, cy-45, cardVal, isVisible)
		}
		
		// Отрисовка текущей ставки рядом с игроком
		if pState.Bet > 0 {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Bet: %d", pState.Bet), int(cx)+40, int(cy))
		}
		
		// Отрисовка комбинации на Showdown
		if g.Phase == PhaseShowdown && !pState.Folded {
			ebitenutil.DebugPrintAt(screen, pState.HandDesc, int(cx)-30, int(cy)-60)
		}

		ebitenutil.DebugPrintAt(screen, label, int(cx-30), int(cy+50))
	}

	drawUI(screen, g)
}

func drawCard(screen *ebiten.Image, x, y float64, cardID int, faceUp bool) {
	rectColor := color.RGBA{240, 240, 240, 255}
	if !faceUp {
		rectColor = color.RGBA{100, 30, 30, 255} // Рубашка
	}
	ebitenutil.DrawRect(screen, x, y, CardWidth, CardHeight, rectColor)
	
	// Рамка
	// ebitenutil.DrawLine(...) - опустим для краткости

	if faceUp && cardID >= 0 {
		txt, col := getCardTextAndColor(cardID)
		// Рисуем ранг
		text.Draw(screen, txt.Rank, basicfont.Face7x13, int(x+5), int(y+15), col)
		// Рисуем масть (буквой, чтобы работало везде)
		text.Draw(screen, txt.Suit, basicfont.Face7x13, int(x+25), int(y+50), col)
	}
}

type CardDisplay struct { Rank, Suit string }

func getCardTextAndColor(id int) (CardDisplay, color.Color) {
	ranks := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
	// Порядок мастей в ID: Spades, Clubs, Diamonds, Hearts
	
	rankIdx := id / 4
	suitIdx := id % 4
	
	rStr := "?"
	if rankIdx < len(ranks) { rStr = ranks[rankIdx] }

	sStr := "?"
	var c color.Color = color.Black

	switch suitIdx {
	case SuitSpades:
		sStr = "S" // Spades (Пики)
		c = color.Black
	case SuitClubs:
		sStr = "C" // Clubs (Трефы)
		c = color.RGBA{0, 100, 0, 255} // Темно-зеленый или черный
	case SuitDiamonds:
		sStr = "D" // Diamonds (Бубны)
		c = color.RGBA{200, 0, 0, 255} // Красный
	case SuitHearts:
		sStr = "H" // Hearts (Черви)
		c = color.RGBA{200, 0, 0, 255} // Красный
	}

	return CardDisplay{rStr, sStr}, c
}

func drawUI(screen *ebiten.Image, g *Game) {
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Phase: %s | %s", g.PhaseToString(), g.InfoText), 10, 10)

	// Кнопки действий
	if g.GameOver {
		drawButton(screen, ScreenWidth-120, ScreenHeight-50, 110, 40, "RESTART", color.RGBA{0, 100, 0, 255})
	} else {
		// Рисуем кнопки только если ход активного игрока (на которого смотрим)
		if g.ActionPlayer == g.ActiveView {
			drawButton(screen, 300, 550, 80, 40, "FOLD", color.RGBA{150, 50, 50, 255})
			
			callText := "CHECK"
			if g.CurrentWager > g.PlayerStates[g.ActionPlayer].Bet {
				callText = fmt.Sprintf("CALL %d", g.CurrentWager - g.PlayerStates[g.ActionPlayer].Bet)
			}
			drawButton(screen, 390, 550, 80, 40, callText, color.RGBA{50, 50, 150, 255})
			
			drawButton(screen, 480, 550, 80, 40, "RAISE", color.RGBA{200, 150, 0, 255})
		} else {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Waiting for Player %d...", g.ActionPlayer), 300, 560)
		}
	}
	
	// Кнопка Switch View
	drawButton(screen, 10, ScreenHeight-50, 140, 40, fmt.Sprintf("View: P%d", g.ActiveView), color.RGBA{80, 80, 80, 255})
}

func drawButton(screen *ebiten.Image, x, y, w, h int, label string, col color.Color) {
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(w), float64(h), col)
	text.Draw(screen, label, basicfont.Face7x13, x+5, y+25, color.White)
}

func (g *Game) PhaseToString() string {
	switch g.Phase {
	case PhasePreFlop: return "Pre-Flop"
	case PhaseFlop: return "Flop"
	case PhaseTurn: return "Turn"
	case PhaseRiver: return "River"
	case PhaseShowdown: return "Showdown"
	default: return "Unknown"
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Texas Hold'em Mental Poker")
	
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}