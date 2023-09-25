package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const alpha int = 3

type Kademlia struct {
	Me           Contact
	RoutingTable *RoutingTable
	Data         map[string][]byte
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

	network := &Network{kademlia: kademlia}

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

		/* "The sequence of parallel searches is continued until either no node in the sets returned is
		closer than the closest node already seen or the initiating node has accumulated k probed and
		known to be active contacts." */
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
	hashID := NewKademliaID(hash)
	target := NewContact(hashID, kademlia.Me.Address)

	var value []byte
	var node *Contact

	// Initialize a shortlist with alpha closest nodes
	closestNodes := kademlia.RoutingTable.FindClosestContacts(target.ID, alpha)
	shortlist := NewShortList(closestNodes)

	// Keep track of the closest node seen so far
	closestNode := closestNodes[0]

	network := &Network{kademlia: kademlia}

	for {
		contactsChannel := make(chan []Contact)
		valueChannel := make(chan []byte)
		nodeChannel := make(chan *Contact)
		numAsked := 0

		for i := 0; i < shortlist.getLength() && numAsked < alpha; i++ {
			if !shortlist.nodes[i].isAsked {
				go kademlia.sendAsyncFindDataMsg(shortlist.nodes[i].contact, &target, hash, contactsChannel, valueChannel, nodeChannel, network)
				numAsked++
				shortlist.addContacts(<-contactsChannel)
				shortlist.dropUnactiveNodes()
				shortlist.sort()
			}

		}

		/* "The sequence of parallel searches is continued until either no node in the sets returned is
		closer than the closest node already seen or the initiating node has accumulated k probed and
		known to be active contacts." */

		value = <-valueChannel
		node = <-nodeChannel

		if *shortlist.nodes[0].contact == closestNode || shortlist.numOfAskedNodes() >= bucketSize || value != nil {
			println("BREAK!!" + kademlia.Me.Address)
			break
		}

		closestNode = *shortlist.nodes[0].contact

	}

	fmt.Printf("\nRetrieved value: %s, from node: %s\n", value, node.Address)

}

func (kademlia *Kademlia) sendAsyncFindDataMsg(contact *Contact, target *Contact, hash string, contactsChannel chan []Contact, valueChannel chan []byte, nodeChannel chan *Contact, network *Network) {
	contacts, value, node, err := network.SendFindDataMessage(&kademlia.Me, contact, target, hash)
	if err != nil {
		contactsChannel <- contacts
		valueChannel <- value
		nodeChannel <- node
	} else {
		contactsChannel <- contacts
		valueChannel <- value
		nodeChannel <- node
	}
}

func (kademlia *Kademlia) Store(data []byte) error {
	var err error
	network := &Network{kademlia: kademlia}
	hashedData := CreateHash(string(data))
	hashID := NewKademliaID(hashedData)

	target := NewContact(hashID, kademlia.Me.Address)
	storingNodes := kademlia.LookupContact(&target)
	for _, node := range storingNodes {
		err = network.SendStoreMessage(&kademlia.Me, &node, data, hashedData)
	}

	if err == nil {
		fmt.Printf("\nData stored with key: %s\n", hashedData)
	}

	return err

}

func (kademlia *Kademlia) JoinNetwork(knownNode *Contact) {
	fmt.Printf("Joining network through %s...\n", knownNode.String())
	kademlia.RoutingTable.AddContact(*knownNode)

	contacts := kademlia.LookupContact(&kademlia.Me)
	for _, contact := range contacts {
		kademlia.RoutingTable.AddContact(contact)
	}

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
