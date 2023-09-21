package main

import (
	"fmt"
	"net"
	"strings"
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
func Listen(ip string) {
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
	buffer := make([]byte, 4000)

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
		addr.Port = port
		handleResponse(string(data), addr, conn)
	}
}

func handleResponse(data string, address *net.UDPAddr, conn *net.UDPConn) {
	network := &Network{}
	messageSplit := strings.Split(data, ";")
	messageType := messageSplit[0]
	switch messageType {
	case pingMessage:
		network.SendPingResponse(address, conn)
	case findContactMessage:
		fmt.Println("Find contact message!")
		results := network.HandleFindContact(messageSplit[1], messageSplit[2])
		network.SendFindContactResponse(&results, address, conn)
	case findContactResponse:
		fmt.Println("Find contact response!")
	case findDataMessage:
		fmt.Println("Find data message!")
	case findDataResponse:
		fmt.Println("Find data response!")
	case storeMessage:
		fmt.Println("Store message!")
	case storeResponse:
		fmt.Println("Store response!")
	default:
	}

}

func (network *Network) SendPingMessage(contact *Contact) {
	// Establish udp connection
	conn, err := net.Dial("udp", contact.Address)
	if err != nil {
		fmt.Printf("Error connecting to contact: %v\n", err)
	}
	defer conn.Close()

	// Send the ping message as bytes to the contact's address
	// TODO: hande response, this is just temporary
	_, err = conn.Write([]byte(pingMessage))
	if err != nil {
		fmt.Printf("Error writing to contact: %v\n", err)
	}

	fmt.Printf("Sent ping message to %s\n", contact.Address)

}

func (network *Network) SendFindContactMessage(sender *Kademlia, receiver *Contact, target *Contact) {
	// Establish udp connection
	conn, err := net.Dial("udp", receiver.Address)
	if err != nil {
		fmt.Printf("Error connecting to contact: %v\n", err)
	}
	defer conn.Close()

	// Send the ping message as bytes to the contact's address
	// TODO: hande response, this is just temporary
	message := findContactMessage + ";" + sender.me.String() + ";" + target.String()
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("Error writing to contact: %v\n", err)
	}

	fmt.Printf("Sent find contact message to %s\n", receiver.Address)

	fmt.Printf("SendFindContactMessage - sender: %s, toFind: %s", sender.me.Address, receiver.Address)
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func (network *Network) SendPingResponse(address *net.UDPAddr, conn *net.UDPConn) {

	// Send the ping response as bytes to the contact's address
	_, err := conn.WriteToUDP([]byte(pingResponse), address)
	if err != nil {
		fmt.Println("Error sending UDP response:", err)
	}

	fmt.Printf("Sent ping response to %s\n", address.String())

}

func (network *Network) SendFindContactResponse(contacts *[]Contact, address *net.UDPAddr, conn *net.UDPConn) {
	message := ""
	// Send the ping response as bytes to the contact's address
	_, err := conn.WriteToUDP([]byte(pingResponse), address)
	if err != nil {
		fmt.Println("Error sending UDP response:", err)
	}

	fmt.Printf("Sent ping response to %s\n", address.String())
}

func (network *Network) SendFindDataResponse(hash string) {
	// TODO
}

func (network *Network) SendStoreResponse(data []byte) {
	// TODO
}

func (network *Network) HandleFindContact(sender string, target string) []Contact {
	// results := sender.LookupContact(target)
	// return results
	return []Contact{}
}
