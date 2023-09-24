package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const alpha int = 3

type Kademlia struct {
	Me           Contact
	RoutingTable *RoutingTable
}

// address is full ip, including port
func NewKademliaNode(address string) (kademlia Kademlia) {
	kademliaID := NewKademliaID(CreateHash(address))
	fmt.Println("My kademlia ID: ", kademliaID)
	kademlia.Me = NewContact(kademliaID, address)
	kademlia.RoutingTable = NewRoutingTable(kademlia.Me)
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	// Initialize a shortlist with alpha closest nodes
	closestNodes := kademlia.RoutingTable.FindClosestContacts(target.ID, alpha)
	shortlist := NewShortList(closestNodes)

	// Keep track of the closest node seen so far
	closestNode := closestNodes[0]

	network := &Network{}

	for {
		responseChannel := make(chan []Contact)
		numAsked := 0

		for i := 0; i < shortlist.getLength() && numAsked < alpha; i++ {
			if !shortlist.nodes[i].isAsked {
				go kademlia.sendAsyncFindContactMsg(shortlist.nodes[i].contact, target, responseChannel, network)
				numAsked++
				shortlist.addContacts(<-responseChannel)
				shortlist.dropUnactiveNodes()
				shortlist.sort()
			}

		}

		if *shortlist.nodes[0].contact == closestNode || shortlist.numOfAskedNodes() >= bucketSize {
			println("BREAK!!" + kademlia.Me.Address)
			break
		}

		closestNode = *shortlist.nodes[0].contact

	}
	return shortlist.getContacts()
}

func (kademlia *Kademlia) sendAsyncFindContactMsg(contact *Contact, target *Contact, responseChannel chan []Contact, network *Network) {
	result, err := network.SendFindContactMessage(contact, target)
	if err != nil {
		responseChannel <- result
	} else {
		responseChannel <- result
	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

func (kademlia *Kademlia) JoinNetwork(knownNode *Contact) {
	fmt.Printf("Joining network through %s...\n", knownNode.String())
	kademlia.RoutingTable.AddContact(*knownNode)

	//network := &Network{}

	//
	//network.SendPingMessage(knownNode)
	kademlia.LookupContact(&kademlia.Me)

	// TODO: refresh k-buckets further away (lookup random node within the k-bucket range)

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
