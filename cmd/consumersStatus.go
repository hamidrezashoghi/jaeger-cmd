package cmd

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// consumersStatusCmd represents the consumersStatus command
var consumersStatusCmd = &cobra.Command{
	Use:   "consumers-status",
	Short: "Number of consumers connected to kafka",
	Run: func(cmd *cobra.Command, args []string) {
		consumersStatus()
	},
}

func init() {
	rootCmd.AddCommand(consumersStatusCmd)
}

func consumersStatus() {
	servers := []string{"192.168.1.2:9092", "192.168.1.3:9092", "192.168.1.4:9092"}

	groupName := "jaeger-ingester"
	timeout := "10000"
	kafkaPath := "/usr/local/kafka/bin/"
	consumerServersFile := "/var/lib/node_exporter/textfile_collector/consumer_servers.prom"
	countConsumers := make(map[string]int, 3)

	// Get number of consumers per kafka
	for _, server := range servers {
		cmdConsumer := exec.Command(kafkaPath+"kafka-consumer-groups.sh", "--dry-run",
			"--bootstrap-server", server, "--timeout", timeout, "--group="+groupName,
			"--members", "--describe")

		// Number of consumers
		out, err := cmdConsumer.Output()
		if err != nil {
			log.Fatalf("Couldn't get number of consumers on %s, %v\n", server, err)
		}

		outLines := strings.Split(string(out), "\n")

		// Drop header and blank lines from output command
		trueLines := outLines[2:len(outLines)]

		for _, line := range trueLines {
			vars := strings.Split(line, " ")
			if len(vars) < 3 {
				continue
			}

			// Change server IP to region
			switch {
			case server == "192.168.1.2:9092":
				countConsumers["local1"] += 1
			case server == "192.168.1.3:9092":
				countConsumers["local2"] += 1
			case server == "192.168.1.4:9092":
				countConsumers["local3"] += 1
			}
		}
	}

	file, err := os.OpenFile("consumer_servers.prom", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Couldn't open consumer_servers.prom file.")
	}

	defer file.Close()

	_ = os.Remove(consumerServersFile)
	_, _ = file.WriteString("# HELP consumer_servers Metric\n")
	_, _ = file.WriteString("# TYPE consumer_servers gauge\n")
	for region, consumers := range countConsumers {
		_, _ = file.WriteString("consumer_servers{region=\"" + region + "\"} " + strconv.Itoa(consumers) + "\n")
	}

	// Change ownership of consumer_servers.prom
	// execute below command to get uid and gid of node_exporter
	// # id node_exporter
	// uid=997(node_exporter) gid=998(node_exporter) groups=998(node_exporter)
	var uidNodeExporter int = 997
	var gidNodeExporter int = 998
	_ = os.Chown("consumer_servers.prom", uidNodeExporter, gidNodeExporter)

	// /var/lib/node_exporter/textfile_collector/
	if err := os.Rename("consumer_servers.prom", consumerServersFile); err != nil {
		log.Fatalln("Couldn't move consumer_servers.prom to /var/lib/node_exporter/textfile_collector/ path.", err)
	}
}