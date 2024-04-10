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

> **WARNING**
>
> In some cases, when the active consumer is newly switched over, 
> it may start receiving new messages. This could result in the duplication of messages or them being received out of order.

## Shared

Messages are delivered in a round-robin distribution across consumers ~ worker pool

> **WARNING**
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

> **WARNING**
> 
> In Go client, the operation Consumer.Unsubscribe() 
> 
> It'll fail when performed on the **shared** subscription
> where more than one consumer are currently connected. 
> 
> Mind to handle the error, the last consumer will perform clear the subscription.

