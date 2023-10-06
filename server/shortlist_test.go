package server

import (
	"fmt"
	"testing"
)

func TestShortlist_print(t *testing.T) {
	fmt.Print("\n--------------------\n shortlist.go\n--------------------\n")
}

func TestGetContacts_shortlist(t *testing.T) {
	contact_1 := NewContact(NewKademliaID("00000000000000000000000000000000FFFFFFFF"), "localhost:8000")
	contact_2 := NewContact(NewKademliaID("000000000000000000000000FFFFFFFF00000000"), "localhost:8000")

	contacts_1 := []Contact{contact_1, contact_2}

	shortList_1 := NewShortList(contacts_1)

	// Why is the first the same as the second?
	//fmt.Print(shortList_1.getContacts(), "\n\n")

	got := shortList_1.getContacts()[1].ID
	want := contacts_1[1].ID

	if got != want {
		t.Errorf("got %s want %s", got, want)
	} else {
		fmt.Println("NewShortList \tPASS")
		fmt.Println("getContacts \tPASS")
	}

}

func TestAddContacts(t *testing.T) {
	contact_1 := NewContact(NewKademliaID("00000000000000000000000000000000FFFFFFFF"), "localhost:8000")
	contact_2 := NewContact(NewKademliaID("000000000000000000000000FFFFFFFF00000000"), "localhost:8000")
	contact_3 := NewContact(NewKademliaID("0000000000000000FFFFFFFF0000000000000000"), "localhost:8000")
	contact_4 := NewContact(NewKademliaID("00000000FFFFFFFF000000000000000000000000"), "localhost:8000")

	contacts_1 := []Contact{contact_1, contact_2}
	contacts_2 := []Contact{contact_3, contact_4}
	contacts_1n2 := []Contact{contact_1, contact_2, contact_3, contact_4}

	shortList_1 := NewShortList(contacts_1)
	shortList_1.addContacts(contacts_2)

	// Create for loop to check each
	got := shortList_1.getContacts()[1].ID
	want := contacts_1n2[1].ID

	if got != want {
		t.Errorf("got %s want %s", got, want)
	} else {
		fmt.Println("addContacts \tPASS")
	}

}

func TestAddContact_shortlist(t *testing.T) {
	contact_1 := NewContact(NewKademliaID("00000000000000000000000000000000FFFFFFFF"), "localhost:8000")
	contact_2 := NewContact(NewKademliaID("000000000000000000000000FFFFFFFF00000000"), "localhost:8000")
	contact_3 := NewContact(NewKademliaID("0000000000000000FFFFFFFF0000000000000000"), "localhost:8000")

	contacts_1 := []Contact{contact_1, contact_2}
	contacts_2 := []Contact{contact_1, contact_2, contact_3}

	shortList_1 := NewShortList(contacts_1)
	shortList_1.addContact(contact_3)

	got := shortList_1.getContacts()[2].ID
	want := contacts_2[2].ID

	if got != want {
		t.Errorf("got %s want %s", got, want)
	} else {
		fmt.Println("addContact \tPASS")
	}

}

func TestGetAlphaNodes(t *testing.T) {

	contact_1 := NewContact(NewKademliaID("00000000FFFFFFFF000000000000000000000000"), "localhost:8000")
	contact_2 := NewContact(NewKademliaID("0000000000000000FFFFFFFF0000000000000000"), "localhost:8000")
	contact_3 := NewContact(NewKademliaID("000000000000000000000000FFFFFFFF00000000"), "localhost:8000")
	contact_4 := NewContact(NewKademliaID("00000000000000000000000000000000FFFFFFFF"), "localhost:8000")

	contacts_1 := []Contact{contact_1}
	contacts_2 := []Contact{contact_2, contact_3, contact_4}

	shortList := NewShortList(contacts_1)
	shortList.addContacts(contacts_2)

	// Create for loop to check each
	got := shortList.getAlphaNodes(4)[0].contact.ID
	want := contact_1.ID

	if got != want {
		t.Errorf("got %s want %s", got, want)
	} else {
		fmt.Println("getAlphaNodes \tPASS")
	}

}

func TestSort(t *testing.T) {

	contact_1 := NewContact(NewKademliaID("00000000FFFFFFFF000000000000000000000000"), "localhost:8000")
	contact_2 := NewContact(NewKademliaID("0000000000000000FFFFFFFF0000000000000000"), "localhost:8000")
	contact_3 := NewContact(NewKademliaID("000000000000000000000000FFFFFFFF00000000"), "localhost:8000")
	contact_4 := NewContact(NewKademliaID("00000000000000000000000000000000FFFFFFFF"), "localhost:8000")

	// Unsorted Shortlist
	contacts_1 := []Contact{contact_3}
	contacts_2 := []Contact{contact_2, contact_1, contact_4}

	shortList := NewShortList(contacts_1)
	shortList.addContacts(contacts_2)

	//Sorted Shortlist
	contacts_3 := []Contact{contact_4}
	contacts_4 := []Contact{contact_3, contact_2, contact_1}

	shortList_sorted := NewShortList(contacts_3)
	shortList_sorted.addContacts(contacts_4)

	// Sorting unsorted shortlist
	shortList.sort(&contact_1)

	// Create for loop to check each
	got := shortList.getContacts()[0].Address
	want := shortList_sorted.getContacts()[0].Address

	if got != want {
		t.Errorf("got %s want %s", got, want)
	} else {
		fmt.Println("sort \t\tPASS")
	}

}

