# Subscription samples

## Exclusive

Allows **only a single** consumer to attach to the subscription

Let's start our consumers
```bash
go run ../cmd/cli/cli.go consumer -n=ConsumerE1 -t=topic-logger --subName=logcollector --subType=exclusive
```

```bash
go run ../cmd/cli/cli.go consumer -n=ConsumerE2 -t=topic-logger --subName=logcollector --subType=exclusive
```

> [!CAUTION]
>
> We'll get an ERROR when trying to subscribe another consumer to the Exclusive subscription type
>

Let's publish the message

```bash
go run ../cmd/cli/cli.go publish -t=topic-logger Only1ConsumerReceiveThisMsg
```

## Failover

Allows multiple consumers attach to the same subscription. A master consumer is picked and receives messages. 

> [!IMPORTANT]
>
> In some cases, when the active consumer is newly switched over, 
> it may start receiving new messages. This could result in the duplication of messages or them being received out of order.

### Failover non-partitioned topic

The broker picks consumers in the order they subscribe to topics

Let's start 2 consumers

```bash
go run ../cmd/cli/cli.go consumer -n=ConsumerFN1 -t=topic-payment --subName=reconciliation --subType=failover
```

```bash
go run ../cmd/cli/cli.go consumer -n=ConsumerFN2 -t=topic-payment --subName=reconciliation --subType=failover
```

Publish some messages, then disconnect the first consumer to switch delivering to the second one

```bash
go run ../cmd/cli/cli.go publish -t=topic-payment TheFirstMessage TheSecondOne TheThirdOne AndTheLast
```

#### If there are multiple non-partitioned topics

A consumer is selected based on consumer name hash and topic name hash, eventually the 
assignment cannot be determined

```bash
go run ../cmd/cli/cli.go consumer -n=ConsumerFN1 -t=topic-01 -t=topic-02 -t=topic-03 -t=topic-04 --subName=checking --subType=failover
```

```bash
go run ../cmd/cli/cli.go consumer -n=ConsumerFN2 -t=topic-01 -t=topic-02 -t=topic-03 -t=topic-04 --subName=checking --subType=failover
```

> [!TIP]
> 
> Consumers can subscribe to multiple topics with the same subscription name

Publish some messages, then disconnect first consumer to switch delivering to the second one

```bash
go run ../cmd/cli/cli.go publish -t=topic-01 msg01Topic01 msg02Topic01 msg03Topic01
```

```bash
go run ../cmd/cli/cli.go publish -t=topic-02 msg01Topic02 msg02Topic02 msg03Topic02
```

```bash
go run ../cmd/cli/cli.go publish -t=topic-03 msg01Topic03 msg02Topic03 msg03Topic03
```

```bash
go run ../cmd/cli/cli.go publish -t=topic-04 msg01Topic04 msg02Topic04 msg03Topic04
```

> [!IMPORTANT]
>
> It does not matter which consumer starts first, the assignment will always follow the name hash.
> The assignment may change if you alter the consumer names.


### Failover partitioned topic

## Shared

Messages are delivered in a round-robin distribution across consumers ~ worker pool

> [!IMPORTANT]
> 
> Shared subscriptions **do not** guarantee message ordering or support cumulative acknowledgment.
> 
> Messages can be acknowledged in two ways:
> - Being acknowledged individually ~ The consumer acknowledges each message
> - Being acknowledged cumulatively ~ The consumer **only** acknowledges the last message it received

Let's start our consumers

```bash
go run ../cmd/cli/cli.go consumer -n=ConsumerS1 -t=topic-payment --subName=reconciliation --subType=shared
```

```bash
go run ../cmd/cli/cli.go consumer -n=ConsumerS2 -t=topic-payment --subName=reconciliation --subType=shared
```

```bash
go run ../cmd/cli/cli.go consumer -n=ConsumerS3 -t=topic-payment --subName=reconciliation --subType=shared
```

Let's publish some messages

```bash
go run ../cmd/cli/cli.go publish -t=topic-payment TheFirstMessage TheSecondOne TheThirdOne AndTheLast
```

> [!WARNING]
> 
> In Go client, the operation Consumer.Unsubscribe() 
> 
> It'll fail when performed on the **shared** subscription
> where more than one consumer are currently connected. 
> 
> Mind to handle the error, the last consumer will perform clear the subscription.

