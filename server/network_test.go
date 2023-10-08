package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNetwork_print(t *testing.T) {
	fmt.Print("\n--------------------\n network.go\n--------------------\n")
}

// Imported from main
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
func GetBootstrapIP(ip string) string {
	IPparts := strings.Split(ip, ".")
	firstPart := IPparts[:len(IPparts)-1]
	bootstrapIP := strings.Join(firstPart, ".") + "." + "2"
	return bootstrapIP
}

// -------------- End Of Import ------------------- //

func TestNewNetwork(t *testing.T) {
	adress := "localhost:8000"
	kademliaNode := NewKademliaNode(adress)

	Network := NewNetwork(&kademliaNode)

	got := Network.kademlia.Me.Address
	want := kademliaNode.Me.Address

	if got != want {
		t.Errorf("got %s want %s", got, want)
	} else {
		fmt.Println("NewNetwork \tPASS")
	}

}

func TestListen(t *testing.T) {
	testAddr := "127.0.0.1:12345"

	invalidTestAddr := "127.0.0.1:123456"

	network := &Network{}

	fail := false

	go func() {
		network.Listen(invalidTestAddr)
	}()

	go func() {
		network.Listen(testAddr)
	}()

	clientConn, err := net.Dial("udp", testAddr)
	if err != nil {
		t.Fatalf("Failed to create UDP client connection: %v", err)
		fail = true
	}

	message := []byte("Hello, UDP server!")
	_, err = clientConn.Write(message)
	if err != nil {
		t.Fatalf("Failed to write to UDP server: %v", err)
		fail = true
	}

	time.Sleep(5 * time.Second)

	clientConn.Close()

	if !fail {
		fmt.Println("Listen \tPASS")
	}

}

func TestHandleResponse(t *testing.T) {
	mockAddr := "127.0.0.1:12344"
	udpAddr, _ := net.ResolveUDPAddr("udp", mockAddr)
	mockConn, _ := net.ListenUDP("udp", udpAddr)

	mockContact := NewContact(NewKademliaID(CreateHash(mockAddr)), mockAddr)
	hash := CreateHash("test")
	data := []byte("test")

	mockKademlia := NewKademliaNode(mockAddr)

	network := &Network{kademlia: &mockKademlia}

	fail := false

	err := network.handleResponse(nil, udpAddr, mockConn)
	if err != "error" {
		t.Errorf("got %s want %s", err, "error")
		fail = true
	}

	var messages []Message
	messages = append(messages, Message{Type: pingMessage},
		Message{Type: findContactMessage, Sender: &mockContact, TargetContact: &mockContact},
		Message{Type: findDataMessage, Sender: &mockContact, TargetContact: &mockContact, DataHashString: hash},
		Message{Type: storeMessage, Sender: &mockContact, TargetContact: &mockContact, HashedData: data, DataHashString: hash},
		Message{Type: resetTTLMessage, Sender: &mockContact, DataHashString: hash},
		Message{Type: "invalid"})

	for _, message := range messages {
		jasonMsg, _ := json.Marshal(message)
		got := network.handleResponse(jasonMsg, udpAddr, mockConn)
		want := message.Type
		if got != want {
			t.Errorf("got %s want %s", got, want)
			fail = true
		}
	}

	mockConn.Close()

	if !fail {
		fmt.Println("HandleResponse \tPASS")
	}

}

func TestSendPingMessage(t *testing.T) {

	var wg sync.WaitGroup
	mockAddr := "127.0.0.1:12344"
	udpAddr, _ := net.ResolveUDPAddr("udp", mockAddr)
	mockConn, _ := net.ListenUDP("udp", udpAddr)
	mockContact := NewContact(NewKademliaID(CreateHash(mockAddr)), mockAddr)
	mockKademlia := NewKademliaNode(mockAddr)
	network := &Network{kademlia: &mockKademlia}

	fail := false
	wg.Add(1)
	go func() {

		defer wg.Done()
		buffer := make([]byte, 1024)
		n, _, err := mockConn.ReadFromUDP(buffer)
		if err != nil {
			t.Errorf("Error reading from UDP server: %v", err)
			fail = true
			return
		}

		data := buffer[:n]
		var receivedMessage Message
		if err := json.Unmarshal(data, &receivedMessage); err != nil {
			t.Errorf("Error decoding received JSON message: %v", err)
			fail = true
			return
		}

		if receivedMessage.Type != pingMessage {
			t.Errorf("Received message type is not 'pingMessage'")
			fail = true
		}
	}()

	network.SendPingMessage(&mockContact, &mockContact)

	invalidAddr := "127.0.0.1:123446"
	invalidContact := NewContact(NewKademliaID(CreateHash(invalidAddr)), invalidAddr)
	network.SendPingMessage(&invalidContact, &invalidContact)

	wg.Wait()
	mockConn.Close()

	if !fail {
		fmt.Println("SendPingMessage \tPASS")
	}

}

