# 1. Overview

This repo contains the Helm chart for the Awakari Core system deployment.
The chart deploys the following minimum set of components:
1. Sharded MongoDB
2. NATS in the Jetstream mode (if NATS usage is enabled).
3. Specific queue wrapper service (e.g. [queue-nats](https://github.com/awakari/queue-nats))
4. Specific condition services (e.g. [kiwi-tree](https://github.com/awakari/kiwi-tree))
5. [Subscriptions](https://github.com/awakari/subscriptions) service
6. [Matches](https://github.com/awakari/matches) service
7. [Messages](https://github.com/awakari/messages) service
8. [Writer](https://github.com/awakari/writer) service

# 2. Configuration

For a component-specific options see the corresponding sub-chart configuration. Here follow own configuration options: 

| Variable             | Default | Description                                                                                                      |
|----------------------|---------|------------------------------------------------------------------------------------------------------------------|
| conditions.kiwi.tree | `true`  | Enables the kiwi-tree conditions usage. May be used together with other conditions implementations.              | 
| queue.backend.nats   | `true`  | Enables the NATS JetStream queue wrapper service. Exclusive, can not be used together with other queue backends. |

# 3. Deployment

Build a helm package:
```shell
helm package helm/core
```

Create the target namespace:
```shell
kubectl create namespace awakari
```

Install the package built locally:
```shell
helm install core core-0.0.0.tgz -n awakari
```

Or use an existing:
```shell
helm repo add awakari-core https://awakari.github.io/helm-core

helm install core awakari-core/core \
  -n awakari
```

# 4. Usage

TODO

# 5. Testing

TODO