package main

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/lmittmann/tint"

	"log/slog"
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
		// ctx.Send() do not set the sender information and hence the recipient can not respond to the sender;
		// ctx.Request() sets the sender information so the receiving actor can refer to the sender's PID and respond.
		ctx.Request(p.pongPid, &ping{})

	case *pong:
		slog.Info("Received a pong message")

	}
}

func main() {
	// Set up a logger to observe the behavior
	logger := slog.New(tint.NewHandler(
		os.Stdout,
		&tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.TimeOnly,
		},
	))
	slog.SetDefault(logger)

	// Set up the actor system
	system := actor.NewActorSystem(
		actor.WithLoggerFactory(func(system *actor.ActorSystem) *slog.Logger {
			return logger.With("system", system.ID)
		}),
	)

	// Run a pong actor that receives a ping payload and sends back a pong payload
	pongProps := actor.PropsFromFunc(func(ctx actor.Context) {
		switch ctx.Message().(type) {
		case *ping:
			slog.Info("Received a ping message")
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
			slog.Info("Finish")
			return

		}
	}
}
