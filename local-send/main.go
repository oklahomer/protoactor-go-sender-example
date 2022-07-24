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
		// Below do not set ctx.Self() as a sender, and hence the recipient has no knowledge of the sender actor
		// even though the message is sent from one actor to another.
		// Use ctx.Request() or ctx.RequestFuture() when expecting a response.
		ctx.Send(p.pongPid, &ping{})

	case *pong:
		// Never comes here.
		// When the pong actor tries to respond, the sender is not set.
		log.Print("Received pong message")

	}
}

func main() {
	// Set up actor system
	system := actor.NewActorSystem()

	// Run a pong actor that receives a ping payload and sends back a pong payload
	pongProps := actor.PropsFromFunc(func(ctx actor.Context) {
		switch ctx.Message().(type) {
		case *ping:
			log.Print("Received ping message")

			// This call fails and ends up with a dead letter because the sender did not set the sender information.
			ctx.Respond(&pong{})

		default:

		}
	})
	pongPid := system.Root.Spawn(pongProps)

	// Run a ping actor that receives an arbitrary payload from outside the actor system, and then sends a ping payload to the pong actor
	pingProps := actor.PropsFromProducer(func() actor.Actor {
		return &pingActor{
			pongPid: pongPid,
		}
	})
	pingPid := system.Root.Spawn(pingProps)

	// Subscribe to a signal to finish the interaction
	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, os.Kill)

	// Periodically send a ping payload till a signal comes
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
