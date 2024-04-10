package root

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/spf13/cobra"

	"gopulsar/pkg/pubsub"
)

var (
	name   = "con-01"
	topics = []string{"default"}

	subName = "sub-01"
	subType = "exclusive"
)

var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Consume messages",
	Long:  `Create a consumer and start pulling messages`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(name) == 0 {
			log.Println("")
		}
		if err := startConsumer(cmd.Context(), name, topics, subName, subType); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(consumerCmd)

	consumerCmd.Flags().StringVarP(&name, "name", "n", name, "consumer name")
	consumerCmd.Flags().StringArrayVarP(&topics, "topic", "t", topics, "topic")
	consumerCmd.Flags().StringVarP(&subName, "subName", "", subName, "subscription name")
	consumerCmd.Flags().StringVarP(&subType, "subType", "", subType, "subscription type [exclusive, shared, failover, key_shared]")

	_ = consumerCmd.MarkFlagRequired("name")
	_ = consumerCmd.MarkFlagRequired("topic")
	_ = consumerCmd.MarkFlagRequired("subName")
	_ = consumerCmd.MarkFlagRequired("subType")

	//consumerCmd.Flags().BoolP("name", "n", false, "Help message for toggle")
}

var (
	SubscriptionTypeMapper = map[string]pulsar.SubscriptionType{
		"exclusive":  pulsar.Exclusive,
		"shared":     pulsar.Shared,
		"failover":   pulsar.Failover,
		"key_shared": pulsar.KeyShared,
	}
)

func startConsumer(
	ctx context.Context,
	name string,
	topics []string,
	subName string,
	subType string,
) error {
	subscriptionType, ok := SubscriptionTypeMapper[subType]
	if !ok {
		return fmt.Errorf("subType have to in [exclusive, shared, failover, key_shared]")
	}

	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: PulsarURL,
	})
	if err != nil {
		return fmt.Errorf("could not instantiate Pulsar (%s) client: %w", PulsarURL, err)
	}
	defer func() {
		client.Close()
		log.Println("Closed pulsar client")
	}()

	consumer, err := pubsub.NewConsumer(client, name, topics, subName, subscriptionType)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatal(err)
		}
		log.Println("Closed consumer")
	}()

	log.Println("Registered consumer on topics", strings.Join(topics, ", "))
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-osSignals
		log.Println("Received interrupt signal. Shutting down...")
		if err := consumer.Close(); err != nil {
			log.Println("Failed to clear consumer: ", err)
		}
		client.Close()
		os.Exit(0)
	}()

	log.Println("Pulling messages")
	if err := consumer.PullMessages(ctx); err != nil {
		if strings.Contains(err.Error(), "ConsumerClosed") {
			log.Println("Consumer Closed")
		} else {
			log.Println("Failed to pull messages: ", err)
		}
	}

	return nil
}
