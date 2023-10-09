package server

import (
	"fmt"
	"testing"
	"time"
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

func TestLookupContact(t *testing.T) {
	addr := "127.0.0.1:4001"
	kademlia := NewKademliaNode(addr)
	network := &Network{kademlia: &kademlia}
	go network.Listen(addr)

	addr2 := "127.0.0.1:4002"
	kademlia2 := NewKademliaNode(addr2)
	network2 := &Network{kademlia: &kademlia2}
	go network2.Listen(addr2)

	addr3 := "127.0.0.1:4003"
	kademlia3 := NewKademliaNode(addr3)
	network3 := &Network{kademlia: &kademlia3}
	go network3.Listen(addr3)

	network.kademlia.JoinNetwork(&kademlia.Me)
	network2.kademlia.JoinNetwork(&kademlia.Me)
	network3.kademlia.JoinNetwork(&kademlia.Me)

	targetContact := &Contact{
		ID:      NewKademliaID(CreateHash("127.0.0.1:4003")),
		Address: "127.0.0.1:4003",
	}

	contacts := network.kademlia.LookupContact(targetContact)

	expectedResponse1 := []Contact{
		{ID: NewKademliaID(CreateHash("127.0.0.1:4001")), Address: "127.0.0.1:4001"}, {ID: NewKademliaID(CreateHash("127.0.0.1:4002")), Address: "127.0.0.1:4002"},
	}

	expectedCondition := ((contacts[0].Address == addr) || (contacts[1].Address == addr) || (contacts[2].Address == addr)) &&
		((contacts[0].Address == addr2) || (contacts[1].Address == addr2) || (contacts[2].Address == addr2)) &&
		((contacts[0].Address == addr3) || (contacts[1].Address == addr3) || (contacts[2].Address == addr3))

	if !expectedCondition {
		t.Errorf("Received response %v, expected %v", contacts, expectedResponse1)
	} else {
		fmt.Println("LookupContact \tPASS")
	}
}

func TestLookupData(t *testing.T) {
	addr := "127.0.0.1:4004"
	kademlia := NewKademliaNode(addr)
	network := &Network{kademlia: &kademlia}
	go network.Listen(addr)

	addr2 := "127.0.0.1:4005"
	kademlia2 := NewKademliaNode(addr2)
	network2 := &Network{kademlia: &kademlia2}
	go network2.Listen(addr2)

	addr3 := "127.0.0.1:4006"
	kademlia3 := NewKademliaNode(addr3)
	network3 := &Network{kademlia: &kademlia3}
	go network3.Listen(addr3)

	network.kademlia.JoinNetwork(&kademlia.Me)
	network2.kademlia.JoinNetwork(&kademlia.Me)
	network3.kademlia.JoinNetwork(&kademlia.Me)

	fail := false

	byteData := []byte("test")

	err := network.kademlia.Store(byteData)

	if err != nil {
		t.Errorf("Error storing")
		fail = true
	}

	_, value := network2.kademlia.LookupData(CreateHash("test"))

	expected := "test"

	if expected != value {
		t.Errorf("Received response %v, expected %v", value, expected)
		fail = true
	}

	_, value2 := network2.kademlia.LookupData(CreateHash("invalid"))

	expected2 := ""

	if expected2 != value2 {
		t.Errorf("Received response %v, expected %v", value2, expected2)
		fail = true
	}

	if !fail {
		fmt.Println("LookupData \tPASS")
	}
}

func TestForget(t *testing.T) {
	addr := "127.0.0.1:4004"
	kademlia := NewKademliaNode(addr)
	key1 := CreateHash("test1")
	key2 := CreateHash("test2")
	key3 := CreateHash("test3")
	kademlia.UploadedData = append(kademlia.UploadedData, key1, key2, key3)

	fail := false

	exists := false
	kademlia.Forget(key1)
	for _, value := range kademlia.UploadedData {
		if value == key1 {
			exists = true
		}
	}
	if exists {
		t.Errorf("Key not removed from uploadedData %v, expected %v", false, exists)
		fail = true
	}
	kademlia.Forget("invalid")
	if len(kademlia.UploadedData) != 2 {
		t.Errorf("Removed invalid value, length is %v, expected %v", len(kademlia.Datastore.Data), 2)
		fail = true
	}
	if !fail {
		fmt.Println("Forget \tPASS")
	}

}

func TestRemoveExpiredKademlia(t *testing.T) {
	addr := "127.0.0.1:4010"
	kademlia := NewKademliaNode(addr)
	data := []byte("test")
	key := CreateHash("test")
	kademlia.Datastore.addDataWithTTL(data, key, 3*time.Second)
	fail := false
	go kademlia.RemoveExpiredData()
	if len(kademlia.Datastore.Data) == 0 {
		t.Errorf("Data removed to early, len was %v, expected %v", len(kademlia.Datastore.Data), 1)
		fail = true
	}

	time.Sleep(5 * time.Second)
	if len(kademlia.Datastore.Data) != 0 {
		t.Errorf("Data not removed, len was %v, expected %v", len(kademlia.Datastore.Data), 0)
		fail = true
	}

	if !fail {
		fmt.Println("RemoveExpired \tPASS")
	}
}

func TestRefreshUploadedData(t *testing.T) {
	addr := "127.0.0.1:4011"
	kademlia := NewKademliaNode(addr)
	network := &Network{kademlia: &kademlia}
	go network.Listen(addr)

	addr2 := "127.0.0.1:4012"
	kademlia2 := NewKademliaNode(addr2)
	network2 := &Network{kademlia: &kademlia2}
	go network2.Listen(addr2)

	network.kademlia.JoinNetwork(&kademlia.Me)
	network2.kademlia.JoinNetwork(&kademlia.Me)

	fail := false

	byteData := []byte("test")

	err := network.kademlia.Store(byteData)

	if err != nil {
		t.Errorf("Error storing")
		fail = true
	}

	go network.kademlia.RefreshUploadedData()

	time.Sleep(TTLconst)
	time.Sleep(5 * time.Second)

	_, value := network2.kademlia.LookupData(CreateHash("test"))

	expected := "test"

	if expected != value {
		t.Errorf("Received response %v, expected %v", value, expected)
		fail = true
	}

	if !fail {
		fmt.Println("RefreshUploadedData \tPASS")
	}
}
