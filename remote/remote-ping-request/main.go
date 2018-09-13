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
		//p.pongPid.Request(ping, ctx.Self())
		ctx.Request(p.pongPid, ping)

	case *messages.Pong:
		log.Print("Received pong message")

	}
}

func main() {
	remote.Start("127.0.0.1:8081")

	remotePong := actor.NewPID("127.0.0.1:8080", "pongActorID")

	pingProps := actor.FromProducer(func() actor.Actor {
		return &pingActor{
			pongPid: remotePong,
		}
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