func TestSendFindContactMessage(t *testing.T) {
	var wg sync.WaitGroup
	mockAddrSender := "127.0.0.1:12344"
	mockAddrContact := "127.0.0.1:12343"
	mockAddrTarget := "127.0.0.1:12342"
	udpAddr, _ := net.ResolveUDPAddr("udp", mockAddrContact)
	mockConn, _ := net.ListenUDP("udp", udpAddr)
	mockSender := NewContact(NewKademliaID(CreateHash(mockAddrSender)), mockAddrSender)
	mockKademlia := NewKademliaNode(mockAddrSender)
	network := &Network{kademlia: &mockKademlia}

	mockReceiver := NewContact(NewKademliaID(CreateHash(mockAddrContact)), mockAddrContact)
	mockTarget := NewContact(NewKademliaID(CreateHash(mockAddrTarget)), mockAddrTarget)

	sender := &mockSender
	contact := &mockReceiver
	target := &mockTarget

	fail := false

	defer mockConn.Close()
	wg.Add(1)
	go func() {

		defer wg.Done()
		buffer := make([]byte, 1024)
		n, addr, err := mockConn.ReadFromUDP(buffer)
		if err != nil {
			t.Errorf("Error reading from UDP server: %v", err)
			fail = true
			return
		}

		data := buffer[:n]
		var receivedMessage Message
		if err := json.Unmarshal(data, &receivedMessage); err != nil {
			t.Errorf("Error decoding received JSON message: %v", err)
			fail = true
			return
		}

		// Simulate processing the received message and constructing a response
		if receivedMessage.Type == findContactMessage {
			response := Message{
				Type:     findContactMessage,
				Contacts: []Contact{*target}, // Simulate the response with the target contact
			}

			jsonResponse, err := json.Marshal(response)
			if err != nil {
				t.Errorf("Error encoding JSON response: %v", err)
				fail = true
				return
			}

			_, err = mockConn.WriteToUDP(jsonResponse, addr)
			if err != nil {
				t.Errorf("Error writing response to UDP server: %v", err)
				fail = true
			}
		}
	}()

	contacts, err := network.SendFindContactMessage(sender, contact, target)

	wg.Wait()

	if err != nil {
		t.Errorf("Error sending FindContactMessage: %v", err)
		fail = true
	}
	if len(contacts) != 1 || contacts[0].ID.String() != target.ID.String() {
		t.Errorf("Unexpected response contacts: %v", contacts)
		fail = true
	}
	if !fail {
		fmt.Println("SendFindContactMessage \tPASS")
	}
}

