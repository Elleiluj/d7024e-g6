package client

import (
	"bufio"
	"errors"
	"fmt"
	"kademlia/server"
	"os"
	"strings"
	"time"
)

// define put, get, exit etc. also in put, define ping

type Client struct {
	kademlia  *server.Kademlia
	sleepTime time.Duration
}

func NewClient(kademlia *server.Kademlia) *Client {
	client := Client{kademlia: kademlia}
	client.sleepTime = 0 * time.Second
	return &client
}

func (client *Client) Start() {
	reader := bufio.NewReader(os.Stdin)
	putString := "put <string> (takes a single argument, the contents of the file you are uploading, and outputs thehash of the object)\n"
	getString := "get <hash> (takes a hash as its only argument, and outputs the contents of the object and the node it was retrieved from)\n"
	forgetString := "forget <hash> (takes hash of the object that is no longer to be refreshed, only works on original uploader)\n"
	exitString := "exit (terminates this node)\n"
	//fmt.Printf("Commands:\n" + putString + getString + forgetString + exitString)

	for {
		fmt.Printf("\nCommands:\n" + putString + getString + forgetString + exitString)
		fmt.Printf("\nEnter a command: \n")

		// Read from terminal
		input, _ := reader.ReadString('\n')
		input = strings.TrimRight(input, "\n")

		splitInput := strings.SplitN(input, " ", 2)

		err := client.HandleInput(splitInput)

		if err != nil {
			fmt.Printf("%s", err)
		}

		time.Sleep(client.sleepTime)
	}
}

func (client *Client) HandleInput(input []string) error {
	var err error
	command := input[0]
	switch command {
	case "":
	case "put":
		data := input[1]
		err = client.put(data)
	case "get":
		data := input[1]
		args := strings.Fields(data)
		if len(args) > 1 {
			return errors.New("\nToo many arguments, can only take one\n")
		}
		if len(data) != 64 {
			return errors.New("\nLength of hash must be exactly 64 characters\n")
		}
		//err = client.get(data)
		client.get(data)
	case "forget":
		data := input[1]
		args := strings.Fields(data)
		if len(args) > 1 {
			return errors.New("\nToo many arguments, can only take one\n")
		}
		if len(data) != 64 {
			return errors.New("\nLength of hash must be exactly 64 characters\n")
		}
		client.forget(data)
	case "exit":
		fmt.Println("\nTerminating node...")
		os.Exit(0)
	default:
		return errors.New("\nInvalid command: " + command + "\n")
	}
	return err
}
