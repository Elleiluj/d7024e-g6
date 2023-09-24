package main

import "sort"

type ShortListNode struct {
	contact  *Contact
	isActive bool
	isAsked  bool
}

type ShortList struct {
	nodes []ShortListNode
}

func NewShortList(contacts []Contact) ShortList {
	var shortlist ShortList
	for _, contact := range contacts {
		shortlist.nodes = append(shortlist.nodes, ShortListNode{
			contact:  &contact,
			isActive: false,
			isAsked:  false,
		})
	}
	return shortlist
}

func (shortList *ShortList) addContacts(contacts []Contact) {
	if contacts != nil {
		for _, contact := range contacts {
			if !shortList.isInShortList(&contact) {
				shortList.nodes = append(shortList.nodes, ShortListNode{
					contact:  &contact,
					isActive: false,
					isAsked:  false,
				})
			}
		}
	}
}

func (shortList *ShortList) dropUnactiveNodes() {
	var updatedNodes []ShortListNode

	for _, n := range shortList.nodes {
		if n.isActive || !n.isAsked {
			// Keep active nodes, append them to the updatedNodes slice
			updatedNodes = append(updatedNodes, n)
		}
	}

	shortList.nodes = updatedNodes
}

func (shortList *ShortList) addContact(contact Contact) {
	shortList.nodes = append(shortList.nodes, ShortListNode{
		contact:  &contact,
		isActive: false,
		isAsked:  false,
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

func (shortList *ShortList) sort() {
	// Define a custom sorting function based on the distance
	sort.Slice(shortList.nodes, func(i, j int) bool {
		// Compare the distances of two contacts
		return shortList.nodes[i].contact.distance.Less(shortList.nodes[j].contact.distance)
	})
}

func (shortList *ShortList) getContacts() []Contact {
	var contacts []Contact

	for _, node := range shortList.nodes {
		contacts = append(contacts, *node.contact)
	}

	return contacts
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
