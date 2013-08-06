package main

import (
	"testing"
)

func TestMain(t *testing.T) {
	main()
}

func TestValueOfEmptyHand(t *testing.T) {
	hand := &Hand{[]Card{}}
	if len(hand.Value()) != 1 {
		t.Fatal("expected 1 value")
	}
	if hand.Value()[0] != 0 {
		t.Fatal("expected value of 0")
	}
}

func TestValueOfAceHand(t *testing.T) {
	hand := &Hand{[]Card{
		Card{Ace, Spades},
	}}

	if len(hand.Value()) != 2 {
		t.Fatal("expected 2 values")
	}
	if hand.Value()[0] != 1 {
		t.Fatalf("expected 1; got %d\n", hand.Value()[0])
	}
	if hand.Value()[1] != 11 {
		t.Fatalf("expected 11; got %d\n", hand.Value()[1])
	}
}

func TestValueOfKingHand(t *testing.T) {
	hand := &Hand{[]Card{
		Card{King, Spades},
	}}
	if len(hand.Value()) != 1 {
		t.Fatal("expected 1 value")
	}
	if hand.Value()[0] != 10 {
		t.Fatalf("expected 10; got %d\n", hand.Value()[0])
	}
}

func TestValueOfKingKingHand(t *testing.T) {
	hand := &Hand{[]Card{
		Card{King, Spades},
		Card{King, Clubs},
	}}
	if len(hand.Value()) != 1 {
		t.Fatal("expected 1 value")
	}
	if hand.Value()[0] != 20 {
		t.Fatalf("expected 20; got %d\n", hand.Value()[0])
	}
}

func TestValueOfAceKingHand(t *testing.T) {
	hand := &Hand{[]Card{
		Card{Ace, Spades},
		Card{King, Clubs},
	}}
	if len(hand.Value()) != 2 {
		t.Fatal("expected 1 value")
	}
	if hand.Value()[0] != 11 {
		t.Fatalf("expected 11; got %d\n", hand.Value()[0])
	}
	if hand.Value()[1] != 21 {
		t.Fatalf("expected 21; got %d\n", hand.Value()[1])
	}
}

func TestValueOfAceAceHand(t *testing.T) {
	hand := &Hand{[]Card{
		Card{Ace, Spades},
		Card{Ace, Clubs},
	}}
	if len(hand.Value()) != 3 {
		t.Fatalf("expected 3 values; got %#v\n", hand.Value())
	}
	if hand.Value()[0] != 2 {
		t.Fatalf("expected 2; got %d\n", hand.Value()[0])
	}
	if hand.Value()[1] != 12 {
		t.Fatalf("expected 12; got %d\n", hand.Value()[1])
	}
	if hand.Value()[2] != 12 {
		t.Fatalf("expected 12; got %d\n", hand.Value()[2])
	}
}

func TestBlackJackOn15(t *testing.T) {
	hand := &Hand{[]Card{
		Card{King, Diamonds},
		Card{Five, Diamonds},
	}}

	if hand.BlackJack() {
		t.Fatal("should not be blackjack")
	}
}

func TestBlackJack(t *testing.T) {
	hand := &Hand{[]Card{
		Card{King, Diamonds},
		Card{Ace, Diamonds},
	}}

	if !hand.BlackJack() {
		t.Fatal("expected blackjack")
	}
}
