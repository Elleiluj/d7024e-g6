package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

const alpha int = 3

type Kademlia struct {
	Me           Contact
	RoutingTable *RoutingTable
	Datastore    *Datastore
	UploadedData []string
}

// address is full ip, including port
func NewKademliaNode(address string) (kademlia Kademlia) {
	kademliaID := NewKademliaID(CreateHash(address))
	fmt.Println("My kademlia ID: ", kademliaID)
	kademlia.Me = NewContact(kademliaID, address)
	kademlia.RoutingTable = NewRoutingTable(kademlia.Me)
	kademlia.Datastore = NewDataStore()
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	closestNodes := kademlia.RoutingTable.FindClosestContacts(target.ID, alpha)
	shortlist := NewShortList(closestNodes)

	closestNode := closestNodes[0]
	network := &Network{kademlia: kademlia}

	for {
		responseChannel := make(chan []Contact)
		var inactiveNode chan *Contact
		numAsked := 0

		for i := 0; i < shortlist.getLength() && numAsked < alpha; i++ {
			if !shortlist.nodes[i].isAsked {
				go kademlia.sendAsyncFindContactMsg(shortlist.nodes[i].contact, target, responseChannel, inactiveNode, network)
				numAsked++
				shortlist.addContacts(<-responseChannel)
				if inactiveNode != nil {
					shortlist.dropNode(shortlist.nodes[i].contact)
				}
				shortlist.nodes[i].isAsked = true
			}
		}
		if shortlist.nodes[0].contact.Address == closestNode.Address || shortlist.numOfAskedNodes() >= bucketSize {
			if shortlist.nodes[0].contact.Address == closestNode.Address {
				unqueriedNodes := shortlist.findUnqueriedNodes(bucketSize)
				for _, node := range unqueriedNodes {
					contact := node.contact
					go kademlia.sendAsyncFindContactMsg(contact, target, responseChannel, inactiveNode, network)
					shortlist.addContacts(<-responseChannel)
					if inactiveNode != nil {
						shortlist.dropNode(node.contact)
					}
					node.isAsked = true
				}
			}
			break
		}
		closestNode = *shortlist.nodes[0].contact
	}
	shortlist.sort(target)
	return shortlist.getContacts()
}

func (kademlia *Kademlia) sendAsyncFindContactMsg(contact *Contact, target *Contact, responseChannel chan []Contact, inactiveNode chan *Contact, network *Network) {
	result, err := network.SendFindContactMessage(&kademlia.Me, contact, target)
	if err != nil {
		responseChannel <- result
		inactiveNode <- contact
	} else {
		responseChannel <- result
	}
}

