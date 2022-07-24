This repository contains supplemental examples for my blog article, [[Golang] Protoactor-go 101: How actors communicate with each other](https://blog.oklahome.net/2018/09/protoactor-go-messaging-protocol.html) , to cover all message passing methods for all kinds of actors provided by protoactor-go.

![](https://raw.githubusercontent.com/oklahomer/protoactor-go-sender-example/master/docs/components.png)

# Local
For local message passing, see the below directories:
- [local-send](./local-send/main.go) ... Use Send() for local message passing. The recipient actor cannot refer to the sender actor.
- [local-request](./local-request/main.go) ... Use Request() for local message passing. The recipient actor can refer to the sender actor.
- [local-future](./local-future/main.go) ... Use RequestFuture() for local message passing. Context.Sender() does not return the PID of sender actor but that of actor.Future.

# Remote
- [remote/messages](./remote/messages/) ... Contain Protobuf serializable message structures. See **README.md** for code generation.
- [rmeote/remote-pong](./remote/remote-pong/main.go) ... A process that returns pong message to sender.
- [remote/remote-ping-send](./remote/remote-ping-send/main.go) ... A process that sends message to pong actor by Send(). The recipient cannot refer to the sender actor.
- [remote/remote-ping-request](./remote/remote-ping-request/main.go) ... A process that sends message to pong actor by Request(). The recipient actor can refer to the sender actor.
- [remote/remote-ping-future](./remote/remote-ping-future/main.go) ... A process that sends message to pong actor by RequestFuture(). Context.Sender() does not return the PID of sender actor but that of actor.Future.

# Cluster Grain
- [cluster/messages](./cluster/messages/) ... Contain Protobuf serializable message structures and generated actor.Actor implementation for gRPC based communication. See **README.md** for code generation.

## Cluster Grain usage with remote communication
Below examples use Consul Cluster Provider for service discovery. Run `docker-compose -f docker-compose.yml up --build -d` or something equivalent to run Consul on your local environment.
- [cluster-consul/cluster-pong](./cluster-consul/cluster-pong/main.go) ... A process that returns pong message to the sender based on remote actor implementation.
- [cluster-consul/cluster-ping-send](./cluster-consul/cluster-ping-send/main.go) ... A process that sends message to pong actor by Send(). The recipient cannot refer to the sender actor.
- [cluster-consul/cluster-ping-request](./cluster-consul/cluster-ping-request/main.go) ... A process that sends message to pong actor by Request(). The recipient actor can refer to the sender actor.
- [cluster-consul/cluster-ping-future](./cluster-consul/cluster-ping-future/main.go) ... A process that sends message to pong actor by RequestFuture(). Context.Sender() does not return the PID of sender actor but that of actor.Future.

## Cluster Grain usage with gRPC based communication
Below examples use Consul Cluster Provider for service discovery. Run `docker-compose -f docker-compose.yml up --build -d` or something equivalent to run Consul on your local environment.
- [cluster-consul/cluster-pong-grpc](./cluster-consul/cluster-pong-grpc/main.go) ... A process that returns pong message to the sender via gRPC service.
- [cluster-consul/cluster-ping-grpc](./cluster-consul/cluster-ping-grpc/main.go) ... A process that sends message to pong actor over gRPC based service.

## Cluster with Automanaged Cluster Provider
Below examples use Automanaged Cluster Provider for service discovery
- [cluster-automanaged/cluster-pong](./cluster-automanaged/cluster-pong/main.go) ... A process that returns pong message to the sender based on remote actor implementation.
- [cluster-automanaged/cluster-ping-future](./cluster-automanaged/cluster-ping-future/main.go) ... A process that sends message to pong actor by Request(). The recipient actor can refer to the sender actor and therefore can successfully respond.

## Cluster with etcd Cluster Provider
Below examples uses etcd Cluster Provider for service discovery. Run `docker-compose -f docker-compose.yml up --build -d` or something equivalent to run etcd on your local environment.
- [cluster-etcd/cluster-pong](./cluster-etcd/cluster-pong/main.go) ... A process that returns pong message to the sender based on remote actor implementation.
- [cluster-etcd/cluster-ping-future](./cluster-etcd/cluster-ping-future/main.go) ... A process that sends message to pong actor by Request(). The recipient actor can refer to the sender actor and therefore can successfully respond.

# References
- [[Golang] Protoactor-go 101: Introduction to golang's actor model implementation](https://blog.oklahome.net/2018/07/protoactor-go-introduction.html)
- [[Golang] Protoactor-go 101: How actors communicate with each other](https://blog.oklahome.net/2018/09/protoactor-go-messaging-protocol.html)
- [[Golang] protoactor-go 101: How actor.Future works to synchronize concurrent task execution](https://blog.oklahome.net/2018/11/protoactor-go-how-future-works.html)
- [[Golang] protoactor-go 201: How middleware works to intercept incoming and outgoing messages](https://blog.oklahome.net/2018/11/protoactor-go-middleware.html)
- [[Golang] protoactor-go 201: Use plugins to add behaviors to an actor](https://blog.oklahome.net/2018/12/protoactor-go-use-plugin-to-add-behavior.html)
- [[Golang] protoactor-go 301: How proto.actor's clustering works to achieve higher availability](https://blog.oklahome.net/2021/05/protoactor-clustering.html)

# Other Example Codes
- [oklahomer/protoactor-go-future-example](https://github.com/oklahomer/protoactor-go-future-example)
  - Some example codes to illustrate how protoactor-go handles Future process
- [oklahomer/protoactor-go-middleware-example](https://github.com/oklahomer/protoactor-go-middleware-example)
  - Some example codes to illustrate how protoactor-go use Middleware