func TestSendFindDataMessage(t *testing.T) {
	var wg sync.WaitGroup
	mockAddrSender := "127.0.0.1:12344"
	mockAddrContact := "127.0.0.1:12343"
	mockAddrTarget := "127.0.0.1:12342"
	udpAddr, _ := net.ResolveUDPAddr("udp", mockAddrContact)
	mockConn, _ := net.ListenUDP("udp", udpAddr)
	mockSender := NewContact(NewKademliaID(CreateHash(mockAddrSender)), mockAddrSender)
	mockKademlia := NewKademliaNode(mockAddrSender)
	network := &Network{kademlia: &mockKademlia}

	mockReceiver := NewContact(NewKademliaID(CreateHash(mockAddrContact)), mockAddrContact)
	mockTarget := NewContact(NewKademliaID(CreateHash(mockAddrTarget)), mockAddrTarget)

	sender := &mockSender
	receiver := &mockReceiver
	target := &mockTarget

	fail := false
	hash := "mockHash"

	defer mockConn.Close()
	wg.Add(1)
	go func() {

		defer wg.Done()
		buffer := make([]byte, 1024)
		n, addr, err := mockConn.ReadFromUDP(buffer)
		if err != nil {
			t.Errorf("Error reading from UDP server: %v", err)
			fail = true
			return
		}

		data := buffer[:n]
		var receivedMessage Message
		if err := json.Unmarshal(data, &receivedMessage); err != nil {
			t.Errorf("Error decoding received JSON message: %v", err)
			fail = true
			return
		}

		// Simulate processing the received message and constructing a response
		if receivedMessage.Type == findDataMessage {
			response := Message{
				Type:       findDataMessage,
				Contacts:   []Contact{*target}, // Simulate the response with the target contact
				HashedData: []byte("mockData"), // Simulate the response with hashed data
			}

			jsonResponse, err := json.Marshal(response)
			if err != nil {
				t.Errorf("Error encoding JSON response: %v", err)
				fail = true
				return
			}

			_, err = mockConn.WriteToUDP(jsonResponse, addr)
			if err != nil {
				t.Errorf("Error writing response to UDP server: %v", err)
				fail = true
			}
		}
	}()

	contacts, hashedData, responseReceiver, err := network.SendFindDataMessage(sender, receiver, target, hash)

	wg.Wait()

	if err != nil {
		t.Errorf("Error sending FindDataMessage: %v", err)
		fail = true
	}
	if len(contacts) != 1 || contacts[0].ID.String() != target.ID.String() {
		t.Errorf("Unexpected response contacts: %v", contacts)
		fail = true
	}
	if string(hashedData) != "mockData" {
		t.Errorf("Unexpected response hashedData: %s", string(hashedData))
		fail = true
	}
	if responseReceiver.ID != receiver.ID {
		t.Errorf("Unexpected response receiver: %v", responseReceiver)
		fail = true
	}

	if !fail {
		fmt.Println("SendFindDataMessage \tPASS")
	}
}

func TestSendStoreMessage(t *testing.T) {
	var wg sync.WaitGroup
	mockAddrSender := "127.0.0.1:12344"
	mockAddrReceiver := "127.0.0.1:12343"
	udpAddr, _ := net.ResolveUDPAddr("udp", mockAddrReceiver)
	mockConn, _ := net.ListenUDP("udp", udpAddr)
	mockSender := NewContact(NewKademliaID(CreateHash(mockAddrSender)), mockAddrSender)
	mockKademlia := NewKademliaNode(mockAddrSender)
	network := &Network{kademlia: &mockKademlia}

	mockReceiver := NewContact(NewKademliaID(CreateHash(mockAddrReceiver)), mockAddrReceiver)

	sender := &mockSender
	receiver := &mockReceiver

	fail := false

	data := []byte("test")
	key := CreateHash("test")

	defer mockConn.Close()
	wg.Add(1)
	go func() {

		defer wg.Done()
		buffer := make([]byte, 1024)
		n, addr, err := mockConn.ReadFromUDP(buffer)
		if err != nil {
			t.Errorf("Error reading from UDP server: %v", err)
			fail = true
			return
		}

		data := buffer[:n]
		var receivedMessage Message
		if err := json.Unmarshal(data, &receivedMessage); err != nil {
			t.Errorf("Error decoding received JSON message: %v", err)
			fail = true
			return
		}

		// Simulate processing the received message and constructing a response
		if receivedMessage.Type == storeMessage {
			jsonResponse, err := json.Marshal(receivedMessage)
			if err != nil {
				t.Errorf("Error encoding JSON response: %v", err)
				fail = true
				return
			}
			_, err = mockConn.WriteToUDP(jsonResponse, addr)
			if err != nil {
				t.Errorf("Error writing response to UDP server: %v", err)
				fail = true
			}
		}
	}()

	err := network.SendStoreMessage(sender, receiver, data, key)
	wg.Wait()

	if err != nil {
		t.Errorf("Error sending StoreMessage: %v", err)
		fail = true
	}
	if !fail {
		fmt.Println("SendStoreMessage \tPASS")
	}
}

