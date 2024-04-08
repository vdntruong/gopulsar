package root

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/spf13/cobra"

	"gopulsar/pkg/pubsub"
)

var (
	name  = "con-01"
	topic = "default"

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
		if err := startConsumer(cmd.Context(), name, topic, subName, subType); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(consumerCmd)

	consumerCmd.Flags().StringVarP(&name, "name", "n", "", "consumer name")
	consumerCmd.Flags().StringVarP(&topic, "topic", "t", "", "topic")
	consumerCmd.Flags().StringVarP(&subName, "subName", "", "", "subscription name")
	consumerCmd.Flags().StringVarP(&subType, "subType", "", "", "subscription type [exclusive, shared, failover, key_shared]")

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

func startConsumer(ctx context.Context, name, topic string, subName, subType string) error {
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
	defer client.Close()

	consumer, err := pubsub.NewConsumer(client, name, topic, subName, subscriptionType)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()
		if err := consumer.PullMessages(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx, &wg)

	wg.Wait()
	return nil
}
