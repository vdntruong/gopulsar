package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"
)

type Producer struct {
	producer pulsar.Producer
}

func NewProducer(client pulsar.Client, topic string) (*Producer, func(), error) {
	p, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create producer: %w", err)
	}

	producer := &Producer{producer: p}
	closeFunc := func() {
		p.Close()
	}

	return producer, closeFunc, nil
}

func (p *Producer) Send(ctx context.Context, payload interface{}) (string, error) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	msgID, err := p.producer.Send(ctx, &pulsar.ProducerMessage{
		Payload: jsonPayload,
	})
	if err != nil {
		return "", fmt.Errorf("failed to send message: %w", err)
	}

	// extract message information here
	return msgID.String(), nil
}
