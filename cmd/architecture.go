package cmd

import (
	"context"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

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
		log.Fatalf("Unable to stop %s container, %s\n", containerName, err)
	}

	removeOptions := container.RemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}

	if err := client.ContainerRemove(ctx, containerName, removeOptions); err != nil {
		log.Fatalf("Unable to remove %s container, %s\n", containerName, err)
	}
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
	stopAndRemoveContainer(client, "ubuntu-agent2")
}

func init() {
	rootCmd.AddCommand(architectureCmd)
}
