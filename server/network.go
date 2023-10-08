package server

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

const pingMessage = "ping"
const pingResponse = "ping response"
const findContactMessage = "find contact"
const findContactResponse = "find contact response"
const findDataMessage = "find data"
const findDataResponse = "find data response"
const storeMessage = "store"
const storeResponse = "store response"
const resetTTLMessage = "resetTTL"
const resetTTLeResponse = "resetTTL response"

const port = 4000

type Network struct {
	kademlia *Kademlia
}

func NewNetwork(kademlia *Kademlia) *Network {
	return &Network{kademlia: kademlia}
}

// full ip
func (network *Network) Listen(ip string) {
	// Resolve UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", ip)
	if err != nil {
		fmt.Println("Error resolving address:", err)
	}

	// Create a UDP connection
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Error creating UDP connection:", err)
	}
	defer conn.Close()

	fmt.Println("UDP server is listening on", ip)

	// Buffer to hold incoming data
	buffer := make([]byte, 6000)

	for {
		// Read data from UDP connection
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP connection:", err)
			continue
		}

		// Process received data
		data := buffer[:n]
		//fmt.Printf("Received message from %s: %s\n", addr, string(data))

		network.handleResponse(data, addr, conn)
	}
}

func (network *Network) handleResponse(data []byte, address *net.UDPAddr, conn *net.UDPConn) {
	var message Message
	if err := json.Unmarshal(data, &message); err != nil {
		fmt.Println("Error decoding JSON message:", err)
		return
	}
	sender := message.Sender
	if sender != nil {
		network.kademlia.RoutingTable.AddContact(*sender)
	}
	switch message.Type {
	case pingMessage:
		address.Port = port
		network.SendPingResponse(address, conn)
	case pingResponse:
		// do nothing
	case findContactMessage:
		network.SendFindContactResponse(message, address, conn)
	case findDataMessage:
		network.SendFindDataResponse(message, address, conn)
	case storeMessage:
		network.SendStoreResponse(message, address, conn)
	case resetTTLMessage:
		network.SendResetTTLResponse(message, address, conn)
	default:
		fmt.Println("Unknown message type:", message.Type)
	}
}

func (network *Network) SendPingMessage(sender *Contact, contact *Contact) {
	message := Message{
		Type:   pingMessage,
		Sender: sender,
	}
	// Serialize message to JSON
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Error encoding JSON message: %v\n", err)
		return
	}

	// Establish UDP connection
	conn, err := net.Dial("udp", contact.Address)
	if err != nil {
		fmt.Printf("Error connecting to contact: %v\n", err)
		return
	}
	defer conn.Close()

	// Send JSON message to the contact's address
	_, err = conn.Write(jsonMessage)
	if err != nil {
		fmt.Printf("Error writing to contact: %v\n", err)
	}

	fmt.Printf("Sent ping message to %s\n", contact.Address)
}

func (network *Network) SendFindContactMessage(sender *Contact, contact *Contact, target *Contact) ([]Contact, error) {

	message := Message{
		Type:          findContactMessage,
		Sender:        sender,
		TargetContact: target,
	}

	responseMessage, err := network.SendMessage(contact, message)

	return responseMessage.Contacts, err
}

func (network *Network) SendFindDataMessage(sender *Contact, receiver *Contact, target *Contact, hash string) ([]Contact, []byte, *Contact, error) {
	message := Message{
		Type:           findDataMessage,
		Sender:         sender,
		TargetContact:  target,
		DataHashString: hash,
	}

	responseMessage, err := network.SendMessage(receiver, message)

	return responseMessage.Contacts, responseMessage.HashedData, receiver, err
}

func (network *Network) SendStoreMessage(sender *Contact, receiver *Contact, data []byte, key string) error {
	message := Message{
		Type:           storeMessage,
		Sender:         sender,
		HashedData:     data,
		DataHashString: key,
	}

	_, err := network.SendMessage(receiver, message)

	return err
}

func (network *Network) SendResetTTLMessage(sender *Contact, receiver *Contact, key string) error {
	message := Message{
		Type:           resetTTLMessage,
		Sender:         sender,
		DataHashString: key,
	}

	_, err := network.SendMessage(receiver, message)

	return err
}

