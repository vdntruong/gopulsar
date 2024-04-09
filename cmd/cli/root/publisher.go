package root

import (
	"log"
	"os"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/spf13/cobra"

	"gopulsar/pkg/pubsub"
)

var publisherTopic string

var publisherCmd = &cobra.Command{
	Use:   "publish",
	Short: "push a message",
	Long:  "just push a message to the topic",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Printf("Need message payload. You passed %d args\n", len(args))
			os.Exit(1)
		}

		client, err := pulsar.NewClient(pulsar.ClientOptions{
			URL:               PulsarURL,
			OperationTimeout:  30 * time.Second,
			ConnectionTimeout: 30 * time.Second,
		})
		if err != nil {
			log.Fatalf("Could not instantiate Pulsar (%s) client: %v", PulsarURL, err)
		}
		defer client.Close()

		producer, closer, err := pubsub.NewProducer(client, publisherTopic)
		if err != nil {
			log.Fatalf("Could not instantiate Pulsar producer: %v", err)
		}
		defer closer()

		for i := 0; i < len(args); i++ {
			messageID, err := producer.Send(cmd.Context(), args[i])
			if err != nil {
				log.Printf("ERROR: Could not send the message, err: %v, payload: %s\n", err, args[i])
				return
			}
			log.Println("Published message:", messageID, "payload:", args[i])
		}
	},
}

func init() {
	rootCmd.AddCommand(publisherCmd)

	publisherCmd.Flags().StringVarP(&publisherTopic, "topic", "t", "default", "Topic to push the message")
}
