package pubsub

import (
	"context"
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"
)

type Consumer struct {
	client   pulsar.Client
	consumer pulsar.Consumer

	name    string
	topic   []string
	subName string
	subType pulsar.SubscriptionType
}

func NewConsumer(
	client pulsar.Client,
	name string,
	topic []string,
	subName string,
	subType pulsar.SubscriptionType,
) (*Consumer, error) {
	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Name:             name,
		Topics:           topic,
		SubscriptionName: subName,
		Type:             subType,
	})
	if err != nil {
		return nil, err
	}

	return &Consumer{
		client:   client,
		name:     name,
		topic:    topic,
		subName:  subName,
		subType:  subType,
		consumer: consumer,
	}, nil
}

func (c *Consumer) Close() error {
	if err := c.consumer.Unsubscribe(); err != nil {
		// Unsubscribe
		// This operation will fail when performed on a shared subscription
		// where more than one consumer are currently connected.
		if c.subType != pulsar.Shared {
			return fmt.Errorf("failed to unsubscribe: %w", err)
		}
	}
	c.consumer.Close()
	return nil
}

func (c *Consumer) PullMessages(ctx context.Context) error {
	for {
		msg, err := c.consumer.Receive(ctx)
		if err != nil {
			return err
		}

		fmt.Printf(
			"Received message msgId: %#v -- content: '%s'\n",
			msg.ID(), string(msg.Payload()),
		)

		c.consumer.Nack(msg)
		if err := c.consumer.Ack(msg); err != nil {
			return err
		}
	}
}
