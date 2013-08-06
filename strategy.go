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
		WizardHardHit,
		WizardSoftHit,
	}}
	dealerStrategy = Strategy{[]Rule{
		HitIfUnder17,
		HitOnSoft17,
	}}
)

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
			return Hit // actually should double
		}

	case 7, 8, 9, 10, 1:
		switch mine.BestValue() {
		case 4, 5, 6, 7, 8, 9:
			return Hit

		case 10, 11:
			return Hit // actually should double

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
			return Skip
		}

	case 7, 8, 9, 10, 1:
		switch mine.BestValue() {
		case 13, 14, 15, 16, 17, 18:
			return Skip
		}

	}

	return Skip
}
