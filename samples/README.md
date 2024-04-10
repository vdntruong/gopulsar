# Subscription samples

## Exclusive

Let's start our consumers
```bash
go run ../cmd/cli/cli.go consumer -n=ConsumerE1 -t=topic-logger --subName=logcollector --subType=exclusive
```

```bash
go run ../cmd/cli/cli.go consumer -n=ConsumerE2 -t=topic-logger --subName=logcollector --subType=exclusive
```

Let's publish the message

```bash
go run ../cmd/cli/cli.go publish -t=topic-logger Only1ConsumerReceiveThisMsg
```

## Failover

### Failover non-partitioned topic

The broker picks consumers in the order they subscribe to topics

#### If there is one non-partitioned topic (consumer subscribes to only 1 topic)

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

#### If there are multiple non-partitioned topics (consumer subscribes to more than 1 topic)

A consumer is selected based on consumer name hash and topic name hash, eventually the 
assignment cannot be determined

```bash
go run ../cmd/cli/cli.go consumer -n=ConsumerFN1 -t=topic-01 -t=topic-02 -t=topic-03 -t=topic-04 --subName=checking --subType=failover
```

```bash
go run ../cmd/cli/cli.go consumer -n=ConsumerFN2 -t=topic-01 -t=topic-02 -t=topic-03 -t=topic-04 --subName=checking --subType=failover
```

> [!NOTE]
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

## Key_shared
