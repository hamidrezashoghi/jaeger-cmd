package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
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

type Alert struct {
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations,omitempty`
	GeneratorURL string
}

func restartCassandra() {
	status := 0
	hostnameWithDomain, _ := os.Hostname()
	hostname := strings.Split(hostnameWithDomain, ".")[0]

	alertManagerAPI := "https://alertmanager.local/api/v1/alerts"
	alert := Alert{
		Labels: map[string]string{
			"alertname":   "CassandraRestarted",
			"severity":    "warning",
			"instance":    hostname,
			"environment": "JaegerDatabase",
		},
		Annotations: map[string]string{
			"summary":     "Not need to inform, Cassandra service restarted due high load",
			"description": "Cassandra service restarted automatically on " + hostname,
		},
	}

	alerts := []Alert{alert}

	// Convert alerts to json
	data, err := json.Marshal(alerts)
	if err != nil {
		log.Println("Couldn't convert alerts to json")
	}

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

					// Send alert to Alertmanager
					resp, err := http.Post(alertManagerAPI, "application/json", bytes.NewBuffer(data))
					if err != nil {
						log.Println("Couldn't sent alert to Alertmanager")
					} else {
						fmt.Println("Alert sent to Alertmanager.")
					}
					defer resp.Body.Close()

					// Check response from Alertmanager
					if resp.StatusCode != http.StatusOK {
						log.Println("Failed to sent alert.")
					}
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
