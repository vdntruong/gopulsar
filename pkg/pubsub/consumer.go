package pubsub

import (
	"context"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
)

type Consumer struct {
	c pulsar.Consumer
}

type (
	MessageHandler func(ctx context.Context, msg pulsar.Message, ackFunc func() error) error
)

func NewConsumer(client pulsar.Client, topic string, name string) (*Consumer, func(), error) {
	c, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            topic,
		SubscriptionName: name,
		Type:             pulsar.Shared,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to subscribe to topic %s: %w", topic, err)
	}

	consumer := &Consumer{c: c}
	closeFunc := func() {
		c.Close()
	}

	return consumer, closeFunc, nil
}

func (c *Consumer) Consume(ctx context.Context) ([]byte, error) {
	msg, err := c.c.Receive(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to receive message: %w", err)
	}

	if err := c.c.Ack(msg); err != nil {
		return nil, fmt.Errorf("failed to ack message: %w", err)
	}
	return msg.Payload(), nil
}

func (c *Consumer) Unsubscribe() error {
	if err := c.c.Unsubscribe(); err != nil {
		return fmt.Errorf("failed to unsubscribe: %w", err)
	}
	return nil
}
