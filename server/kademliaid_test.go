package server

import (
	"fmt"
	"testing"
)

func TestKademliaid_print(t *testing.T) {
	fmt.Print("\n--------------------\n kademliaid.go\n--------------------\n")
}

func TestCalcDistance(t *testing.T) {
	id_1 := NewKademliaID("0000000000ffffffffffffffffffff0000000000")
	id_2 := NewKademliaID("ffffffffffffffffffff00000000000000000000")

	// Same
	got := id_1.CalcDistance(id_2)
	want := NewKademliaID("ffffffffff0000000000ffffffffff0000000000")

	if *got != *want {
		t.Errorf("got %s want %s", got, want)
	} else {
		fmt.Println("NewKademliaID \tPASS")
		fmt.Println("CalDistance \tPASS")
	}
}

func TestString_kadid(t *testing.T) {
	id_1 := NewKademliaID("0000000000ffffffffffffffffffff0000000000")

	// Same
	got := id_1.String()
	want := "0000000000ffffffffffffffffffff0000000000"

	if got != want {
		t.Errorf("got %s want %s", got, want)
	} else {
		fmt.Println("String \t\tPASS")
	}
}

// NewRandomKademliaID
// RandomKademliaIDInBucket

func TestLess_true(t *testing.T) {
	id_min := NewKademliaID("0000000000000000000000000000000000000000")
	id_max := NewKademliaID("ffffffffffffffffffffffffffffffffffffffff")

	got := id_min.Less(id_max)
	want := true

	if got != want {
		t.Errorf("got %t want %t", got, want)
	} else {
		fmt.Println("Less_T \t\tPASS")
	}
}

func TestLess_false(t *testing.T) {
	id_min := NewKademliaID("0000000000000000000000000000000000000000")
	id_max := NewKademliaID("ffffffffffffffffffffffffffffffffffffffff")

	got := id_max.Less(id_min)
	want := false

	if got != want {
		t.Errorf("got %t want %t", got, want)
	} else {
		fmt.Println("Less_F \t\tPASS")
	}
}

func TestEquals_true(t *testing.T) {
	id_1 := NewKademliaID("0000000000ffffffffffffffffffff0000000000")
	id_2 := NewKademliaID("0000000000ffffffffffffffffffff0000000000")

	got := id_1.Equals(id_2)
	want := true

	if got != want {
		t.Errorf("got %t want %t", got, want)
	} else {
		fmt.Println("Equals_T \tPASS")
	}
}

func TestEquals_false(t *testing.T) {
	id_1 := NewKademliaID("0000000000ffffffffffffffffffff0000000000")
	id_2 := NewKademliaID("ffffffffffffffffffff00000000000000000000")

	got := id_1.Equals(id_2)
	want := false

	if got != want {
		t.Errorf("got %t want %t", got, want)
	} else {
		fmt.Println("Equals_F \tPASS")
	}
}
