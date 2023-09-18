package main

import (
	"fmt"
	"net"
)

type Network struct {
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
	buffer := make([]byte, 1024)

	for {
		// Read data from UDP connection
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP connection:", err)
			continue
		}

		// Process received data
		data := buffer[:n]
		fmt.Printf("Received UDP message from %s: %s\n", addr, string(data))

		// Respond to the client
		// maybe do in separate function, to simplify testing
		// depending on what the request was, we want different responses
		response := []byte("Hello from UDP server")
		_, err = conn.WriteToUDP(response, addr)
		if err != nil {
			fmt.Println("Error sending UDP response:", err)
		}
	}
}

func handleResponse(data string, address string) {

}

func (network *Network) SendPingMessage(contact *Contact) {
	// Establish udp connection
	conn, err := net.Dial("udp", contact.Address)
	if err != nil {
		fmt.Printf("Error connecting to contact: %v\n", err)
	}
	defer conn.Close()

	pingMessage := "Ping"

	// Send the ping message as bytes to the contact's address
	// TODO: hande response, this is just temporary
	_, err = conn.Write([]byte(pingMessage))
	if err != nil {
		fmt.Printf("Error writing to contact: %v\n", err)
	}

	fmt.Printf("Sent ping message to %s\n", contact.Address)

}

func (network *Network) SendFindContactMessage(contact *Contact) []Contact {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
