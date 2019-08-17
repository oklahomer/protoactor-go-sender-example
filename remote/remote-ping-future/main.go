package main

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/oklahomer/protoactor-go-sender-example/remote/messages"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var cnt uint64 = 0

type pingActor struct {
	cnt     uint
	pongPid *actor.PID
}

func (p *pingActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case struct{}:
		cnt += 1
		ping := &messages.Ping{
			Cnt: cnt,
		}

		// Below both work.
		//
		//future := p.pongPid.RequestFuture(ping, time.Second)
		future := ctx.RequestFuture(p.pongPid, ping, time.Second)
		result, err := future.Result()
		if err != nil {
			log.Print(err.Error())
			return
		}
		log.Printf("Received %v", result)

	case *messages.Pong:
		// Never comes here.
		// When the pong actor responds to the sender,
		// the sender is not a ping actor but a future process.
		log.Print("Received pong message")

	}
}

func main() {
	remote.Start("127.0.0.1:8081")

	rootContext := actor.EmptyRootContext

	remotePong := actor.NewPID("127.0.0.1:8080", "pongActorID")
	pingProps := actor.PropsFromProducer(func() actor.Actor {
		return &pingActor{
			pongPid: remotePong,
		}
	})
	pingPid := rootContext.Spawn(pingProps)

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, syscall.SIGINT)
	signal.Notify(finish, syscall.SIGTERM)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rootContext.Send(pingPid, struct{}{})

		case <-finish:
			log.Print("Finish")
			return

		}
	}
}
