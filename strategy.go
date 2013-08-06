package main

import (
	"fmt"
)

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

func HitIfUnder17(mine *Hand, theirs *Hand) bool {
	if mine.BestValue() < 17 {
		if *verbose {
			fmt.Println("\t< 17 - Hit")
		}
		return true

	} else {
		return false
	}
}

func HitOnSoft17(mine *Hand, theirs *Hand) bool {
	if mine.BestValue() == 17 && len(mine.Value()) == 2 {
		if *verbose {
			fmt.Println("\tSoft 17 - Hit")
		}
		return true

	} else {
		return false
	}
}

// http://wizardofodds.com/games/blackjack/
func WizardHardHit(mine *Hand, theirs *Hand) bool {
	switch theirs.cards[0].rank.values[0] {

	case 2, 3, 4, 5, 6:
		switch mine.BestValue() {
		case 4, 5, 6, 7, 8:
			return true

		case 9, 10, 11:
			return true // actually should double
		}

	case 7, 8, 9, 10, 1:
		switch mine.BestValue() {
		case 4, 5, 6, 7, 8, 9:
			return true

		case 10, 11:
			return true // actually should double

		case 12, 13, 14, 15, 16:
			return true
		}

	}

	return false
}

func WizardSoftHit(mine *Hand, theirs *Hand) bool {
	if len(mine.Value()) == 1 {
		return false // only check for soft hits
	}

	switch theirs.cards[0].rank.values[0] {

	case 2, 3, 4, 5, 6:
		switch mine.BestValue() {
		case 13, 14, 15:
			return true
		}

	case 7, 8, 9, 10, 1:
		switch mine.BestValue() {
		case 13, 14, 15, 16, 17, 18:
			return true
		}

	}

	return false
}