func (kademlia *Kademlia) LookupData(hash string) (*Contact, string) {
	hashID := NewKademliaID(hash)
	target := NewContact(hashID, kademlia.Me.Address)
	var value []byte
	var node *Contact
	var inactiveNode *Contact

	closestNodes := kademlia.RoutingTable.FindClosestContacts(target.ID, alpha)
	shortlist := NewShortList(closestNodes)
	closestNode := closestNodes[0]
	network := &Network{kademlia: kademlia}

	for {
		contactsChannel := make(chan []Contact, bucketSize)
		valueChannel := make(chan []byte, 1)
		nodeChannel := make(chan *Contact, 1)
		inactiveNodeChannel := make(chan *Contact, 1)
		numAsked := 0

		for i := 0; i < shortlist.getLength() && numAsked < alpha; i++ {
			if !shortlist.nodes[i].isAsked {
				go kademlia.sendAsyncFindDataMsg(shortlist.nodes[i].contact, &target, hash, contactsChannel, valueChannel, nodeChannel, inactiveNodeChannel, network)
				numAsked++
				shortlist.addContacts(<-contactsChannel)
				inactiveNode = <-inactiveNodeChannel
				if inactiveNode != nil {
					shortlist.dropNode(shortlist.nodes[i].contact)
				}
				shortlist.nodes[i].isAsked = true
				value = <-valueChannel
				node = <-nodeChannel
			}
		}
		if shortlist.nodes[0].contact.Address == closestNode.Address || shortlist.numOfAskedNodes() >= bucketSize || value != nil {
			if shortlist.nodes[0].contact.Address == closestNode.Address {
				unqueriedNodes := shortlist.findUnqueriedNodes(bucketSize)
				for _, unqueriedNode := range unqueriedNodes {
					contact := unqueriedNode.contact
					go kademlia.sendAsyncFindDataMsg(contact, &target, hash, contactsChannel, valueChannel, nodeChannel, inactiveNodeChannel, network)
					shortlist.addContacts(<-contactsChannel)
					inactiveNode = <-inactiveNodeChannel
					if inactiveNode != nil {
						shortlist.dropNode(inactiveNode)
					}
					unqueriedNode.isAsked = true
					value = <-valueChannel
					node = <-nodeChannel
				}
				shortlist.sort(&target)
			}

			if value != nil {
				kademlia.resetTTL(network, kademlia.LookupContact(&target), hash)
			}
			break
		}
		closestNode = *shortlist.nodes[0].contact
	}

	var str string
	if value == nil {
		str = "\nNo value found with hash: " + hash + "\n"
	} else {
		str = "\nRetrieved value: " + string(value) + ", from node: " + node.Address + "\n"
	}
	fmt.Println(str)
	//fmt.Printf("\nRetrieved value: %s, from node: %s\n", value, node.Address)
	return node, string(value)
}

func (kademlia *Kademlia) sendAsyncFindDataMsg(contact *Contact, target *Contact, hash string, contactsChannel chan []Contact, valueChannel chan []byte, nodeChannel chan *Contact, inactiveNode chan *Contact, network *Network) {
	contacts, value, node, err := network.SendFindDataMessage(&kademlia.Me, contact, target, hash)
	if err != nil {
		contactsChannel <- contacts
		inactiveNode <- contact
		valueChannel <- value
		nodeChannel <- node
	} else {
		contactsChannel <- contacts
		inactiveNode <- nil
		valueChannel <- value
		nodeChannel <- node
	}
}

func (kademlia *Kademlia) resetTTL(network *Network, contacts []Contact, key string) {
	count := 0
	for _, contact := range contacts {
		count++
		network.SendResetTTLMessage(&kademlia.Me, &contact, key)
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

	kademlia.UploadedData = append(kademlia.UploadedData, hashedData)

	if err == nil {
		fmt.Printf("\nData stored with key: %s\n", hashedData)
	}

	return err

}

func (kademlia *Kademlia) Forget(hash string) {
	indexToRemove := -1
	for index, key := range kademlia.UploadedData {
		if hash == key {
			indexToRemove = index
			break
		}
	}
	if indexToRemove >= 0 && indexToRemove < len(kademlia.UploadedData) {
		kademlia.UploadedData = append(kademlia.UploadedData[:indexToRemove], kademlia.UploadedData[indexToRemove+1:]...)
		fmt.Println("Forgot data with hash: " + hash)

	} else {
		fmt.Println("Hash: \"" + hash + "\" has not been uploaded by this node")
	}

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

func (kademlia *Kademlia) RemoveExpiredData() {
	for {
		kademlia.Datastore.removeExpired()
		time.Sleep(5 * time.Second)
	}
}

func (kademlia *Kademlia) RefreshUploadedData() {
	for {
		for _, key := range kademlia.UploadedData {
			network := &Network{kademlia: kademlia}
			hashID := NewKademliaID(key)
			target := NewContact(hashID, kademlia.Me.Address)
			contacts := kademlia.LookupContact(&target)
			kademlia.resetTTL(network, contacts, key)
		}

		time.Sleep(20 * time.Second)
	}
}

func (kademlia *Kademlia) refreshData(key string) {
	kademlia.Datastore.resetTTL(key)
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
