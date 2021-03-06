package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	main()
}

func TestValueOfEmptyHand(t *testing.T) {
	hand := &Hand{[]Card{}, 0}
	if len(hand.Values()) != 1 {
		t.Fatal("expected 1 value")
	}
	if hand.Values()[0] != 0 {
		t.Fatal("expected value of 0")
	}
}

func TestValueOfAceHand(t *testing.T) {
	hand := &Hand{[]Card{
		Card{Ace, Spades},
	}, 0}

	if len(hand.Values()) != 2 {
		t.Fatal("expected 2 values")
	}
	if hand.Values()[0] != 1 {
		t.Fatalf("expected 1; got %d\n", hand.Values()[0])
	}
	if hand.Values()[1] != 11 {
		t.Fatalf("expected 11; got %d\n", hand.Values()[1])
	}
}

func TestValueOfKingHand(t *testing.T) {
	hand := &Hand{[]Card{
		Card{King, Spades},
	}, 0}
	if len(hand.Values()) != 1 {
		t.Fatal("expected 1 value")
	}
	if hand.Values()[0] != 10 {
		t.Fatalf("expected 10; got %d\n", hand.Values()[0])
	}
}

func TestValueOfKingKingHand(t *testing.T) {
	hand := &Hand{[]Card{
		Card{King, Spades},
		Card{King, Clubs},
	}, 0}
	if len(hand.Values()) != 1 {
		t.Fatal("expected 1 value")
	}
	if hand.Values()[0] != 20 {
		t.Fatalf("expected 20; got %d\n", hand.Values()[0])
	}
}

func TestValueOfAceKingHand(t *testing.T) {
	hand := &Hand{[]Card{
		Card{Ace, Spades},
		Card{King, Clubs},
	}, 0}
	if len(hand.Values()) != 2 {
		t.Fatal("expected 1 value")
	}
	if hand.Values()[0] != 11 {
		t.Fatalf("expected 11; got %d\n", hand.Values()[0])
	}
	if hand.Values()[1] != 21 {
		t.Fatalf("expected 21; got %d\n", hand.Values()[1])
	}
}

func TestValueOfAceAceHand(t *testing.T) {
	hand := &Hand{[]Card{
		Card{Ace, Spades},
		Card{Ace, Clubs},
	}, 0}
	if len(hand.Values()) != 3 {
		t.Fatalf("expected 3 values; got %#v\n", hand.Values())
	}
	if hand.Values()[0] != 2 {
		t.Fatalf("expected 2; got %d\n", hand.Values()[0])
	}
	if hand.Values()[1] != 12 {
		t.Fatalf("expected 12; got %d\n", hand.Values()[1])
	}
	if hand.Values()[2] != 12 {
		t.Fatalf("expected 12; got %d\n", hand.Values()[2])
	}
}

func TestBlackJackOn15(t *testing.T) {
	hand := &Hand{[]Card{
		Card{King, Diamonds},
		Card{Five, Diamonds},
	}, 0}

	if hand.BlackJack() {
		t.Fatal("should not be blackjack")
	}
}

func TestBlackJack(t *testing.T) {
	hand := &Hand{[]Card{
		Card{King, Diamonds},
		Card{Ace, Diamonds},
	}, 0}

	if !hand.BlackJack() {
		t.Fatal("expected blackjack")
	}
}

func TestSplit(t *testing.T) {
	hand := &Hand{[]Card{
		Card{King, Diamonds},
		Card{Ace, Diamonds},
	}, 0}

	split := hand.Split()
	if len(split[0].cards) != 1 {
		t.Fatalf("expected 1 card in split[0]; actual was %d\n", len(split[0].cards))
	}
	if split[0].cards[0].rank.label != hand.cards[0].rank.label {
		t.Fatalf("expected first card first; actual was %s\n", split[0].cards[0].rank.label)
	}
	if split[1].cards[0].rank.label != hand.cards[1].rank.label {
		t.Fatalf("expected second card second; actual was %s\n", split[1].cards[0].rank.label)
	}
	if len(split[1].cards) != 1 {
		t.Fatalf("expected 1 card in split[1]; actual was %d\n", len(split[1].cards))
	}
}
