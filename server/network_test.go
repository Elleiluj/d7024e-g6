package server

import (
	"fmt"
	"log"
	"net"
	"strings"
	"testing"
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

// Listen
// func TestListen(t *testing.T) {
// 	adress := "localhost:8000"
// 	kademliaNode := NewKademliaNode(adress)

// 	Network := NewNetwork(&kademliaNode)

// 	//New
// 	localIP := GetOutboundIP()
// 	fmt.Println("Local IP:", localIP.String())

// 	port := 4000
// 	portStr := strconv.Itoa(port)

// 	localIPFull := localIP.String() + ":" + portStr
// 	//old

// 	Network.Listen(localIPFull) // <- works, but stops to listen

// 	got := Network.kademlia.Me.Address
// 	want := kademliaNode.Me.Address

// 	if got != want {
// 		t.Errorf("got %s want %s", got, want)
// 	} else {
// 		fmt.Println("NewNetwork \tPASS")
// 	}

// }

// handleResponse
// 		SendPingResponse
// 		SendFindContactResponse
// 		SendFindDataResponse
// 		SendStoreResponse
// 			sendResponse

// SendPingMessage
// SendFindContactMessage
// SendFindDataMessage
// SendStoreMessage
// 		SendMessage
