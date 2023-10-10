package client

func (client *Client) put(data string) error {
	err := client.kademlia.Store([]byte(data))
	return err
}

func (client *Client) get(data string) {
	client.kademlia.LookupData(data)
}

func (client *Client) forget(data string) {
	client.kademlia.Forget(data)
}
