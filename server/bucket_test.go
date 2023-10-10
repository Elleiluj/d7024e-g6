package server

import (
	"fmt"
	"testing"
)

func TestBucket_print(t *testing.T) {
	fmt.Print("\n--------------------\n bucket.go\n--------------------\n")
}

func TestAddContact(t *testing.T) {
	contact := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	bucket := newBucket()

	bucket.AddContact(contact)

	got := bucket.list.Front().Value
	want := contact

	if got != want {
		t.Errorf("got %s want %s", got, want)
	} else {
		fmt.Print("AddContact \tPASS\n")
	}
}

func TestLen_Bucket(t *testing.T) {
	contact := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	bucket := newBucket()

	bucket.AddContact(contact)

	got := bucket.Len()
	want := 1

	if got != want {
		t.Errorf("got %d want %d", got, want)
	} else {
		fmt.Print("Len \t\tPASS\n")
	}
}

func TestGetContactAndCalcDistance(t *testing.T) {
	id_1 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	id_2 := NewKademliaID("00000000FFFFFFFF000000000000000000000000")

	contact_1 := NewContact(id_1, "localhost:8000")
	contact_2 := NewContact(id_2, "localhost:8000")

	bucket := newBucket()

	bucket.AddContact(contact_1)
	bucket.AddContact(contact_2)

	got := bucket.GetContactAndCalcDistance(id_1)[0].Distance.String()
	want := NewKademliaID("ffffffffffffffff000000000000000000000000").String()

	if got != want {
		t.Errorf("got %s want %s", got, want)
	} else {
		fmt.Print("GetContactAndCalcDistance \tPASS\n")
	}

}
