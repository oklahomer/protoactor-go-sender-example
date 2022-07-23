This repository contains supplemental examples for my blog article, [[Golang] Protoactor-go 101: How actors communicate with each other](https://blog.oklahome.net/2018/09/protoactor-go-messaging-protocol.html) , to cover all message passing methods for all kinds of actors provided by protoactor-go.

![](https://raw.githubusercontent.com/oklahomer/protoactor-go-sender-example/master/docs/components.png)

# Local
For local message passing, see the below directories:
- [local-send](https://github.com/oklahomer/protoactor-go-sender-example/blob/master/local-send/main.go) ... Use Send() for local message passing. The recipient actor cannot refer to the sender actor.
- [local-request](https://github.com/oklahomer/protoactor-go-sender-example/blob/master/local-request/main.go) ... Use Request() for local message passing. The recipient actor can refer to the sender actor.
- [local-future](https://github.com/oklahomer/protoactor-go-sender-example/tree/master/local-future) ... Use RequestFuture() for local message passing. Context.Sender() does not return the PID of sender actor but that of actor.Future.

# Remote
- remote/messages ... Contain Protobuf serializable message structures.
- rmeote/remote-pong ... A process that returns pong message to sender.
- remote/remote-ping-send ... A process that sends message to pong actor by Send(). The recipient cannot refer to the sender actor.
- remote/remote-ping-request ... A process that sends message to pong actor by Request(). The recipient actor can refer to the sender actor.
- remote/remote-ping-future ... A process that sends message to pong actor by RequestFuture(). Context.Sender() does not return the PID of sender actor but that of actor.Future.

# Cluster Grain
- cluster/messages ... Contain Protobuf serializable message structures and generated actor.Actor implementation for gRPC based communication.

## Cluster Grain usage with remote communication
Below implementation uses Consul Cluster Provider for service discovery
- cluster-consul/cluster-pong ... A process that returns pong message to the sender based on remote actor implementation.
- cluster-consul/cluster-ping-send ... A process that sends message to pong actor by Send(). The recipient cannot refer to the sender actor.
- cluster-consul/cluster-ping-request ... A process that sends message to pong actor by Request(). The recipient actor can refer to the sender actor.
- cluster-consul/cluster-ping-future ... A process that sends message to pong actor by RequestFuture(). Context.Sender() does not return the PID of sender actor but that of actor.Future.

## Cluster Grain usage with gRPC based communication
Below implementation uses Consul Cluster Provider for service discovery
- cluster-consul/cluster-pong-grpc ... A process that returns pong message to the sender via gRPC service.
- cluster-consul/cluster-ping-grpc ... A process that sends message to pong actor over gRPC based service.

## Cluster with Automanaged Cluster Provider
Below implementation uses Automanaged Cluster Provider for service discovery
- cluster-automanaged/cluster-pong ... A process that returns pong message to the sender based on remote actor implementation.
- cluster-automanaged/cluster-ping-future ... A process that sends message to pong actor by Request(). The recipient actor can refer to the sender actor.

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
