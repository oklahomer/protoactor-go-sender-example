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

		// Below do not set ctx.Self() as sender,
		// and hence the recipient has no knowledge of the sender
		// even though the message is sent from one actor to another.
		//
		ctx.Send(p.pongPid, ping)

	case *messages.Pong:
		// Never comes here.
		// The recipient can not refer to the sender.
		log.Print("Received pong message")

	}
}

func main() {
	system := actor.NewActorSystem()

	config := remote.Configure("127.0.0.1", 8081)
	remoting := remote.NewRemote(system, config)
	remoting.Start()

	remotePong := actor.NewPID("127.0.0.1:8080", "pongActorID")
	pingProps := actor.PropsFromProducer(func() actor.Actor {
		return &pingActor{
			pongPid: remotePong,
		}
	})
	pingPid := system.Root.Spawn(pingProps)

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, syscall.SIGINT)
	signal.Notify(finish, syscall.SIGTERM)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			system.Root.Send(pingPid, struct{}{})

		case <-finish:
			log.Print("Finish")
			remoting.Shutdown(false)
			return

		}
	}
}