func TestIsInShortList(t *testing.T) {
	contact_1 := NewContact(NewKademliaID("00000000FFFFFFFF000000000000000000000000"), "localhost:8000")
	contact_2 := NewContact(NewKademliaID("0000000000000000FFFFFFFF0000000000000000"), "localhost:8000")
	contact_3 := NewContact(NewKademliaID("000000000000000000000000FFFFFFFF00000000"), "localhost:8000")
	contact_4 := NewContact(NewKademliaID("00000000000000000000000000000000FFFFFFFF"), "localhost:8000")

	contacts_1 := []Contact{contact_1}
	contacts_2 := []Contact{contact_2, contact_3, contact_4}

	shortList := NewShortList(contacts_1)
	shortList.addContacts(contacts_2)

	//-------- clearly exists
	// fmt.Print(shortList.nodes[0].contact.String() == contact_1.String(), "\n")
	// fmt.Print(shortList.nodes[0].contact.String(), "\n")
	// fmt.Print(contact_1.String(), "\n")

	// Create for loop to check each
	got := shortList.isInShortList(&contact_1) // <- gives false
	want := true

	if got != want {
		t.Errorf("got %t want %t", got, want)
	} else {
		fmt.Println("isInShortList \tPASS")
	}

}

func TestDropNode(t *testing.T) {

	contact_1 := NewContact(NewKademliaID("00000000FFFFFFFF000000000000000000000000"), "localhost:8000")
	contact_2 := NewContact(NewKademliaID("0000000000000000FFFFFFFF0000000000000000"), "localhost:8000")
	contact_3 := NewContact(NewKademliaID("000000000000000000000000FFFFFFFF00000000"), "localhost:8000")
	contact_4 := NewContact(NewKademliaID("00000000000000000000000000000000FFFFFFFF"), "localhost:8000")

	contacts_1 := []Contact{contact_1}
	contacts_2 := []Contact{contact_2, contact_3, contact_4}
	contacts_3 := []Contact{contact_3, contact_4}

	shortList := NewShortList(contacts_1)
	shortList.addContacts(contacts_2)

	shortList_1 := NewShortList(contacts_1)
	shortList_1.addContacts(contacts_3)

	// Create for loop to check each
	shortList.dropNode(&contact_2)
	got := shortList.getContacts()[2]
	want := shortList_1.getContacts()[2]

	if got != want {
		t.Errorf("got %s want %s", got, want)
	} else {
		fmt.Println("dropNode \tPASS")
	}

}

func TestGetLength(t *testing.T) {
	contact_1 := NewContact(NewKademliaID("00000000FFFFFFFF000000000000000000000000"), "localhost:8000")
	contact_2 := NewContact(NewKademliaID("0000000000000000FFFFFFFF0000000000000000"), "localhost:8000")
	contact_3 := NewContact(NewKademliaID("000000000000000000000000FFFFFFFF00000000"), "localhost:8000")
	contact_4 := NewContact(NewKademliaID("00000000000000000000000000000000FFFFFFFF"), "localhost:8000")

	contacts_1 := []Contact{contact_1}
	contacts_2 := []Contact{contact_2, contact_3, contact_4}

	shortList := NewShortList(contacts_1)
	shortList.addContacts(contacts_2)

	// Create another test where nodes are asked
	got := shortList.getLength()
	want := 4

	if got != want {
		t.Errorf("got %d want %d", got, want)
	} else {
		fmt.Println("getLength \tPASS")
	}

}

func TestFindUnqueriedNodes(t *testing.T) {
	contact_1 := NewContact(NewKademliaID("00000000FFFFFFFF000000000000000000000000"), "localhost:8000")
	contact_2 := NewContact(NewKademliaID("0000000000000000FFFFFFFF0000000000000000"), "localhost:8000")
	contact_3 := NewContact(NewKademliaID("000000000000000000000000FFFFFFFF00000000"), "localhost:8000")
	contact_4 := NewContact(NewKademliaID("00000000000000000000000000000000FFFFFFFF"), "localhost:8000")

	contacts_1 := []Contact{contact_1}
	contacts_2 := []Contact{contact_2, contact_3, contact_4}

	shortList := NewShortList(contacts_1)
	shortList.addContacts(contacts_2)

	//fmt.Print(shortList.findUnqueriedNodes(4)[2].contact)

	// Create for loop to check each
	got := *shortList.findUnqueriedNodes(4)[2].contact
	want := shortList.getContacts()[2]

	if got != want {
		t.Errorf("got %s want %s", got, want)
	} else {
		fmt.Println("findUnqueriedNodes \tPASS")
	}
}

func TestNumOfAskedNodes(t *testing.T) {
	contact_1 := NewContact(NewKademliaID("00000000FFFFFFFF000000000000000000000000"), "localhost:8000")
	contact_2 := NewContact(NewKademliaID("0000000000000000FFFFFFFF0000000000000000"), "localhost:8000")
	contact_3 := NewContact(NewKademliaID("000000000000000000000000FFFFFFFF00000000"), "localhost:8000")
	contact_4 := NewContact(NewKademliaID("00000000000000000000000000000000FFFFFFFF"), "localhost:8000")

	contacts_1 := []Contact{contact_1}
	contacts_2 := []Contact{contact_2, contact_3, contact_4}

	shortList := NewShortList(contacts_1)
	shortList.addContacts(contacts_2)

	//fmt.Print(shortList.findUnqueriedNodes(4)[2].contact)

	// Create another test where nodes are asked
	got := shortList.numOfAskedNodes()
	want := 0

	if got != want {
		t.Errorf("got %d want %d", got, want)
	} else {
		fmt.Println("numOfAskedNodes \tPASS")
	}
}
