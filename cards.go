package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

type Suite string

const (
	Clubs    Suite = "Clubs"
	Hearts   Suite = "Hearts"
	Diamonds Suite = "Diamonds"
	Spades   Suite = "Spades"
)

type Rank struct {
	values []int
	label  string
}

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

type Hand struct {
	cards []Card
	bet   int
}

func (h *Hand) Split() []*Hand {
	if len(h.cards) != 2 {
		panic("can only split with 2 cards")
	}

	return []*Hand{
		&Hand{[]Card{h.cards[0]}, h.bet},
		&Hand{[]Card{h.cards[1]}, h.bet},
	}
}

func (h *Hand) String() string {
	labels := make([]string, 0)
	for _, card := range h.cards {
		labels = append(labels, card.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(labels, ", "))
}

// value of hand and whether its soft or hard
func (h *Hand) Values() []int {
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
	for _, value := range h.Values() {
		if value > best && value <= 21 {
			best = value
		}
	}
	return best
}

func (h *Hand) BlackJack() bool {
	for _, sum := range h.Values() {
		if sum == 21 && len(h.cards) == 2 {
			return true
		}
	}
	return false
}

func (h *Hand) Bust() bool {
	bust := len(h.Values()) == 0
	if bust {
		Log("\tBust!\n")
	}
	return bust // no values under 21
}

func (h *Hand) Beats(other *Hand) bool {
	return h.BestValue() > other.BestValue()
}

type Deck struct {
	index int
	cards []Card
}

func (d *Deck) Shuffle() {
	Spacer("=")
	Log("Shuffling")

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

func NewDeck() *Deck {
	deck := new(Deck)
	for _, rank := range ranks {
		for _, suite := range suites {
			deck.cards = append(deck.cards, Card{rank, suite})
		}
	}
	return deck
}