func (network *Network) SendMessage(receiver *Contact, message Message) (Message, error) {

	udpAddr, err := net.ResolveUDPAddr("udp", receiver.Address)
	if err != nil {
		fmt.Println("Error resolving address:", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("\n\n ERROR DIALING!!!!")
		return Message{}, err
	}
	defer conn.Close()

	// Serialize the message into JSON
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Println("\n\n ERROR MARSHAL!!!!")
		return Message{}, err
	}

	_, err = conn.Write(jsonMessage)
	if err != nil {
		fmt.Println("\n\n ERROR WRITING!!!!")
		return Message{}, err
	}

	//fmt.Printf("Sent message to %s,\n message: %s\n", receiver.Address, jsonMessage)

	timeout := time.Now().Add(5 * time.Second) // Set a 5-second timeout
	conn.SetDeadline(timeout)

	// Wait for a response
	responseBuffer := make([]byte, 6000)
	n, err := conn.Read(responseBuffer)
	if err != nil {
		fmt.Println("\n\n ERROR READING!!!!")
		return Message{}, err
	}

	// Deserialize the response into a Message
	var responseMessage Message
	if err := json.Unmarshal(responseBuffer[:n], &responseMessage); err != nil {
		fmt.Println("\n\n ERROR UNMARSHAL!!!!")
		return Message{}, err
	}

	//fmt.Printf("Received response: %s\n", string(responseBuffer))

	return responseMessage, nil

}

func (network *Network) SendPingResponse(address *net.UDPAddr, conn *net.UDPConn) {
	response := Message{
		Type: pingResponse,
	}

	// Serialize response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Error encoding JSON response: %v\n", err)
		return
	}

	// Send JSON response to the contact's address
	_, err = conn.WriteToUDP(jsonResponse, address)
	if err != nil {
		fmt.Println("Error sending UDP response:", err)
	}

	fmt.Printf("Sent ping response to %s\n", address.String())
}

func (network *Network) SendFindContactResponse(message Message, address *net.UDPAddr, conn *net.UDPConn) {

	closestContacts := network.kademlia.RoutingTable.FindClosestContacts(message.TargetContact.ID, bucketSize)
	response := Message{
		Type:     findContactResponse,
		Contacts: closestContacts,
	}

	network.sendResponse(response, address, conn)

}

func (network *Network) SendFindDataResponse(message Message, address *net.UDPAddr, conn *net.UDPConn) {
	value := network.kademlia.Datastore.getData(message.DataHashString)
	var response Message

	// if no value is found, return k-closest
	if value == nil {
		closestContacts := network.kademlia.RoutingTable.FindClosestContacts(message.TargetContact.ID, bucketSize)
		response = Message{
			Type:     findDataResponse,
			Contacts: closestContacts,
		}
	} else {
		response = Message{
			Type:       findDataResponse,
			HashedData: value,
		}
	}
	network.sendResponse(response, address, conn)
}

func (network *Network) SendStoreResponse(message Message, address *net.UDPAddr, conn *net.UDPConn) {
	network.kademlia.Datastore.addData(message.HashedData, message.DataHashString)
	response := Message{
		Type: storeResponse,
	}

	network.sendResponse(response, address, conn)

}

func (network *Network) SendResetTTLResponse(message Message, address *net.UDPAddr, conn *net.UDPConn) {
	network.kademlia.refreshData(message.DataHashString)
	fmt.Println("\n\n\nReset TTL with key: " + message.DataHashString + "\n\n")
	response := Message{
		Type: resetTTLeResponse,
	}

	network.sendResponse(response, address, conn)

}

func (network *Network) sendResponse(response Message, address *net.UDPAddr, conn *net.UDPConn) {
	// Serialize response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		fmt.Printf("Error encoding JSON response: %v\n", err)
		return
	}

	// Send JSON response to the contact's address
	_, err = conn.WriteToUDP(jsonResponse, address)
	if err != nil {
		fmt.Println("Error sending UDP response:", err)
	}

	//fmt.Printf("Sent response to %s,\n response: %s\n", address.String(), jsonResponse)

}

type Message struct {
	Type           string
	Sender         *Contact
	Contacts       []Contact
	TargetContact  *Contact
	HashedData     []byte // data to store (value)
	DataHashString string // data to retrieve (key)
}
