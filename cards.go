package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

type Suite string
type Rank struct {
	values []int
	label  string
}

const (
	Clubs    Suite = "Clubs"
	Hearts   Suite = "Hearts"
	Diamonds Suite = "Diamonds"
	Spades   Suite = "Spades"
)

var (
	Ace    = Rank{[]int{1, 11}, "Ace"}
	King   = Rank{[]int{10}, "King"}
	Queen  = Rank{[]int{10}, "Queen"}
	Jack   = Rank{[]int{10}, "Jack"}
	Ten    = Rank{[]int{10}, "Ten"}
	Nine   = Rank{[]int{9}, "Nine"}
	Eight  = Rank{[]int{8}, "Eight"}
	Seven  = Rank{[]int{7}, "Seven"}
	Six    = Rank{[]int{6}, "Six"}
	Five   = Rank{[]int{5}, "Five"}
	Four   = Rank{[]int{4}, "Four"}
	Three  = Rank{[]int{3}, "Three"}
	Two    = Rank{[]int{2}, "Two"}
	suites = []Suite{Spades, Hearts, Diamonds, Clubs}
	ranks  = []Rank{Ace, King, Queen, Jack, Ten, Nine, Eight, Seven, Six, Five, Four, Three, Two}
)

type Card struct {
	rank  Rank
	suite Suite
}

func (c Card) String() string {
	return fmt.Sprintf("%s of %s", c.rank.label, c.suite)
}

type Deck struct {
	index int
	cards []Card
}

type Hand struct {
	cards []Card
}

func (h *Hand) String() string {
	labels := make([]string, 0)
	for _, card := range h.cards {
		labels = append(labels, card.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(labels, ", "))
}

// value of hand and whether its soft or hard
func (h *Hand) Value() []int {
	sums := []int{0}

	for _, card := range h.cards {
		tmp := make([]int, 0)
		for _, sum := range sums {
			for _, value := range card.rank.values {
				tmp = append(tmp, sum+value)
			}
		}
		sums = tmp
	}

	sort.Ints(sums) // sort hand values
	for index, sum := range sums {
		if sum > 21 {
			return sums[:index]
		}
	}
	return sums
}

// the best value we can make with this hand <= 21
func (h *Hand) BestValue() int {
	best := -1
	for _, value := range h.Value() {
		if value > best && value <= 21 {
			best = value
		}
	}
	return best
}

func (h *Hand) BlackJack() bool {
	for _, sum := range h.Value() {
		if sum == 21 && len(h.cards) == 2 {
			return true
		}
	}
	return false
}

func (h *Hand) Bust() bool {
	bust := len(h.Value()) == 0
	if *verbose && bust {
		fmt.Println("\tBust!")
	}
	return bust // no values under 21
}

func (d *Deck) Shuffle() {
	if *verbose {
		Spacer("=")
		fmt.Println("Shuffling")
	}

	rand.Seed(time.Now().UTC().UnixNano())
	for i1, card1 := range d.cards {
		i2 := rand.Intn(len(d.cards))
		if i1 != i2 {
			card2 := d.cards[i2]
			d.cards[i1] = card2
			d.cards[i2] = card1
		}
	}

	d.index = 0
}

func (d *Deck) DealTo(hand *Hand) *Card {
	card := d.cards[d.index]
	d.index = d.index + 1

	hand.cards = append(hand.cards, card)

	return &card
}

func (d *Deck) Remaining() int {
	return len(d.cards) - d.index
}

func (d *Deck) Dump() {
	for _, card := range d.cards {
		fmt.Println(card)
	}
}

type Rule func(mine *Hand, theirs *Hand) bool

type Strategy struct {
	rules []Rule
}

func (s Strategy) Hit(mine *Hand, theirs *Hand) bool {
	for _, rule := range s.rules {
		if rule(mine, theirs) {
			return true
		}
	}

	return false
}

func Spacer(s string) {
	if *verbose {
		fmt.Println("")
		for i := 0; i < 65; i++ {
			fmt.Print(s)
		}
		fmt.Println("\n")
	}
}

func Play(deck *Deck, bet int) int {
	Spacer("-")

	if *verbose {
		fmt.Printf("Playing Hand - %d cards remain\n", deck.Remaining())
	}

	playerStrategy := Strategy{[]Rule{
		WizardHardHit,
		WizardSoftHit,
	}}
	dealerStrategy := Strategy{[]Rule{
		HitIfUnder17,
		HitOnSoft17,
	}}

	player := new(Hand)
	dealer := new(Hand)

	deck.DealTo(player)
	deck.DealTo(dealer)

	deck.DealTo(player)
	deck.DealTo(dealer)

	if *verbose {
		fmt.Printf("\tPlayer => %+v %+v\n", player.Value(), player)
		fmt.Printf("\tDealer => %+v %+v\n", dealer.Value(), dealer)
	}

	if player.BlackJack() && !dealer.BlackJack() {
		if *verbose {
			fmt.Println("Player Wins!  Blackjack!")
		}
		return int(float32(bet)*1.5) * *turbo
	}

	if dealer.BlackJack() && !player.BlackJack() {
		if *verbose {
			fmt.Println("Dealer Wins!  Blackjack!")
		}
		return -1 * bet
	}

	if *verbose {
		fmt.Println("Player's turn ...")
	}
	for playerStrategy.Hit(player, dealer) {
		card := deck.DealTo(player)
		if *verbose {
			fmt.Printf("\tPlayer draws %s => %d\n", card.String(), player.BestValue())
		}
		if player.Bust() {
			if *verbose {
				fmt.Println("Dealer wins!")
			}
			return -1 * bet
		}
	}

	if *verbose {
		fmt.Println("Dealer's turn ...")
	}
	for dealerStrategy.Hit(dealer, player) {
		card := deck.DealTo(dealer)
		if *verbose {
			fmt.Printf("\tDealer draws %s => %d\n", card.String(), dealer.BestValue())
		}

		if dealer.Bust() {
			if *verbose {
				fmt.Println("Player wins!")
			}
			return bet * *turbo
		}
	}

	if player.Beats(dealer) {
		if *verbose {
			fmt.Println("Player wins!")
		}
		return bet * *turbo

	} else if dealer.Beats(player) {
		if *verbose {
			fmt.Println("Dealer wins!")
		}
		return -1 * bet

	} else {
		if *verbose {
			fmt.Println("Push!")
		}
		return 0
	}
}

func NewDeck() *Deck {
	deck := new(Deck)
	for _, rank := range ranks {
		for _, suite := range suites {
			deck.cards = append(deck.cards, Card{rank, suite})
		}
	}
	return deck
}

func (h *Hand) Beats(other *Hand) bool {
	return h.BestValue() > other.BestValue()
}

var (
	turbo   = flag.Int("turbo", 3, "the level of turbo boost")
	rounds  = flag.Int("rounds", 1, "the number of complete decks to play")
	verbose = flag.Bool("verbose", false, "print out the hands as they're being played")
)

func main() {
	flag.Parse()
	deck := NewDeck()

	hands := 0 // how many hands have been played
	cash := 0  // how much money do we have
	bet := 100 // bet size

	for i := 0; i < *rounds; i++ {
		deck.Shuffle()
		for deck.Remaining() > 20 {
			result := Play(deck, bet)
			cash = cash + result
			hands = hands + 1
		}
	}

	Spacer("-")
	fmt.Printf("Win/Loss per Hand: %3.1f\n\n", float32(cash)/float32(hands))
}
