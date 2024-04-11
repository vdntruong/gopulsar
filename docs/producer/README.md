# Producer

## Message

The default max size of a message is 5 MB

We can configure the max size of a message with the `broker.conf` file
```
maxMessageSize=5242880
```

or `bookkeeper.conf` files

```
nettyMaxFrameSizeBytes=5253120
```

### Binary protocol

Pulsar uses **a custom binary protocol for communications** between producers/consumers and brokers.

Clients and brokers exchange commands with each other. 
Commands are formatted as **binary protocol buffer (aka protobuf)** messages.

[Refer article](https://pulsar.apache.org/docs/3.2.x/developing-binary-protocol)
