package client

import (
	"fmt"
	"os"
)

func (client *Client) put(data string) error {
	err := client.kademlia.Store([]byte(data))
	return err
}

func (client *Client) get(data string) {
	client.kademlia.LookupData(data)
}

func (client *Client) exit() {
	fmt.Println("\nTerminating node...")
	os.Exit(0)
}
