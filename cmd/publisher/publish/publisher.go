package publish

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/spf13/cobra"

	"gopulsar/pkg/pubsub"
)

var topic string
var host string
var port string

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "publish a message",
	Long:  "just publish a message",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Args: %v\n", args)
		if len(args) < 1 {
			log.Printf("Need message payload. You passed %d args\n", len(args))
			os.Exit(1)
		}

		message := args[0]

		clusterAddr := fmt.Sprintf("pulsar://%s:%s", host, port)
		client, err := pulsar.NewClient(pulsar.ClientOptions{
			URL:               clusterAddr,
			OperationTimeout:  30 * time.Second,
			ConnectionTimeout: 30 * time.Second,
		})
		if err != nil {
			log.Fatalf("Could not instantiate Pulsar (%s) client: %v", clusterAddr, err)
		}
		defer client.Close()

		producer, closer, err := pubsub.NewProducer(client, topic)
		if err != nil {
			log.Fatalf("Could not instantiate Pulsar producer: %v", err)
		}
		defer closer()

		messageID, err := producer.Send(context.Background(), message)
		if err != nil {
			log.Fatalf("Could not send message: %v", err)
		}

		log.Println("Published message:", messageID)
	},
}

func AddCommand(rootCmd *cobra.Command) {
	log.Println("zo2")
	rootCmd.AddCommand(publishCmd)
	publishCmd.Flags().StringVar(&topic, "topic", "default", "Topic to publish messages to")
	publishCmd.Flags().StringVar(&host, "host", "localhost", "Cluster host to publish messages to")
	publishCmd.Flags().StringVar(&port, "port", "6650", "Cluster port to publish messages to")
}
