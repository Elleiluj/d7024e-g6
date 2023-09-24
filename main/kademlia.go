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
	result, err := network.SendFindContactMessage(&kademlia.Me, contact, target)
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
	contacts := kademlia.LookupContact(&kademlia.Me)
	for _, contact := range contacts {
		kademlia.RoutingTable.AddContact(contact)
	}

	// TODO: refresh k-buckets further away (lookup random node within the k-bucket range)
	kademlia.Refresh(contacts[0])
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

// Refresh looks up random nodes in the range of its incomplete buckets to fill them
func (kademlia *Kademlia) Refresh(closestContact Contact) {
	// it refreshes all buckets further away than its closest neighbor, which will be in the
	// occupied bucket with the lowest index.
	// the node selects a random number in that range and does a refresh

	closestBucketIndex := kademlia.RoutingTable.getBucketIndex(closestContact.ID)
	for i := closestBucketIndex + 1; i < IDLength*8; i++ {
		if kademlia.RoutingTable.buckets[i].list.Len() < bucketSize {
			randomNodeToRefresh := NewContact(RandomKademliaIDInBucket(kademlia.Me.ID, i), "")
			contacts := kademlia.LookupContact(&randomNodeToRefresh)
			for _, contact := range contacts {
				kademlia.RoutingTable.AddContact(contact)
			}
		}
	}

}
