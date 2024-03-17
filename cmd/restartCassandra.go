package cmd

import (
	"github.com/spf13/cobra"
)

// restartCassandraCmd represents the restartCassandra command
var restartCassandraCmd = &cobra.Command{
	Use:   "restart-cassandra",
	Short: "Restart Cassandra when one of them doesn't respond",
	Run: func(cmd *cobra.Command, args []string) {
		restartCassandra()
	},
}

func init() {
	rootCmd.AddCommand(restartCassandraCmd)
}

func restartCassandra() {
	
}
