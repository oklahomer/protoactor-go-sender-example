package main

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/clusterproviders/etcd"
	"github.com/asynkron/protoactor-go/remote"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"os"
	"os/signal"
	"protoactor-go-sender-example/cluster/messages"
	"time"
)

type ponger struct {
}

func (*ponger) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *messages.PingMessage:
		pong := &messages.PongMessage{Cnt: msg.Cnt}
		log.Print("Received ping message")
		ctx.Respond(pong)

	default:

	}
}

func main() {
	// Setup actor system
	system := actor.NewActorSystem()

	// Prepare remote env that listens to 8080
	remoteConfig := remote.Configure("127.0.0.1", 8080)

	// Configure cluster on top of the above remote env
	clusterKind := cluster.NewKind(
		"Ponger",
		actor.PropsFromProducer(func() actor.Actor {
			return &ponger{}
		}))
	// Configure cluster on top of the above remote env
	cp, err := etcd.NewWithConfig("/protoactor", clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		panic(err)
	}
	clusterConfig := cluster.Configure("cluster-example", cp, remoteConfig, clusterKind)
	c := cluster.New(system, clusterConfig)
	c.Start()

	// Run till signal comes
	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, os.Kill)
	<-finish
}
