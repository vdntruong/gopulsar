# Consumer

## Message retention and expiry

By default, Pulsar message brokers:

1. immediately delete all messages that have been acknowledged by a consumer, and
2. persistently store all unacknowledged messages in a message backlog

Pulsar has two features:
1. Message **retention** allows for the storage of **acknowledged messages**
2. Message **expiry** allows you to set a time-to-live (TTL) for 
 **messages that have not been acknowledged**

[Refer article](https://pulsar.apache.org/docs/3.2.x/cookbooks-retention-expiry)
## Acknowledgment

Consumed message will be **permanently stored** and **deleted only after all the subscriptions 
have acknowledged it**.

Acknowledgment (ack) is Pulsar's way of knowing that the message can be deleted from the system.

Messages can be acknowledged in 2 following ways: **individually** & **cumulatively**

```
// individually
if err := consumer.Ack(msg); err != nil {
    return err
}
```

```
// cumulatively
if err := consumer.AckCumulative(msg); err != nil {
    return err
}
```

With **cumulative acknowledgment**, the consumer **only acknowledges the last message** it received.

> [!NOTE]  
> In **Shared** subscription type, messages are acknowledged individually

## Message redelivery

It is important to have a built-in mechanism that handles failure,
particularly in asynchronous messaging as highlighted in the following examples:

- Consumers get disconnected from the database or the HTTP server.
    - database is temporarily offline while the consumer is writing the data to it.
    - external HTTP server that the consumer calls are momentarily unavailable (500 or Timeout...).

- Consumers get disconnected from a broker
    - consumer crashes (panic, coding...)
    - broken connections (network...)

Message redelivery in Apache Pulsar using **at-least-once delivery semantics**
that ensure Pulsar **processes a message more than once.**

To activate the message redelivery mechanism in Apache Pulsar using three methods:

#### Negative Acknowledgment

When a consumer fails to consume a message and needs to re-consume it,
the consumer sends **a negative acknowledgment (nack)** to the broker,
triggering the broker to redeliver this message to the consumer.

```
msg, _ := consumer.Receive(ctx)
// ... error occur
consumer.Nack(msg)
```

#### Acknowledgment Timeout

We can set a time range during which the client tracks the unacknowledged messages.
After this acknowledgment timeout (**ackTimeout**) period,
the client sends redeliver unacknowledged messages request to the broker,
thus the broker resends the unacknowledged messages to the consumer.

> [!IMPORTANT]  
> Golang client do not support, [refer issue](https://github.com/apache/pulsar-client-go/issues/403)

#### Retry letter topic

Retry letter topic allows you to store the messages that failed to be consumed
and retry consuming them later.

Consumers on the original topic are **automatically subscribed to the retry letter topic** as well.
Once the maximum number of retries has been reached,
the unconsumed messages are moved to a dead letter topic for manual processing.

```
consumer, _ := client.Subscribe(pulsar.ConsumerOptions{
    Name:             name,
    Type:             subType,
    Topic:            topic,
    SubscriptionName: subName,

    RetryEnable: true,
    DLQ: &pulsar.DLQPolicy{
        MaxDeliveries:    5,
        RetryLetterTopic: fmt.Sprintf("%s-RETRY", topic),
    },
})

msg, _ := consumer.Receive(ctx)
// ... error occur
consumer.ReconsumeLater(msg)
```

> [!TIP]  
> Redelivery backoff mechanism 
