package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const alpha int = 4

type Kademlia struct {
	me           Contact
	routingTable RoutingTable
}

// address is full ip, including port
func NewKademliaNode(address string) (kademlia Kademlia) {
	kademliaID := NewKademliaID(CreateHash(address))
	fmt.Println("My kademlia ID: ", kademliaID)
	kademlia.me = NewContact(kademliaID, address)
	kademlia.routingTable = *NewRoutingTable(kademlia.me)
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	closestNodes := kademlia.routingTable.FindClosestContacts(target.ID, alpha) // Find K closest nodes
	var responses []Contact
	network := &Network{}

	for range closestNodes {
		response := network.SendFindContactMessage(target)
		if response != nil {
			responses = append(responses, response...)
		}
	}

	fmt.Printf("Lookup contact of: %s, found: %s.", target.ID, responses)

	return responses
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

func (kademlia *Kademlia) JoinNetwork(knownNode *Contact) {
	fmt.Printf("Joining network through %s...", knownNode.String())
	kademlia.routingTable.AddContact(*knownNode)

	kademlia.LookupContact(knownNode)

	// TODO: refresh k-buckets further away (lookup random node withing the k-bucket range)

}

func CreateHash(data string) (hash string) {
	// Create a new SHA-256 hash object
	hasher := sha256.New()

	// Write the data to the hash object
	hasher.Write([]byte(data))

	// Calculate the hash and store it as a byte slice
	hashBytes := hasher.Sum(nil)

	// Convert the byte slice to a hexadecimal string
	hash = hex.EncodeToString(hashBytes)

	return hash
}
