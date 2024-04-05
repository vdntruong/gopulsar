# Go Pulsar

## Setup

### Install Go package

```bash
go mod tidy && go mod download
```

### Start a standalone Pulsar cluster in Docker 

```bash
docker run -it -d \
-e PULSAR_STANDALONE_USE_ZOOKEEPER=1 \
-p 6650:6650 \
-p 8080:8080 \
--mount source=pulsardata,target=/pulsar/data \
--mount source=pulsarconf,target=/pulsar/conf \
apachepulsar/pulsar:3.2.2 \
bin/pulsar standalone
```
