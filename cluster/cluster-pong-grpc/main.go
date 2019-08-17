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
	"syscall"
)

type ponger struct {
	cluster.Grain
}

var _ messages.Ponger = (*ponger)(nil)

func (*ponger) Terminate() {
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
	messages.PongerFactory(func() messages.Ponger {
		return &ponger{}
	})

	// kind name must be the same one as protos.proto's service name
	// This name is set in PongerGrain.SendPingWithOpts
	remote.Register("Ponger", actor.PropsFromProducer(func() actor.Actor {
		return &messages.PongerActor{}
	}))

	cp, err := consul.New()
	if err != nil {
		log.Fatal(err)
	}
	cluster.Start("cluster-grpc-example", "127.0.0.1:8080", cp)

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, syscall.SIGINT)
	signal.Notify(finish, syscall.SIGTERM)
	<-finish
	log.Print("Finish")
}
