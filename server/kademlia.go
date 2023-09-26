package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
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
	kademlia.Data = make(map[string][]byte)
	return kademlia
}

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	// Initialize a shortlist with alpha closest nodes
	closestNodes := kademlia.RoutingTable.FindClosestContacts(target.ID, alpha)
	shortlist := NewShortList(closestNodes)

	// Keep track of the closest node seen so far
	closestNode := closestNodes[0]

	network := &Network{kademlia: kademlia}

	//var wg sync.WaitGroup

	for {

		responseChannel := make(chan []Contact)
		//isActive := make(chan bool)
		numAsked := 0

		for i := 0; i < shortlist.getLength() && numAsked < alpha; i++ {
			if !shortlist.nodes[i].isAsked {
				//wg.Add(1)

				/*go func(node ShortListNode) {
					defer wg.Done() // Decrement the WaitGroup when the query is done
					kademlia.sendAsyncFindContactMsg(node.contact, target, responseChannel, isActive, network)
				}(shortlist.nodes[i])*/
				go kademlia.sendAsyncFindContactMsg(shortlist.nodes[i].contact, target, responseChannel, network)
				numAsked++
				shortlist.addContacts(<-responseChannel)
				fmt.Println(shortlist.getLength())
				//shortlist.nodes[i].isActive = isActive
				shortlist.nodes[i].isAsked = true

			}

		}

		//wg.Wait()

		//shortlist.dropUnactiveNodes()

		fmt.Println("\n\nIs asked:")
		fmt.Println(shortlist.numOfAskedNodes())
		fmt.Println("\n\nShortlist length:")
		fmt.Println(shortlist.getLength())

		//shortlist.sort(target)

		/* "The sequence of parallel searches is continued until either no node in the sets returned is
		closer than the closest node already seen or the initiating node has accumulated k probed and
		known to be active contacts." */
		if *shortlist.nodes[0].contact == closestNode || shortlist.numOfAskedNodes() >= bucketSize {

			fmt.Println("\n\nENTERED BREAK CONDITION CLOSEST NODE!!")

			if *shortlist.nodes[0].contact == closestNode {

				fmt.Println("\n\nUNCHANGED CLOSEST NODE!!")

				// Find the k closest nodes that haven't been queried yet
				unqueriedNodes := shortlist.findUnqueriedNodes(bucketSize)

				fmt.Println("\n\nUNQUERIED NODES: " + strconv.Itoa(len(unqueriedNodes)))

				// Send FIND_* RPCs to unqueried nodes
				for _, node := range unqueriedNodes {
					contact := node.contact

					fmt.Println("\n\nFOR LOOP!")

					// Increment the WaitGroup to track the ongoing query
					//wg.Add(1)

					/*go func(contact Contact) {
						defer wg.Done() // Decrement the WaitGroup when the query is done
						kademlia.sendAsyncFindContactMsg(contact, target, responseChannel, network)
					}(node)*/

					go kademlia.sendAsyncFindContactMsg(contact, target, responseChannel, network)

					shortlist.addContacts(<-responseChannel)
					//node.isActive = isActive
					node.isAsked = true
				}

				fmt.Println("\n\nFOR LOOP FINISHED!")

				//shortlist.dropUnactiveNodes()
			}

			fmt.Println("BREAK!!" + kademlia.Me.Address)
			fmt.Println("\n\nShortlist closest: " + shortlist.nodes[0].contact.Address + "\nClosestNode: " + closestNode.Address)
			fmt.Println("\nNumOfAskedNodes: " + strconv.Itoa(shortlist.numOfAskedNodes()) + "\nBucketSize: " + strconv.Itoa(bucketSize) + "\n\n")
			break
		}

		closestNode = *shortlist.nodes[0].contact

	}
	shortlist.sort(target)
	return shortlist.getContacts()
}

func (kademlia *Kademlia) sendAsyncFindContactMsg(contact *Contact, target *Contact, responseChannel chan []Contact, network *Network) {
	result, err := network.SendFindContactMessage(&kademlia.Me, contact, target)
	fmt.Println("GOT MESSAGE RESPONSE; FUNCTION RETURNED")
	if err != nil {
		//isActive <- false
		responseChannel <- result
	} else {
		//isActive <- true
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
		//isActive := make(chan bool)
		numAsked := 0

		for i := 0; i < shortlist.getLength() && numAsked < alpha; i++ {
			if !shortlist.nodes[i].isAsked {
				go kademlia.sendAsyncFindDataMsg(shortlist.nodes[i].contact, &target, hash, contactsChannel, valueChannel, nodeChannel, network)
				numAsked++
				shortlist.addContacts(<-contactsChannel)
			}

		}

		//shortlist.dropUnactiveNodes()
		shortlist.sort(&target)

		/* "The sequence of parallel searches is continued until either no node in the sets returned is
		closer than the closest node already seen or the initiating node has accumulated k probed and
		known to be active contacts." */

		value = <-valueChannel
		node = <-nodeChannel

		fmt.Println("\n\nValue: " + string(value) + "\n\n")
		fmt.Println("\n\nNode: " + string(node.Address) + "\n\n")

		if *shortlist.nodes[0].contact == closestNode || shortlist.numOfAskedNodes() >= bucketSize || value != nil {

			// TODO:
			//If a cycle doesn't find a closer node, if closestNode is unchanged,
			// then the initiating node sends a FIND_* RPC to each of the k closest nodes that it has not already queried.

			if *shortlist.nodes[0].contact == closestNode {

				fmt.Println("\n\nUNCHANGED CLOSEST NODE!!")

				// Find the k closest nodes that haven't been queried yet
				unqueriedNodes := shortlist.findUnqueriedNodes(bucketSize)

				fmt.Println("\n\nUNQUERIED NODES: " + strconv.Itoa(len(unqueriedNodes)))

				// Send FIND_* RPCs to unqueried nodes
				for _, node := range unqueriedNodes {
					fmt.Println("\n\nFOR LOOP!")
					contact := node.contact
					go kademlia.sendAsyncFindDataMsg(contact, &target, hash, contactsChannel, valueChannel, nodeChannel, network)
					shortlist.addContacts(<-contactsChannel)
					//node.isActive = isActive
					node.isAsked = true
				}

				fmt.Println("\n\nFOR LOOP FINISHED!")

				//shortlist.dropUnactiveNodes()
				shortlist.sort(&target)
			}

			value = <-valueChannel
			node = <-nodeChannel

			fmt.Println("BREAK!!" + kademlia.Me.Address)
			fmt.Println("\n\nShortlist closest: " + shortlist.nodes[0].contact.Address + "\nClosestNode: " + closestNode.Address)
			fmt.Println("\nNumOfAskedNodes: " + strconv.Itoa(shortlist.numOfAskedNodes()) + "\nBucketSize: " + strconv.Itoa(bucketSize))
			fmt.Println("\nValue: " + string(value) + "\n\n")
			break
		}

		closestNode = *shortlist.nodes[0].contact

	}

	fmt.Printf("\nRetrieved value: %s, from node: %s\n", value, node.Address)

}

func (kademlia *Kademlia) sendAsyncFindDataMsg(contact *Contact, target *Contact, hash string, contactsChannel chan []Contact, valueChannel chan []byte, nodeChannel chan *Contact, network *Network) {
	contacts, value, node, err := network.SendFindDataMessage(&kademlia.Me, contact, target, hash)
	if err != nil {
		//isActive <- false
		contactsChannel <- contacts
		valueChannel <- value
		nodeChannel <- node
	} else {
		//isActive <- true
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
