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

func TestNewRandomKademliaID(t *testing.T) {
	fail := false
	newKademliaID := NewRandomKademliaID()
	if len(newKademliaID) != IDLength {
		fail = true
		t.Errorf("Expected Kademlia ID length of %d, but got %d", IDLength, len(newKademliaID))
	}
	for i, b := range newKademliaID {
		if b > 255 {
			fail = true
			t.Errorf("Byte at index %d is out of valid range: %d", i, b)
		}
	}
	if !fail {
		fmt.Println("NewRandomKademliaID \tPASS")
	}

}

func TestRandomKademliaIDInBucket(t *testing.T) {
	currentId := NewRandomKademliaID()
	bucketIndex := 5 // example bucket index

	minID, maxID := BucketRange(bucketIndex, currentId)

	newKademliaID := RandomKademliaIDInBucket(currentId, bucketIndex)

	distance := currentId.CalcDistance(newKademliaID)

	minDis := currentId.CalcDistance(minID)
	maxDis := currentId.CalcDistance(maxID)

	// Check if the distance is within the valid range for the specified bucket
	if distance.Less(minDis) || maxDis.Less(distance) {
		t.Errorf("Generated Kademlia ID is not within the specified bucket range")
	} else {
		fmt.Println("RandomKademliaIDInBucket \tPASS")
	}
}

func BucketRange(bucketIndex int, currentID *KademliaID) (minID, maxID *KademliaID) {
	minID = NewKademliaID("0000000000000000000000000000000000000000")
	maxID = NewKademliaID("0000000000000000000000000000000000000000")

	wholeBytes := (IDLength - 1) - (bucketIndex / 8)
	leftOverBits := uint(bucketIndex % 8)

	// Set the minimum ID
	minID[wholeBytes] = 1 << leftOverBits

	// Set the maximum ID
	for i := wholeBytes; i < IDLength; i++ {
		maxID[i] = 0xFF
	}
	maxID[wholeBytes] ^= (1 << leftOverBits) - 1

	// XOR the IDs with the current ID
	for i := 0; i < IDLength; i++ {
		minID[i] ^= currentID[i]
		maxID[i] ^= currentID[i]
	}

	return minID, maxID
}
