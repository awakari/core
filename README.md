# Contents

1. [Overview](#1-overview)<br/>
2. [Configuration](#2-configuration)<br/>
3. [Deployment](#3-deployment)<br/>
4. [Usage](#4-usage)<br/>
5. [Design](#5-design)<br/>
6. [Contributing](#6-contributing)<br/>
   6.1. [Versioning](#61-versioning)<br/>
   6.2. [Issue Reporting](#62-issue-reporting)<br/>
   6.3. [Building](#63-building)<br/>
   6.4. [Testing](#64-testing)<br/>
   &nbsp;&nbsp;&nbsp;6.4.1. [Functional](#641-functional)<br/>
   &nbsp;&nbsp;&nbsp;6.4.2. [Performance](#642-performance)<br/>
   6.5. [Releasing](#65-releasing)<br/>


# 1. Overview

This repo contains the Helm chart for the Awakari Core system deployment.
The chart deploys the following minimum set of components:
1. Sharded MongoDB (if external MongoDB option is not used)
2. NATS in the Jetstream mode (if NATS usage is enabled).
3. Specific queue wrapper service (e.g. [queue-nats](https://github.com/awakari/queue-nats))
4. Specific distributed semaphore service (e.g. [semaphore-nats](https://github.com/awakari/semaphore-nats))
5. Specific condition services (e.g. [kiwi-tree](https://github.com/awakari/kiwi-tree))
6. [Subscriptions](https://github.com/awakari/subscriptions) service
7. [Matches](https://github.com/awakari/matches) service
8. [Messages](https://github.com/awakari/messages) service
9. [Reader](https://github.com/awakari/reader) service
10. [Writer](https://github.com/awakari/writer) service

# 2. Configuration

For a component-specific options see the corresponding sub-chart configuration. Here follow own configuration options: 

| Variable               | Default | Description                                                                                                              |
|------------------------|---------|--------------------------------------------------------------------------------------------------------------------------|
| conditions.kiwi.tree   | `true`  | Enables the kiwi-tree conditions usage. May be used together with other conditions implementations.                      | 
| queue.backend.nats     | `true`  | Enables the NATS JetStream queue wrapper service. Exclusive, can not be used together with other queue backends.         |
| semaphore.backend.nats | `true`  | Enables the NATS-based distributed semaphore service. Exclusive, can not be used together with other semaphore backends. |

# 3. Deployment

Create the target namespace:
```shell
kubectl create namespace awakari
```

> **Note**
> 
> To use external MongoDB, use the values file [values-mongodb-ext.yaml](helm/core/values-mongodb-ext.yaml) for the
> reference and substitute these with own values.

Install the package built locally:
```shell
helm install core core-0.0.0.tgz -n awakari
```

> **Warning**
> 
> Do not change the "core" release name

Or use an existing:
```shell
helm repo add awakari-core https://awakari.github.io/core

helm install core awakari-core/core \
  -n awakari
```

# 4. Usage

## 4.1. Client SDK

Refer to [Client SDK Usage](https://github.com/awakari/client-sdk-go#3-usage).

> **Note**:
>
> Usage Limits and Permits APIs are not available in the Core.

## 4.2. API

* [Subscriptions](https://github.com/awakari/subscriptions#4-usage)
* [Reader](https://github.com/awakari/reader#4-usage)
* [Writer](https://github.com/awakari/writer#4-usage)

> **Note**:
>
> Invoking Reader and Writer APIs without Client SDK requires the corresponding proto files to be available on the client side.


# 5. Design

The core of Awakari consist of:
* Storages
  * [Subscriptions](https://github.com/awakari/subscriptions)
  * Conditions, e.g. [Kiwi Tree](https://github.com/awakari/kiwi-tree)
  * [Matches](https://github.com/awakari/matches)
  * [Messages](https://github.com/awakari/messages)
* Stateless components
  * [Writer](https://github.com/awakari/writer)
  * [Reader](https://github.com/awakari/reader)


![components](doc/components-core.png)

# 6. Contributing

## 6.1. Versioning

The service uses the [semantic versioning](http://semver.org/).
The single source of the version info is the git tag:
```shell
git describe --tags --abbrev=0
```

## 6.2. Issue Reporting

TODO

## 6.3. Building

Build a helm package:
```shell
helm package helm/core
```

## 6.4. Testing

### 6.4.1. Functional

The repo contains core functional end-to-end tests.

To run the tests locally:

1. Port-forward the reader API to local port 50051
2. Port-forward the subscriptions API to local port 50052
3. Port-forward the writer API to local port 50053
4. 
```shell 
make test
```

To run the tests in K8s cluster:
```shell
helm test core -n awakari --filter name=core-test
```

### 6.4.2. Performance

TODO

## 6.5. Releasing

To release a new version (e.g. `1.2.3`) it's enough to put a git tag:
```shell
git tag -v1.2.3
git push --tags
```

The corresponding CI job is started to build a helm chart and publish it with the specified tag (+latest).
