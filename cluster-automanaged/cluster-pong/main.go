package main

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/cluster/automanaged"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/oklahomer/protoactor-go-sender-example/cluster/messages"
	"log"
	"os"
	"os/signal"
	"time"
)

type ponger struct {
}

func (*ponger) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *messages.Ping:
		pong := &messages.Pong{Cnt: msg.Cnt}
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
	cp := automanaged.NewWithConfig(1*time.Second, 6331, "localhost:6330", "localhost:6331")
	clusterConfig := cluster.Configure("cluster-example", cp, remoteConfig, clusterKind)
	c := cluster.New(system, clusterConfig)
	c.Start()

	// Run till signal comes
	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, os.Kill)
	<-finish
}
