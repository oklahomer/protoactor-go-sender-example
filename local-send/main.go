package main

import (
	"github.com/asynkron/protoactor-go/actor"
	"log"
	"os"
	"os/signal"
	"time"
)

type pong struct {
}

type ping struct {
}

type pingActor struct {
	pongPid *actor.PID
}

func (p *pingActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case struct{}:
		// Below do not set ctx.Self() as sender,
		// and hence the recipient has no knowledge of the sender
		// even though the message is sent from one actor to another.
		//
		ctx.Send(p.pongPid, &ping{})

	case *pong:
		// Never comes here.
		// When the pong actor tries to respond, the sender is not set.
		log.Print("Received pong message")

	}
}

func main() {
	// Setup actor system
	system := actor.NewActorSystem()

	// Run pong actor that receives ping payload and send back pong payload
	pongProps := actor.PropsFromFunc(func(ctx actor.Context) {
		switch ctx.Message().(type) {
		case *ping:
			log.Print("Received ping message")
			ctx.Respond(&pong{})

		default:

		}
	})
	pongPid := system.Root.Spawn(pongProps)

	// Run ping actor that receives an arbitrary payload from outside of actor system, and then send ping payload to pong actor
	pingProps := actor.PropsFromProducer(func() actor.Actor {
		return &pingActor{
			pongPid: pongPid,
		}
	})
	pingPid := system.Root.Spawn(pingProps)

	// Subscribe to signal to finish interaction
	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, os.Kill)

	// Periodically send ping payload till signal comes
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			system.Root.Send(pingPid, struct{}{})

		case <-finish:
			log.Print("Finish")
			return

		}
	}
}