func TestSendResetTTLMessage(t *testing.T) {
	var wg sync.WaitGroup
	mockAddrSender := "127.0.0.1:12344"
	mockAddrReceiver := "127.0.0.1:12343"
	udpAddr, _ := net.ResolveUDPAddr("udp", mockAddrReceiver)
	mockConn, _ := net.ListenUDP("udp", udpAddr)
	mockSender := NewContact(NewKademliaID(CreateHash(mockAddrSender)), mockAddrSender)
	mockKademlia := NewKademliaNode(mockAddrSender)
	network := &Network{kademlia: &mockKademlia}

	mockReceiver := NewContact(NewKademliaID(CreateHash(mockAddrReceiver)), mockAddrReceiver)

	sender := &mockSender
	receiver := &mockReceiver

	fail := false
	key := CreateHash("test")

	defer mockConn.Close()
	wg.Add(1)
	go func() {
		defer wg.Done()
		buffer := make([]byte, 1024)
		n, addr, err := mockConn.ReadFromUDP(buffer)
		if err != nil {
			t.Errorf("Error reading from UDP server: %v", err)
			fail = true
			return
		}

		data := buffer[:n]
		var receivedMessage Message
		if err := json.Unmarshal(data, &receivedMessage); err != nil {
			t.Errorf("Error decoding received JSON message: %v", err)
			fail = true
			return
		}

		// Simulate processing the received message and constructing a response
		if receivedMessage.Type == resetTTLMessage {
			jsonResponse, err := json.Marshal(receivedMessage)
			if err != nil {
				t.Errorf("Error encoding JSON response: %v", err)
				fail = true
				return
			}
			_, err = mockConn.WriteToUDP(jsonResponse, addr)
			if err != nil {
				t.Errorf("Error writing response to UDP server: %v", err)
				fail = true
			}
		}
	}()
	err := network.SendResetTTLMessage(sender, receiver, key)

	wg.Wait()

	if err != nil {
		t.Errorf("Error sending ResetTTLMessage: %v", err)
		fail = true
	}
	if !fail {
		fmt.Println("SendResetTTLMessage \tPASS")
	}
}

func TestSendMessage(t *testing.T) {
	var wg sync.WaitGroup
	mockAddrSender := "127.0.0.1:12344"
	mockAddrReceiver := "127.0.0.1:12343"
	udpAddr, _ := net.ResolveUDPAddr("udp", mockAddrReceiver)
	mockConn, _ := net.ListenUDP("udp", udpAddr)
	mockSender := NewContact(NewKademliaID(CreateHash(mockAddrSender)), mockAddrSender)
	mockKademlia := NewKademliaNode(mockAddrSender)
	network := &Network{kademlia: &mockKademlia}

	mockReceiver := NewContact(NewKademliaID(CreateHash(mockAddrReceiver)), mockAddrReceiver)

	sender := &mockSender
	receiver := &mockReceiver

	fail := false

	message := Message{Type: "testing", Sender: sender}

	wg.Add(1)
	go func() {
		defer wg.Done()
		buffer := make([]byte, 1024)
		n, addr, err := mockConn.ReadFromUDP(buffer)
		if err != nil {
			t.Errorf("Error reading from UDP server: %v", err)
			fail = true
			return
		}

		data := buffer[:n]
		var receivedMessage Message
		if err := json.Unmarshal(data, &receivedMessage); err != nil {
			t.Errorf("Error decoding received JSON message: %v", err)
			fail = true
			return
		}

		// Simulate processing the received message and constructing a response
		if receivedMessage.Type == "testing" {
			jsonResponse, err := json.Marshal(receivedMessage)
			if err != nil {
				t.Errorf("Error encoding JSON response: %v", err)
				fail = true
				return
			}
			_, err = mockConn.WriteToUDP(jsonResponse, addr)
			if err != nil {
				t.Errorf("Error writing response to UDP server: %v", err)
				fail = true
			}
		}
	}()

	_, err := network.SendMessage(receiver, message)

	wg.Wait()

	if err != nil {
		t.Errorf("Error sending message: %v", err)
		fail = true
	}

	invalidMsg := Message{test: "test"}

	network.SendMessage(receiver, invalidMsg)

	wg.Wait()

	invalidContact := NewContact(NewKademliaID(CreateHash("1.1.1.1.1.1.1")), "1.1.1.1.1.1.1")

	network.SendMessage(&invalidContact, message)

	wg.Wait()

	mockConn.Close()

	mockConnInvalidResponseMsg, _ := net.ListenUDP("udp", udpAddr)
	wg.Add(1)
	go func() {
		defer wg.Done()
		buffer := make([]byte, 1024)
		_, addr, _ := mockConnInvalidResponseMsg.ReadFromUDP(buffer)
		_, err = mockConnInvalidResponseMsg.WriteToUDP(nil, addr)
	}()

	_, err = network.SendMessage(receiver, message)

	wg.Wait()

	mockConnInvalidResponseMsg.Close()

	if err == nil {
		t.Errorf("Unmarshalling nil")
		fail = true
	}

	if !fail {
		fmt.Println("SendMessage \tPASS")
	}

}

