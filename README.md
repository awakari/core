# 1. Overview

Awakari Core is a system to be used with custom 3rd-party components:
1. Custom message producer (e.g. [producer-rss](https://github.com/awakari/producer-rss))
2. Custom message consumer (e.g. [consumer-log](https://github.com/awakari/consumer-log))
3. Custom queue implementation (e.g. NATS JetStream, Kafka, AWS SQS, etc)

This repo contains the Helm chart for the Awakari Core system deployment.
The chart deploys the following minimum set of components:
1. Sharded MongoDB
2. Specific queue wrapper service (e.g. [queue-nats](https://github.com/awakari/queue-nats))
3. Specific condition services (e.g. [kiwi-tree](https://github.com/awakari/kiwi-tree))
4. [Subscriptions](https://github.com/awakari/subscriptions) service
5. [Matches](https://github.com/awakari/matches) service
6. [Resolver](https://github.com/awakari/resolver) service
7. [Router](https://github.com/awakari/router) service

# 2. Configuration

For a component-specific options see the corresponding sub-chart configuration. Here follow own configuration options: 

| Variable             | Default | Description                                                                                                      |
|----------------------|---------|------------------------------------------------------------------------------------------------------------------|
| conditions.kiwi.tree | `true`  | Enables the kiwi-tree conditions usage. May be used together with other conditions implementations.              | 
| queue.backend.nats   | `true`  | Enables the NATS JetStream queue wrapper service. Exclusive, can not be used together with other queue backends. |

# 3. Deployment

## 3.1. Dependencies

### 3.1.1. Queue

It's necessary to choose a queue backend implementation to be used by Awakari Core internally.

#### 3.1.1.1. NATS JetStream

NATS JetStream is the default queue option used by Awakari Core. The following example command may be used to deploy the 
NATS JetStream service if there's no existing one to reuse: 

```shell
helm install nats bitnami/nats \
  -n awakari \
  --set jetstream.enabled=true \
  --set persistence.enabled=true \
  --set resourceType="statefulset" \
  --set auth.user=awakari \
  --set auth.password=awakari \
  --set replicaCount=3
```

### 3.1.2. Consumer

The system sends the processed messages to a consumer. By default, it expects a 
[consumer gRPC service](https://github.com/awakari/consumer-log/blob/master/api/grpc/service.proto) 
to be present in the same K8s namespace with the URI "consumer:8080".

For the testing purposes the [dummy implementation](https://github.com/awakari/consumer-log) may be deployed and used. 

## 3.2. Self

Build a helm package:
```shell
helm package .
```

Install the package:
```shell
helm install core core-0.0.0.tgz -n awakari
```

# 4. Usage

## 4.1. Submit Messages

The resolver service is the entrypoint to feed the messages to the core system. See the 
[message send example](https://github.com/awakari/resolver#4-usage).

## 4.2. Connect Producer

A custom producer may be used to feed the messages to the system. For the demo purpose the 
[producer-rss](https://github.com/producer-rss) may be deployed and connected to the core system.
