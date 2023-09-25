package client

import (
	"bufio"
	"errors"
	"fmt"
	"kademlia/server"
	"os"
	"strings"
)

// define put, get, exit etc. also in put, define ping

type Client struct {
	kademlia *server.Kademlia
}

type ClientMessage struct {
	Type    string
	Data    []byte
	Contact server.Contact
}

func NewClient(kademlia *server.Kademlia) *Client {
	return &Client{kademlia: kademlia}
}

func (client *Client) Start() {
	reader := bufio.NewReader(os.Stdin)
	putString := "put <string> (takes a single argument, the contents of the file you are uploading, and outputs thehash of the object)\n"
	getString := "get <hash> (takes a hash as its only argument, and outputs the contents of the object and the node it was retrieved from)\n"
	exitString := "exit (terminates this node)\n"
	fmt.Printf("Commands:\n" + putString + getString + exitString)

	for {
		fmt.Printf("\nCommands:\n" + putString + getString + exitString)
		fmt.Printf("\nEnter a command: \n")

		// Read from terminal
		input, _ := reader.ReadString('\n')
		input = strings.TrimRight(input, "\n")

		splitInput := strings.SplitN(input, " ", 2)

		err := client.HandleInput(splitInput)

		if err != nil {
			fmt.Printf("%s", err)
		}
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
		//err = client.get(data)
		client.get(data)
	case "exit":
		client.exit()
	default:
		return errors.New("\nInvalid command: " + command + "\n")
	}
	return err
}