func TestSendPingResponse(t *testing.T) {
	var wg sync.WaitGroup
	mockAddr := "127.0.0.1:12344"
	udpAddr, _ := net.ResolveUDPAddr("udp", mockAddr)
	mockConn, _ := net.ListenUDP("udp", udpAddr)

	network := &Network{}

	fail := false
	expectedResponse := Message{
		Type: pingResponse,
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		buffer := make([]byte, 1024)
		n, _, err := mockConn.ReadFromUDP(buffer)
		if err != nil {
			t.Errorf("Error reading from UDP server: %v", err)
			fail = true
			return
		}

		data := buffer[:n]
		var receivedMessage Message
		if err := json.Unmarshal(data, &receivedMessage); err != nil {
			t.Errorf("Error decoding received JSON message: %v", err)
			fail = true
			return
		}
		if receivedMessage.Type != expectedResponse.Type {
			t.Errorf("Received message does not match expected response")
			fail = true
			return
		}
	}()
	network.SendPingResponse(udpAddr, mockConn)
	wg.Wait()

	mockConn.Close()

	invalidAddr := "127.0.0.1:123445"
	invalidUdpAddr, _ := net.ResolveUDPAddr("udp", invalidAddr)
	invldMockConn, _ := net.ListenUDP("udp", invalidUdpAddr)
	wg.Add(1)
	go func() {
		defer wg.Done()
		buffer := make([]byte, 1024)
		timeout := time.Now().Add(3 * time.Second)
		invldMockConn.SetDeadline(timeout)
		invldMockConn.ReadFromUDP(buffer)
	}()

	network.SendPingResponse(invalidUdpAddr, invldMockConn)

	wg.Wait()

	invldMockConn.Close()
	if !fail {
		fmt.Println("SendPingResponse \tPASS")
	}
}

func TestSendResponse(t *testing.T) {
	addr := "127.0.0.1:12344"
	udpAddr, _ := net.ResolveUDPAddr("udp", addr)
	conn, _ := net.ListenUDP("udp", udpAddr)

	network := &Network{}

	response := Message{
		Type: "testResponse",
	}

	fail := false

	var receivedMessage Message
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		buffer := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			t.Errorf("Error reading from UDP server: %v", err)
			fail = true
			return
		}

		data := buffer[:n]
		if err := json.Unmarshal(data, &receivedMessage); err != nil {
			t.Errorf("Error decoding received JSON message: %v", err)
			fail = true
			return
		}
	}()

	network.sendResponse(response, udpAddr, conn)

	wg.Wait()

	if receivedMessage.Type != response.Type {
		t.Errorf("Received message does not match expected response")
		fail = true
	}

	invalidMsg := Message{test: "test"}

	network.sendResponse(invalidMsg, udpAddr, conn)

	wg.Wait()

	conn.Close()

	invldAddr := "127.0.0.1:123445"
	invldUdpAddr, _ := net.ResolveUDPAddr("udp", invldAddr)
	invldConn, _ := net.ListenUDP("udp", invldUdpAddr)
	wg.Add(1)
	go func() {
		defer wg.Done()
		buffer := make([]byte, 1024)
		timeout := time.Now().Add(3 * time.Second)
		invldConn.SetDeadline(timeout)
		invldConn.ReadFromUDP(buffer)
	}()

	network.sendResponse(response, invldUdpAddr, invldConn)

	wg.Wait()

	invldConn.Close()

	if !fail {
		fmt.Println("SendResponse \tPASS")
	}
}
