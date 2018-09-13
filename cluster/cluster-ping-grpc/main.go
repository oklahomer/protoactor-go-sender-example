package main

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/cluster/consul"
	"github.com/oklahomer/protoactor-go-sender-example/cluster/messages"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var cnt uint64 = 0

type pingActor struct {
	cnt uint
}

func (p *pingActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case struct{}:
		cnt += 1
		ping := &messages.Ping{
			Cnt: cnt,
		}

		grain := messages.GetPongerGrain("ponger-1")
		pong, err := grain.SendPing(ping)
		if err != nil {
			log.Print(err.Error())
			return
		}
		log.Printf("Received %#v", pong)

	case *messages.Pong:
		// Never comes here.
		// When the pong grain responds to the sender's gRPC call,
		// the sender is not a ping actor but a future process.
		log.Print("Received pong message")

	}
}

func main() {
	cp, err := consul.New()
	if err != nil {
		log.Fatal(err)
	}
	cluster.Start("cluster-grpc-example", "127.0.0.1:8081", cp)

	pingProps := actor.FromProducer(func() actor.Actor {
		return &pingActor{}
	})
	pingPid := actor.Spawn(pingProps)

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt)
	signal.Notify(finish, syscall.SIGTERM)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pingPid.Tell(struct{}{})

		case <-finish:
			return
			log.Print("Finish")

		}
	}
}
