package server

import (
	"sort"
)

type ShortListNode struct {
	contact *Contact
	//isActive chan bool
	isAsked bool
}

type ShortList struct {
	nodes []ShortListNode
}

func NewShortList(contacts []Contact) ShortList {
	var shortlist ShortList
	for _, contact := range contacts {
		shortlist.nodes = append(shortlist.nodes, ShortListNode{
			contact: &contact,
			//isActive: make(chan bool),
			isAsked: false,
		})
	}
	return shortlist
}

func (shortList *ShortList) addContacts(contacts []Contact) {
	if contacts != nil {
		for _, contact := range contacts {
			// Create a copy of the contact inside the loop
			copyOfContact := contact

			// Check if the copy of the contact is already in the ShortList
			if !shortList.isInShortList(&copyOfContact) {
				shortList.nodes = append(shortList.nodes, ShortListNode{
					contact: &copyOfContact, // Use the address of the copy
					isAsked: false,
				})
			}
		}
	}
}

/*func (shortList *ShortList) dropUnactiveNodes() {
	var updatedNodes []ShortListNode

	for _, n := range shortList.nodes {
		if <-n.isActive || !n.isAsked {
			// Keep active nodes, append them to the updatedNodes slice
			updatedNodes = append(updatedNodes, n)
		}
	}

	shortList.nodes = updatedNodes

	fmt.Println("\n\nLENGTH AFTER DROPPING: " + strconv.Itoa(shortList.getLength()))
}*/

func (shortList *ShortList) addContact(contact Contact) {
	shortList.nodes = append(shortList.nodes, ShortListNode{
		contact: &contact,
		isAsked: false,
	})
}

func (shortList *ShortList) getAlphaNodes(alpha int) []ShortListNode {

	var alphaNodes []ShortListNode

	for _, n := range shortList.nodes {
		if !n.isAsked {
			// Keep active nodes, append them to the updatedNodes slice
			alphaNodes = append(alphaNodes, n)
			alpha--
		}
		if alpha == 0 {
			break
		}
	}

	return alphaNodes
}

/*func (shortList *ShortList) sort(target *Contact) {
	// Define a custom sorting function based on the distance
	sort.Slice(shortList.nodes, func(i, j int) bool {
		// Compare the distances of two contacts
		return shortList.nodes[i].contact.Distance.Less(shortList.nodes[j].contact.Distance)
	})
}*/

func (shortList *ShortList) sort(target *Contact) {
	// Define a custom sorting function based on the distance
	sort.Slice(shortList.nodes, func(i, j int) bool {
		// Compare the distances of two contacts to the target
		contact1 := shortList.nodes[i].contact
		contact1.CalcDistance(target.ID)
		distance1 := contact1.Distance
		contact2 := shortList.nodes[j].contact
		contact2.CalcDistance(target.ID)
		distance2 := contact2.Distance
		return distance1.Less(distance2)
	})
}

func (shortList *ShortList) getContacts() []Contact {
	var contacts []Contact

	for _, node := range shortList.nodes {
		contacts = append(contacts, *node.contact)
	}

	return contacts
}

func (shortList *ShortList) findUnqueriedNodes(k int) []ShortListNode {
	var unqueriedNodes []ShortListNode
	for i := 0; i < k && i < shortList.getLength(); i++ {
		node := shortList.nodes[i]
		if !node.isAsked {
			unqueriedNodes = append(unqueriedNodes, node)
		}
	}
	return unqueriedNodes
}

func (shortList *ShortList) numOfAskedNodes() int {
	askedNodes := 0
	for _, node := range shortList.nodes {
		if node.isAsked {
			askedNodes++
		}
	}
	return askedNodes
}

func (shortList *ShortList) getLength() int {
	return len(shortList.nodes)
}

func (shortList *ShortList) isInShortList(contact *Contact) bool {
	for _, node := range shortList.nodes {
		if node.contact != nil && node.contact == contact {
			return true
		}
	}
	return false
}
