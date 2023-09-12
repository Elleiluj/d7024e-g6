package main

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func main() {

	// Test
	fmt.Println("Hello!!")

	// Create a Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	// List containers on the same network, without specifying network
	/*containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}*/

	// List containers on the same network, specifying network
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "network",
			Value: "comdockerdevenvironmentscode_kademlia_network",
		}),
	})
	if err != nil {
		panic(err)
	}

	// Print all containers
	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}

	// Container you want to ping
	targetContainerID := containers[0].ID

	sourceContainerID := containers[1].ID

	fmt.Println("Source Container ID:", sourceContainerID)
	fmt.Println("Target Container ID:", targetContainerID)

	// Run a command in the source container to ping the target container
	response, err := cli.ContainerExecCreate(context.Background(), sourceContainerID, types.ExecConfig{
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		Cmd:          []string{"ping", targetContainerID, "-c", "3"},
	})
	if err != nil {
		panic(err)
	}

	// Start the exec instance to run the command
	execResp, err := cli.ContainerExecAttach(context.Background(), response.ID, types.ExecStartCheck{})
	if err != nil {
		panic(err)
	}

	defer execResp.Close()

	buffer := make([]byte, 1024) // Create a buffer for reading
	bytesRead, err := execResp.Reader.Read(buffer)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Process the data
	data := buffer[:bytesRead]
	fmt.Println(string(data))

}
