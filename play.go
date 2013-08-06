package main

import (
	"flag"
	"fmt"
)

func Play(deck *Deck, bet int) int {
	Spacer("-")

	Log("Playing Hand - %d cards remain\n", deck.Remaining())

	player := new(Hand)
	dealer := new(Hand)

	deck.DealTo(player)
	deck.DealTo(dealer)

	deck.DealTo(player)
	deck.DealTo(dealer)

	Log("\tPlayer => %+v %+v\n", player.Values(), player)
	Log("\tDealer => %+v %+v\n", dealer.Values(), dealer)

	if player.BlackJack() && !dealer.BlackJack() {
		Log("Player Wins!  Blackjack!")
		return int(float32(bet)*1.5) * *turbo
	}

	if dealer.BlackJack() && !player.BlackJack() {
		Log("Dealer Wins!  Blackjack!")
		return -1 * bet
	}

	Log("Player's turn ...")
	for playerStrategy.Hit(player, dealer) {
		card := deck.DealTo(player)
		Log("\tPlayer draws %s => %d\n", card.String(), player.BestValue())
		if player.Bust() {
			Log("Dealer wins!")
			return -1 * bet
		}
	}

	Log("Dealer's turn ...")
	for dealerStrategy.Hit(dealer, player) {
		card := deck.DealTo(dealer)
		Log("\tDealer draws %s => %d\n", card.String(), dealer.BestValue())

		if dealer.Bust() {
			Log("Player wins!")
			return bet * *turbo
		}
	}

	if player.Beats(dealer) {
		Log("Player wins!")
		return bet * *turbo

	} else if dealer.Beats(player) {
		Log("Dealer wins!")
		return -1 * bet

	} else {
		Log("Push!")
		return 0
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
			cash = cash + result
			hands = hands + 1
		}
	}

	Spacer("-")
	fmt.Printf("Win/Loss per Hand: %3.1f\n\n", float32(cash)/float32(hands))
}
