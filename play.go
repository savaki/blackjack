package main

import (
	"flag"
	"fmt"
)

func Play(deck *Deck, bet int) int {
	Spacer("-")

	Log("Playing Hand - %d cards remain\n", deck.Remaining())

	dealer := new(Hand)
	player := new(Hand)
	player.bet = bet

	deck.DealTo(player)
	deck.DealTo(dealer)

	deck.DealTo(player)
	deck.DealTo(dealer)

	Log("\tPlayer => %+v %+v\n", player.Values(), player)
	Log("\tDealer => %+v %+v\n", dealer.Values(), dealer)

	if dealer.BlackJack() && !player.BlackJack() {
		Log("Dealer Wins!  Blackjack!\n")
		return -1 * bet
	}

	Log("Player's turn ...\n")
	hands := player.PlayStrategy(playerStrategy, dealer, deck)

	Log("Dealer's turn ...\n")
	for decision := dealerStrategy.Decide(dealer, player); decision != Skip; {
		card := deck.DealTo(dealer)
		Log("\tDealer draws %s => %d\n", card.String(), dealer.BestValue())

		if dealer.Bust() {
			sum := 0
			for index, hand := range hands {
				Log("Hand %d: Player wins!\n", index)
				sum = sum + (hand.bet * *turbo)
			}
			return sum
		}

		decision = dealerStrategy.Decide(dealer, player)
	}

	// calculate how we did for each of the hands
	sum := 0
	for index, hand := range hands {
		if hand.BlackJack() {
			Log("Hand %d: Blackjack!\n", index)
			sum = sum + (hand.bet * *turbo * 3 / 2)

		} else if hand.Beats(dealer) {
			Log("Hand %d: Player wins!\n", index)
			sum = sum + (hand.bet * *turbo)

		} else if dealer.Beats(&hand) {
			Log("Hand %d: Dealer wins!\n", index)
			sum = sum - hand.bet

		} else {
			Log("Hand %d: Push!\n", index)
		}
	}
	return sum
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

func Log(format string, args ...interface{}) (int, error) {
	if *verbose {
		return fmt.Printf(format, args...)
	} else {
		return 0, nil
	}
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

			if result > 0 {
				Log("+%d\n", result)
			} else if result < 0 {
				Log("%d\n", result)
			}
			cash = cash + result
			hands = hands + 1
		}
	}

	Spacer("-")
	fmt.Printf("Win/Loss per Hand: %3.1f\n\n", float32(cash)/float32(hands))
}
