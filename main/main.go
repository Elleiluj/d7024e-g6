package main

import (
	"fmt"
	"kademlia/client"
	"kademlia/router"
	"kademlia/server"
	"log"
	"net"
	"strconv"
	"strings"
)

func main() {

	localIP := GetOutboundIP()
	fmt.Println("Local IP:", localIP.String())

	port := 4000
	portStr := strconv.Itoa(port)

	localIPFull := localIP.String() + ":" + portStr
	fmt.Println("Full local ip: ", localIPFull)

	me := server.NewKademliaNode(localIPFull)

	bootstrapIP := GetBootstrapIP(localIP.String())
	fmt.Println("bootstrapIP: ", bootstrapIP)

	bootstrapIPFull := bootstrapIP + ":" + portStr
	fmt.Println("Full bootstrapIP: ", bootstrapIPFull)

	if localIP.String() != bootstrapIP {
		fmt.Println("Retrieving bootstrap node...")
		bootstrapAddress := server.CreateHash(bootstrapIPFull)
		bootstrapContact := server.NewContact(server.NewKademliaID(bootstrapAddress), bootstrapIPFull)
		fmt.Println("Joining network...")
		me.JoinNetwork(&bootstrapContact)
	} else {
		fmt.Println("Initializing network with bootstrap node...")
	}

	network := server.NewNetwork(&me)

	go network.Listen(localIPFull)
	go me.RemoveExpiredData()
	go me.RefreshUploadedData()

	// RESTful
	router := router.NewRouter(&me)
	router.DefineHandleFunc()
	go router.StartHTTP()

	// CLI
	client := client.NewClient(&me)
	client.Start()

}

// Get preferred outbound ip of this machine
// https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
// establishes udp connection, and gets local ip from that
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

// IP is without port
func GetBootstrapIP(ip string) string {
	IPparts := strings.Split(ip, ".")
	firstPart := IPparts[:len(IPparts)-1]
	bootstrapIP := strings.Join(firstPart, ".") + "." + "2"
	return bootstrapIP
}
