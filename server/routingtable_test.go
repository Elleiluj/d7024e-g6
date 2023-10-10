package server

import (
	"fmt"
	"testing"
)

func TestRoutingTable_print(t *testing.T) {
	fmt.Print("\n--------------------\n routingtable.go\n--------------------\n")
}

func TestRoutingTable(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002"))

	contacts := rt.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 20)

	// for i := range contacts {
	// 	fmt.Println(contacts[i].String())
	// }

	got := contacts[0].ID.String()
	want := NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002").ID.String()

	if got != want {
		t.Errorf("got %s want %s", got, want)
	} else {
		fmt.Println("NewRoutingTable \tPASS")
		fmt.Println("AddContact \t\tPASS")
		fmt.Println("FindClosestContacts \tPASS")
	}

}

func TestGetBucketIndex(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002"))

	id := NewKademliaID("2111111400000000000000000000000000000000")

	got := rt.getBucketIndex(id)
	want := 0

	if got != want {
		t.Errorf("got %d want %d", got, want)
	} else {
		fmt.Println("getBucketIndex \t\tPASS")
	}

}
