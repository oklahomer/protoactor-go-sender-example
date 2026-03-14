package main

import (
	"log/slog"
	"sync/atomic"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lmittmann/tint"

	"os"
	"os/signal"
	"time"
)

var counter atomic.Int64

type pong struct {
	PingMessageNum int64
}

type ping struct {
	MessageNum int64
}

type pingActor struct {
	pongPid *actor.PID
}

func (p *pingActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case struct{}:
		// Unlike ctx.Request(), this sends a request and waits for a response,
		// expecting the receiving actor to respond with ctx.Respond().
		// The call to ctx.RequestFuture() itself is not blocking because its returning value is Future.
		// However, Future.Result() blocks until the receiving actor responds
		// or the timeout interval specified by the call to ctx.RequestFuture() passes.
		future := ctx.RequestFuture(p.pongPid, &ping{MessageNum: counter.Add(1)}, time.Second)
		result, err := future.Result()
		if err != nil {
			slog.Error("Failed to receive a result", "error", err)
			return
		}
		slog.Info("Received a pong message", "message", result)

	case *pong:
		// Never comes here.
		// When the pong actor responds to the sender,
		// the sender is not a ping actor but a future process.
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
		switch msg := ctx.Message().(type) {
		case *ping:
			slog.Info("Received a ping message", "message", msg)
			ctx.Respond(&pong{PingMessageNum: msg.MessageNum})

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
