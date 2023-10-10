package client

import (
	"fmt"
	"kademlia/server"
	"testing"
	"time"
)

func TestClient_print(t *testing.T) {
	fmt.Print("\n--------------------\n client.go\n--------------------\n")
}

func TestNewClient(t *testing.T) {
	kademlia := server.NewKademliaNode("127.0.0.1:12345")

	client := NewClient(&kademlia)
	if client.kademlia != &kademlia {
		t.Errorf("Expected client.kademlia to be %v, but got %v", kademlia, client.kademlia)
	} else {
		fmt.Println("NewClient \tPASS")
	}
}

func TestHandleInput(t *testing.T) {
	addr := "127.0.0.1:4021"
	kademlia := server.NewKademliaNode(addr)
	network := server.NewNetwork(&kademlia)
	kademlia.JoinNetwork(&kademlia.Me)
	go network.Listen(addr)

	client := NewClient(&kademlia)
	client.sleepTime = 1 * time.Second

	go client.Start()

	fail := false

	err := client.HandleInput([]string{"put", "test"})
	if err != nil {
		t.Errorf("Expected no error, but had")
		fail = true
	}
	err = client.HandleInput([]string{"get", server.CreateHash("test")})
	if err != nil {
		t.Errorf("Expected no error, but had")
		fail = true
	}
	err = client.HandleInput([]string{"get", "aaaaaaaaaaaa aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"})
	if err == nil {
		t.Errorf("Expected error, but had no error")
		fail = true
	}
	err = client.HandleInput([]string{"get", "invalid"})
	if err == nil {
		t.Errorf("Expected error, but had no error")
		fail = true
	}
	err = client.HandleInput([]string{"forget", server.CreateHash("test")})
	if err != nil {
		t.Errorf("Expected no error, but had")
		fail = true
	}
	err = client.HandleInput([]string{"forget", "aaaaaaaaaaaa aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"})
	if err == nil {
		t.Errorf("Expected error, but had no error")
		fail = true
	}
	err = client.HandleInput([]string{"forget", "invalid"})
	if err == nil {
		t.Errorf("Expected error, but had no error")
		fail = true
	}
	err = client.HandleInput([]string{"invalid"})
	if err == nil {
		t.Errorf("Expected error, but had no error")
		fail = true
	}
	if !fail {
		fmt.Println("HandleInput \tPASS")
	}
}
