package server

import (
	"fmt"
	"testing"
)

func TestKademlia_print(t *testing.T) {
	fmt.Print("\n--------------------\n kademlia.go\n--------------------\n")
}

// NewKademliaNode
func TestNewKademiaNode(t *testing.T) {
	adress := "localhost:8000"
	id := NewKademliaNode(adress)

	got := id.Me.Address
	want := adress

	if got != want {
		t.Errorf("got %s want %s", got, want)
	} else {
		fmt.Println("NewKademliaNode \tPASS")
	}

}

// LookupContact
func TestLookupContact(t *testing.T) {
	// adress := "localhost:8000"
	// data := "00000000000000000000000000000000FFFFFFFF"
	// kademliaNode := NewKademliaNode(adress)

	// hashedData := CreateHash(string(data))
	// hashID := NewKademliaID(hashedData)

	// //Network := NewNetwork(&kademliaNode)

	// target := NewContact(hashID, kademliaNode.Me.Address)
	// //storingNodes := kademliaNode.LookupContact(&target)

	// kademliaNode.LookupContact(&target)

}

// sendAsyncFindContactMsg
func TestSendAsyncFindContactMsg(t *testing.T) {
	// adress := "localhost:8000"
	// data := "00000000000000000000000000000000FFFFFFFF"
	// kademliaNode := NewKademliaNode(adress)

	// hashedData := CreateHash(string(data))
	// hashID := NewKademliaID(hashedData)

	// //Network := NewNetwork(&kademliaNode)

	// target := NewContact(hashID, kademliaNode.Me.Address)
	// //storingNodes := kademliaNode.LookupContact(&target)

	// kademliaNode.LookupContact(&target)

}

// LookupData
// sendAsyncFindDataMsg
// Store
// JoinNetwork
// CreateHash
// Refresh
