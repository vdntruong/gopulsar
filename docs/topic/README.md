# Topic

A Pulsar topic is a unit of storage (logical & abstract) that organizes messages into a stream.

Topic names are URLs that have a well-defined structure:
```
{persistent|non-persistent}://tenant/namespace/topic
```

By default, Pulsar persistently stores all unacknowledged messages 
on multiple BookKeeper bookies (storage nodes). 
Data for messages on persistent topics can thus survive broker restarts and subscriber failover.

```
{persistent}://tenant/namespace/topic
```

## Namespaces

A Pulsar namespace is a logical grouping of topics as well as a logical nomenclature 
within a tenant.

The topic `persistent://my-tenant/my-app/my-topic` is a namespace for 
the application `my-app` for `my-tenant`.

## Multi-topic subscriptions

Pulsar consumers can simultaneously subscribe to multiple topics.
Go client can define a list of topic in two ways:

- a regular expression (regex)
- explicitly defining a list of topic (pulsar.ConsumerOptions.Topics)

> [!WARNING]   
> No ordering guarantees across multiple topics

## Partitioned topics

**For higher throughput**

Normal topics are served only by a single broker, which limits the maximum throughput of the topic. 
Partitioned topic is a special type of topic handled by multiple brokers.

A partitioned topic is implemented as **N internal topics**. 
The distribution of partitions across brokers is handled automatically by Pulsar.

> [!TIP]   
> **Throughput** concerns should guide **partitioning/routing** decisions 
> while **subscription** decisions should be guided by **application semantics**

### Routing modes

The routing mode determines each message should be published to which partition

#### Round Robin Partition mode

> [!NOTE]   
> RoundRobinPartition is the default mode

The **key is provided**, 
the partitioned producer will **hash the key and assign** message to a particular partition.

If **no key is provided**, 
the producer will publish messages across all partitions in **round-robin fashion** 
to achieve maximum throughput.

> [!IMPORTANT]  
> Round-robin is not done per individual message 
> but rather it's set to the same boundary of batching delay, to ensure batching is effective

```
msg, _ := pulsar.NewDefaultRouter(args...)
msg.Payload = payload
msg.Key = key

msgID, _ := product.Send(ctx, msg)
```

#### Single Partition mode

The **key is provided**,
the partitioned producer will **hash the key and assign** message to a particular partition.

If **no key is provided**,
the producer will **randomly pick one single partition** and publish all the messages into
that partition.

```
msg, topicMetadata, err := pulsar.NewSinglePartitionRouter()
msg.Payload = payload
msg.Key = key

msgID, _ := product.Send(ctx, msg)
```

#### Custom Partition mode

Users have to **implement the MessageRouter interface** to determine the partition 
for a particular message

### Ordering guarantee

#### Per-key-partition

Use either `SinglePartition` or `RoundRobinPartition` mode, and **Key is provided** by each message.

#### Per-producer

All the messages from the same producer will be in order.

Use `SinglePartition` mode, and **no Key** is provided for each message.

### Hashing scheme

Standard hashing functions available when choosing the partition to use for a particular message.

- JavaStringHash - The default hashing function
- Murmur3_32Hash

> [!WARNING]  
> `JavaStringHash` is not useful when producers can be from different multiple language clients,
> under this use case, `Murmur3_32Hash` is recommended

## Non-persistent topics

Non-persistent topics are Pulsar topics in which message data is **never persistently stored to disk 
and kept only in memory.**

Killing the broker or disconnecting a subscriber to a topic means that all messages go to hell.

```
{non-persistent}://tenant/namespace/topic
```

Brokers don't persist messages and immediately send acks back to the producer 
as soon as that message is delivered to connected brokers, so **non-persistent messaging is usually 
faster than persistent messaging.**
