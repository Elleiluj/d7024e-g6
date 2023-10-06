package server

import (
	"fmt"
	"testing"
)

func TestContac_printt(t *testing.T) {
	fmt.Print("\n--------------------\n Contact.go\n--------------------\n")
}

func TestString(t *testing.T) {
	contact := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")

	got := contact.String()
	want := "contact(\"ffffffff00000000000000000000000000000000\", \"localhost:8000\")"

	if got != want {
		t.Errorf("got %s want %s", got, want)
	} else {
		fmt.Print("NewContact \tPASS\n")
		fmt.Print("String \t\tPASS\n")
	}

}

func TestCalcDistance_con(t *testing.T) {
	contact := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	id := NewKademliaID("00000000FFFFFFFF000000000000000000000000")

	got := contact.ID.CalcDistance(id).String()
	want := NewKademliaID("FFFFFFFFFFFFFFFF000000000000000000000000").String()

	if got != want {
		t.Errorf("got %s want %s", got, want)
	} else {
		fmt.Print("CalcDistance \tPASS\n")
	}

}

func TestAppend(t *testing.T) {
	contact_1 := NewContact(NewKademliaID("00000000000000000000000000000000FFFFFFFF"), "localhost:8000")
	contact_2 := NewContact(NewKademliaID("000000000000000000000000FFFFFFFF00000000"), "localhost:8000")

	var candidates ContactCandidates
	contacts := []Contact{contact_1, contact_2}
	candidates.Append(contacts)

	got := candidates.contacts[0]
	want := contact_1

	if got != want {
		t.Errorf("got %s want %s", got, want)
	} else {
		fmt.Print("Append \t\tPASS\n")
	}

}

func TestLen_Contact(t *testing.T) {
	contact_1 := NewContact(NewKademliaID("00000000000000000000000000000000FFFFFFFF"), "localhost:8000")
	contact_2 := NewContact(NewKademliaID("000000000000000000000000FFFFFFFF00000000"), "localhost:8000")

	var candidates ContactCandidates
	contacts := []Contact{contact_1, contact_2}
	candidates.Append(contacts)

	got := candidates.Len()
	want := 2

	if got != want {
		t.Errorf("got %d want %d", got, want)
	} else {
		fmt.Print("Len \t\tPASS\n")
	}

}

func TestGetContacts(t *testing.T) {
	contact_1 := NewContact(NewKademliaID("00000000000000000000000000000000FFFFFFFF"), "localhost:8000")
	contact_2 := NewContact(NewKademliaID("000000000000000000000000FFFFFFFF00000000"), "localhost:8000")

	var candidates ContactCandidates
	contacts := []Contact{contact_1, contact_2}
	candidates.Append(contacts)

	got := candidates.GetContacts(2)[0].String()
	want := contacts[0].String()

	if got != want {
		t.Errorf("got %s want %s", got, want)
	} else {
		fmt.Print("GetContacts \tPASS\n")
	}

}

func TestSwap(t *testing.T) {
	contact_1 := NewContact(NewKademliaID("00000000000000000000000000000000FFFFFFFF"), "localhost:8000")
	contact_2 := NewContact(NewKademliaID("000000000000000000000000FFFFFFFF00000000"), "localhost:8000")

	var candidates ContactCandidates
	contacts := []Contact{contact_1, contact_2}
	candidates.Append(contacts)

	candidates.Swap(0, 1)

	got := candidates.GetContacts(2)[0].String()
	want := contact_2.String()

	if got != want {
		t.Errorf("got %s want %s", got, want)
	} else {
		fmt.Print("Swap \t\tPASS\n")
	}

}

// sorting test for contact ContactCandidates
// func TestSort_contact(t *testing.T) {
// 	contact_1 := NewContact(NewKademliaID("00000000000000000000000000000000FFFFFFFF"), "localhost:8000")
// 	contact_2 := NewContact(NewKademliaID("000000000000000000000000FFFFFFFF00000000"), "localhost:8000")
// 	contact_3 := NewContact(NewKademliaID("0000000000000000FFFFFFFF0000000000000000"), "localhost:8000")

// 	var candidates_sorted ContactCandidates
// 	contacts_sorted := []Contact{contact_1, contact_2, contact_3}
// 	candidates_sorted.Append(contacts_sorted)

// 	var candidates_unsorted ContactCandidates
// 	contacts_unsorted := []Contact{contact_1, contact_3, contact_2}
// 	candidates_unsorted.Append(contacts_unsorted)

// 	candidates_unsorted.Sort() // <- Creates an error

// 	got := candidates_unsorted.GetContacts(2)[0]
// 	want := candidates_sorted.GetContacts(2)[0]

// 	if got != want {
// 		t.Errorf("got %s want %s", got, want)
// 	} else {
// 		fmt.Print("Sort \t\tPASS\n")
// 	}
// }

// Less for contact
// func TestLess_con_T(t *testing.T) {
// 	contact_max := NewContact(NewKademliaID("000000000000000000000000FFFFFFFF00000000"), "localhost:8000")
// 	contact_min := NewContact(NewKademliaID("00000000000000000000000000000000FFFFFFFF"), "localhost:8000")

// 	got := contact_max.Less(&contact_min)
// 	want := true

// 	if got != want {
// 		t.Errorf("got %t want %t", got, want)
// 	} else {
// 		fmt.Print("PASS")
// 	}
// }

// Less for Candidates
// func TestLess_can_T(t *testing.T) {
// 	contact_1 := NewContact(NewKademliaID("00000000000000000000000000000000FFFFFFFF"), "localhost:8000")
// 	contact_2 := NewContact(NewKademliaID("000000000000000000000000FFFFFFFF00000000"), "localhost:8000")

// 	var candidates ContactCandidates
// 	contacts := []Contact{contact_1, contact_2}
// 	candidates.Append(contacts)

// 	got := candidates.Less(0, 1)
// 	want := true

// 	print("Hello?")

// 	if got != want {
// 		t.Errorf("got %t want %t", got, want)
// 	} else {
// 		fmt.Print("PASS")
// 	}
// }
