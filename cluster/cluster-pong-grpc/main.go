package main

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/cluster/consul"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/oklahomer/protoactor-go-sender-example/cluster/messages"
	"log"
	"os"
	"os/signal"
)

type ponger struct {
	cluster.Grain
}

var _ messages.Ponger = (*ponger)(nil)

func (*ponger) Terminate() {
}

func (*ponger) ReceiveDefault(ctx actor.Context) {
}

func (*ponger) SendPing(ping *messages.Ping, ctx cluster.GrainContext) (*messages.Pong, error) {
	// The sender process is not a sending actor, but a future process
	log.Printf("Sender: %+v", ctx.Sender())

	pong := &messages.Pong{
		Cnt: ping.Cnt,
	}
	return pong, nil
}

func main() {
	// Setup actor system
	system := actor.NewActorSystem()

	messages.PongerFactory(func() messages.Ponger {
		return &ponger{}
	})

	// Prepare remote env that listens to 8080
	remoteConfig := remote.Configure("127.0.0.1", 8080)

	// Configure cluster on top of the above remote env
	cp, err := consul.New()
	if err != nil {
		log.Fatal(err)
	}
	clusterKind := cluster.NewKind(
		"Ponger",
		actor.PropsFromProducer(func() actor.Actor {
			return &messages.PongerActor{}
		}))
	clusterConfig := cluster.Configure("cluster-grpc-example", cp, remoteConfig, clusterKind)
	c := cluster.New(system, clusterConfig)
	c.Start()

	// Run till signal comes
	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, os.Kill)
	<-finish
}
