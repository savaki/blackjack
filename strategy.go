package main

import (
	"fmt"
)

type Decision int

const (
	Skip Decision = iota
	Hit
	Double
	Split
)

type Rule func(mine *Hand, theirs *Hand) Decision

type Strategy struct {
	rules []Rule
}

func (s Strategy) Decide(mine *Hand, theirs *Hand) Decision {
	for _, rule := range s.rules {
		if action := rule(mine, theirs); action != Skip {
			return action
		}
	}

	return Skip
}

var (
	playerStrategy = Strategy{[]Rule{
		WizardSplit,
		WizardHardHit,
		WizardSoftHit,
	}}
	dealerStrategy = Strategy{[]Rule{
		HitIfUnder17,
		HitOnSoft17,
	}}
)

func (h *Hand) PlayStrategy(strategy Strategy, other *Hand, deck *Deck) []Hand {
	hands := make([]Hand, 0)  // the ultimate number of hands that are played
	hands = append(hands, *h) // include the current hand

	for decision := strategy.Decide(h, other); decision != Skip; {
		switch decision {
		case Hit:
			Log("\t%d => Hit!\n", h.BestValue())
			card := deck.DealTo(h)
			Log("\tPlayer Draws %+v => %d\n", card, h.BestValue())
			if h.Bust() {
				decision = Skip
			} else {
				decision = strategy.Decide(h, other)
			}

		case Double:
			Log("\t%d => Double!\n", h.BestValue())
			card := deck.DealTo(h)
			Log("\tGood luck!  Player Draws %+v => %d\n", card, h.BestValue())
			h.bet = h.bet * 2
			decision = Skip

		case Split:
			Log("\t%d => Split!\n", h.BestValue())

			// split the one hand into two hands and give each new hand a card
			split := h.Split()

			Log("\t-------\n")
			for index, hand := range split {
				deck.DealTo(hand)
				Log("\tHand %d: %2d => %s\n", index, hand.BestValue(), hand.String())
			}
			Log("\t-------\n")

			hands = make([]Hand, 0)
			for index, hand := range split {
				Log("\tHand %d: %d => %s\n", index, hand.BestValue(), hand.String())
				for _, child := range hand.PlayStrategy(strategy, other, deck) {
					hands = append(hands, child)
				}
			}
			decision = Skip
		}
	}

	Log("\t(returning %d hands)\n", len(hands))
	return hands
}

func HitIfUnder17(mine *Hand, theirs *Hand) Decision {
	if mine.BestValue() < 17 {
		if *verbose {
			fmt.Println("\t< 17 - Hit")
		}
		return Hit

	} else {
		return Skip
	}
}

func HitOnSoft17(mine *Hand, theirs *Hand) Decision {
	if mine.BestValue() == 17 && len(mine.Values()) == 2 {
		if *verbose {
			fmt.Println("\tSoft 17 - Hit")
		}
		return Hit

	} else {
		return Skip
	}
}

// http://wizardofodds.com/games/blackjack/
func WizardHardHit(mine *Hand, theirs *Hand) Decision {
	switch theirs.cards[0].rank.values[0] {

	case 2, 3, 4, 5, 6:
		switch mine.BestValue() {
		case 4, 5, 6, 7, 8:
			return Hit

		case 9, 10, 11:
			if len(mine.cards) == 2 {
				return Double // actually should double
			} else {
				return Hit
			}
		}

	case 7, 8, 9, 10, 1:
		switch mine.BestValue() {
		case 4, 5, 6, 7, 8, 9:
			return Hit

		case 10, 11:
			if len(mine.cards) == 2 {
				return Double // actually should double
			} else {
				return Hit
			}

		case 12, 13, 14, 15, 16:
			return Hit
		}

	}

	return Skip
}

func WizardSoftHit(mine *Hand, theirs *Hand) Decision {
	if len(mine.Values()) == 1 {
		return Skip // only check for soft hits
	}

	switch theirs.cards[0].rank.values[0] {

	case 2, 3, 4, 5, 6:
		switch mine.BestValue() {
		case 13, 14, 15:
			return Hit
		}

	case 7, 8, 9, 10, 1:
		switch mine.BestValue() {
		case 13, 14, 15, 16, 17, 18:
			return Hit
		}

	}

	return Skip
}

func WizardSplit(mine *Hand, theirs *Hand) Decision {
	if len(mine.Values()) == 2 {
		return Skip // only check if we have 2 cards
	}

	if mine.cards[0].rank.values[0] != mine.cards[1].rank.values[0] {
		return Skip // these aren't the same cards
	}

	switch theirs.cards[0].rank.values[0] {

	case 2, 3, 4, 5, 6:
		switch mine.cards[0].rank.values[0] {
		case 2, 3, 6, 7, 9:
			return Split

		case 8, 1:
			return Split
		}

	case 7, 8, 9, 10, 1:
		switch mine.cards[0].rank.values[0] {
		case 8, 1:
			return Split
		}

	}

	return Skip
}
