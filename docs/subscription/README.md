# Subscription

A Pulsar subscription is **a named configuration rule** that determines
**how messages are delivered to consumers.**

> [!TIP]  
> Pub-Sub or Queuing In Pulsar, you can use different subscriptions flexibly.
> 
> - traditional "fan-out pub-sub messaging", specify a **unique subscription** 
> name for each consumer. It is an **exclusive** subscription type
> -  "message queuing", share the **same subscription name** among multiple consumers 
> (**shared**, **failover**, **key_shared**)
> - for both effects simultaneously, combine **exclusive** subscription types with other 
> subscription types for consumers

## Subscription modes

The subscription mode indicates the cursor belongs to **durable** type or **non-durable** type.

When subscription is created, an associated cursor is created to record the last consumed position.

### Durable

When a broker restarts from a failure, it can recover the cursor from 
the persistent storage **(BookKeeper)**

> [!NOTE]  
> Durable is the default subscription mode
### NonDurable

Once a broker stops, the cursor is lost and can never be recovered

> [!NOTE]  
> Reader's subscription mode is NonDurable, and it does not prevent data in a topic
> from being deleted. 
> 
> Reader's subscription mode can **not** be changed

## Subscription type

The subscription type determines which messages go to which consumers.

### Exclusive

> [!NOTE]  
> Exclusive is the default subscription type

Allows **only a single consumer** to attach to the subscription. 

> [!CAUTION]  
> If another consumer subscribes to the topic, an error occurs

### Failover

Allows multiple consumers attach to the same subscription.

A master consumer is picked and receives messages. **When the master consumer disconnects**, 
all (non-acknowledged and subsequent) messages are delivered to **the next consumer in line**.

> [!WARNING]  
> In some cases, when the switchover operation is performed, it could result in messages **being 
> duplicated or received out of order**

#### Failover | Partitioned topics

#### Failover | Non-partitioned topics

### Shared

Allows multiple consumers attach to the same subscription.

Messages are **delivered in a round-robin distribution** across consumers, 
and any given message is delivered to only one consumer, similar to worker pool concept.

> [!IMPORTANT]  
> Shared subscriptions **do not guarantee** message ordering or support cumulative acknowledgment

### Key_shared

Allows multiple consumers attach to the same subscription.

Messages are delivered in distribution across consumers and **messages with the same key or 
same ordering key are delivered to only one consumer.**

```
Message Key --> /Hash Function/ -- hash (32-bit) --> /Algorithm/ --> Consumer
```

The mapping algorithms to select a consumer for a given message key:

- Auto-split Hash Range
- Auto-split Consistent Hashing
- Sticky

> [!TIP]   
> **Throughput** concerns should guide **partitioning/routing** decisions
> while **subscription** decisions should be guided by **application semantics**
