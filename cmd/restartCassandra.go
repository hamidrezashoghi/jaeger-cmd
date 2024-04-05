package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

// restartCassandraCmd represents the restartCassandra command
var restartCassandraCmd = &cobra.Command{
	Use:   "restart-cassandra",
	Short: "Restart Cassandra when leaves the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		restartCassandra()
	},
}

func init() {
	rootCmd.AddCommand(restartCassandraCmd)
}

func restartCassandra() {
	status := 0

	for true {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := exec.CommandContext(ctx, "/usr/bin/nodetool", "status").CombinedOutput()
		if err != nil {
			if status == 0 {
				fmt.Println("Error running nodetool status:", err)
				err := exec.Command("systemctl", "restart", "cassandra").Run()
				if err != nil {
					fmt.Println("Error restarting Cassandra:", err)
				} else {
					fmt.Println("Cassandra restarted successfully.")
					status = 1
				}
			}
		} else {
			fmt.Println("Cassandra is healthy.")
			status = 0
		}

		time.Sleep(60 * time.Second)
	}
}
