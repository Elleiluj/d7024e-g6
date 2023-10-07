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
func TestSort_contact(t *testing.T) {
	targetId := NewKademliaID("0000000000000000000000000000FFFFFFFFFFFF")
	contact_1 := NewContact(NewKademliaID("00000000000000000000000000000000FFFFFFFF"), "")
	contact_2 := NewContact(NewKademliaID("000000000000000000000000000000FFFFFFFFFF"), "")
	contact_3 := NewContact(NewKademliaID("0000000000000000000000000000000FFFFFFFFF"), "")

	contact_1.CalcDistance(targetId) // TODO: should probably mock
	contact_2.CalcDistance(targetId) // TODO: should probably mock
	contact_3.CalcDistance(targetId) // TODO: should probably mock

	var candidates ContactCandidates
	candidates.Append([]Contact{contact_1, contact_2, contact_3})

	candidates.Sort()

	got := candidates.contacts[1]
	want := contact_3

	if got != want {
		t.Errorf("got %s want %s", got, want)
	} else {
		fmt.Print("PASS")
	}
}

// Less for contact
func TestLess_con_T(t *testing.T) {
	targetId := NewKademliaID("0000000000000000000000000000FFFFFFFFFFFF")
	contact1 := NewContact(NewKademliaID("00000000000000000000000000000000FFFFFFFF"), "")
	contact2 := NewContact(NewKademliaID("000000000000000000000000000000FFFFFFFFFF"), "")

	contact1.CalcDistance(targetId) // TODO: should probably mock
	contact2.CalcDistance(targetId) // TODO: should probably mock

	got := contact1.Less(&contact2)
	want := false

	if got != want {
		t.Errorf("got %t want %t", got, want)
	} else {
		fmt.Print("PASS")
	}
}

// Less for Candidates
func TestLess_can_T(t *testing.T) {
	targetId := NewKademliaID("0000000000000000000000000000FFFFFFFFFFFF")
	contact_1 := NewContact(NewKademliaID("00000000000000000000000000000000FFFFFFFF"), "")
	contact_2 := NewContact(NewKademliaID("000000000000000000000000000000FFFFFFFFFF"), "")

	contact_1.CalcDistance(targetId) // TODO: should probably mock
	contact_2.CalcDistance(targetId) // TODO: should probably mock

	var candidates ContactCandidates
	candidates.Append([]Contact{contact_1, contact_2})

	got := candidates.Less(1, 0)
	want := true

	if got != want {
		t.Errorf("got %t want %t", got, want)
	} else {
		fmt.Print("PASS")
	}
}
