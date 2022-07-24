package main

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
	"log"
	"os"
	"os/signal"
	"protoactor-go-sender-example/remote/messages"
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

		// Below do not set ctx.Self() as the sender, and hence the recipient has no knowledge of the sender
		// even though the message is sent from one actor to another.
		ctx.Send(p.pongPid, ping)

	case *messages.Pong:
		// Never comes here.
		// The recipient can not refer to the sender.
		log.Print("Received pong message")

	}
}

func main() {
	// Set up actor system
	system := actor.NewActorSystem()

	// Set up a remote env that listens to a randomly chosen port.
	// To specify a port number, pass a desired port number as a second argument of remote.Configure instead of 0.
	// Note that a ponger implementation at ../remote-pong/main.go specifies its actor system's port as 8080 so the pinger can refer to it.
	config := remote.Configure("127.0.0.1", 0)
	remoting := remote.NewRemote(system, config)
	remoting.Start()

	// Declare the remote pong actor's address, and let the ping actor send a ping payload to it
	remotePong := actor.NewPID("127.0.0.1:8080", "pongActorID")
	pingProps := actor.PropsFromProducer(func() actor.Actor {
		return &pingActor{
			pongPid: remotePong,
		}
	})
	pingPid := system.Root.Spawn(pingProps)

	// Subscribe to a signal to finish the interaction
	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, os.Kill)

	// Periodically send a ping payload till signal comes
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
