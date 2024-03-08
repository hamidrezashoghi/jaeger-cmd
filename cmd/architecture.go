package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var ctx = context.Background()
var containerName string = "jaeger-architecture"

// architectureCmd represents the architecture command
var architectureCmd = &cobra.Command{
	Use:   "architecture",
	Short: "Generate relationship between traces - System Architecture in Jaeger UI",
	Run: func(cmd *cobra.Command, args []string) {
		architecture()
	},
}

func stopAndRemoveContainer(client *client.Client, containerName string) {
	ctx := context.Background()

	stopOptions := container.StopOptions{
		Signal:  "SIGTERM",
		Timeout: new(int),
	}

	// Set the Timeout to 0
	*stopOptions.Timeout = 0

	if err := client.ContainerStop(ctx, containerName, stopOptions); err != nil {
		fmt.Printf("Unable to stop %s container, %s\n", containerName, err)
	} else {
		fmt.Printf("%s container: stopped\n", containerName)
	}

	removeOptions := container.RemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}

	time.Sleep(2 * time.Second)
	if err := client.ContainerRemove(ctx, containerName, removeOptions); err != nil {
		fmt.Printf("Unable to remove %s container, %s\n", containerName, err)
	}
}

func runContainer(client *client.Client) (isRun bool) {
	storage := "cassandra"
	CASSANDRA_CONTACT_POINTS := "192.1.2.2,192.1.2.3,192.1.2.4,192.1.2.5,192.1.2.6"

	// Create a container config
	containerConfig := &container.Config{
		Image:     "jaegertracing/spark-dependencies:latest",
		Tty:       true,
		OpenStdin: true,
		Env: []string{
			"STORAGE=" + storage,
			"CASSANDRA_CONTACT_POINTS=" + CASSANDRA_CONTACT_POINTS,
		},
	}

	hostConfig := &container.HostConfig{
		AutoRemove: true,
	}

	// Create the container
	containerCreate, err := client.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, containerName)
	if err != nil {
		log.Fatalf("Couldn't create %s container, %s", containerName, err)
	}

	// Start the container
	if err := client.ContainerStart(ctx, containerCreate.ID, container.StartOptions{}); err != nil {
		log.Fatalf("Couldn't start %s container, %s", containerName, err)
	}

	return true
}

func architecture() {

	client, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)

	if err != nil {
		log.Fatalf("Unable to create docker client, %s", err)
	}

	// Stops and remove a container
	stopAndRemoveContainer(client, containerName)
	isRun := runContainer(client)
	if isRun == true {
		fmt.Printf("%s container: started\n", containerName)
	}
}

func init() {
	rootCmd.AddCommand(architectureCmd)
}
