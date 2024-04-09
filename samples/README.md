# Subscription samples

## Exclusive

Only allows a single consumer to attach to the subscription

```bash
go run cmd/cli/cli.go consumer -n=ConsumerA -t=topic-logger --subName=logcollector --subType=exclusive
```

```bash
go run cmd/cli/cli.go consumer -n=ConsumerB -t=topic-logger --subName=logcollector --subType=exclusive
```

> **Warning**
>
> We'll get an ERROR when trying to subscribe another consumer to the Exclusive subscription type
>

Let's publish the message

```bash
go run cmd/cli/cli.go publish -t=topic-logger Only1ConsumerReceiveThisMsg
```
