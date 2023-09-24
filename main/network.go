package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
)

const pingMessage = "ping"
const pingResponse = "ping response"
const findContactMessage = "find contact"
const findContactResponse = "find contact response"
const findDataMessage = "find data"
const findDataResponse = "find data response"
const storeMessage = "store"
const storeResponse = "store response"
const messageDivider = ";"

const port = 4000

type Network struct {
	kademlia *Kademlia
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
		fmt.Printf("Received message from %s: %s\n", addr, string(data))

		// Respond to the client
		// maybe do in separate function, to simplify testing
		// depending on what the request was, we want different responses
		/*response := []byte("Hello from UDP server")
		addr.Port = port
		_, err = conn.WriteToUDP(response, addr)
		if err != nil {
			fmt.Println("Error sending UDP response:", err)
		}*/
		network.handleResponse(data, addr, conn)
	}
}

func (network *Network) handleResponse(data []byte, address *net.UDPAddr, conn *net.UDPConn) {
	var message Message
	if err := json.Unmarshal(data, &message); err != nil {
		fmt.Println("Error decoding JSON message:", err)
		return
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
		fmt.Println("Find data message!")
	case findDataResponse:
		fmt.Println("Find data response!")
	case storeMessage:
		fmt.Println("Store message!")
	case storeResponse:
		fmt.Println("Store response!")
	default:
		fmt.Println("Unknown message type:", message.Type)
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	message := Message{
		Type: pingMessage,
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

func (network *Network) SendFindContactMessage(contact *Contact, target *Contact) ([]Contact, error) {

	udpAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	if err != nil {
		fmt.Println("Error resolving address:", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	message := Message{
		Type:          findContactMessage,
		TargetContact: target,
	}

	fmt.Printf("TEST send!! messagetype: %s, target: %s", message.Type, message.TargetContact.ID)

	// Serialize the message into JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	fmt.Printf("TEST send json!! marshalled: %s", messageJSON)

	_, err = conn.Write(messageJSON)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Sent find contact message to %s\n", contact.Address)

	// Wait for a response
	responseBuffer := make([]byte, 6000)
	n, err := conn.Read(responseBuffer)
	if err != nil {
		return nil, err
	}

	// Deserialize the response into a Message
	var responseMessage Message
	if err := json.Unmarshal(responseBuffer[:n], &responseMessage); err != nil {
		return nil, err
	}

	if responseMessage.Type != findContactResponse {
		return nil, errors.New("Unexpected response type")
	}

	fmt.Printf("Received find contact response: %s\n", string(responseBuffer))

	return responseMessage.Contacts, nil
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
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

	fmt.Printf("TEST response!! messagetype: %s, target: %s", message.Type, message.TargetContact.ID)

	closestContacts := network.kademlia.RoutingTable.FindClosestContacts(message.TargetContact.ID, bucketSize)
	response := Message{
		Type:     findContactResponse,
		Contacts: closestContacts,
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

	fmt.Printf("Sent find contact response to %s\n", address.String())

}

func (network *Network) SendFindDataResponse(hash string) {
	// TODO
}

func (network *Network) SendStoreResponse(data []byte) {
	// TODO
}

type Message struct {
	Type           string
	Contacts       []Contact
	TargetContact  *Contact
	HashedData     []byte
	DataHashString string
}
