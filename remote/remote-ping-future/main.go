package main

import (
	"log/slog"
	"sync/atomic"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/lmittmann/tint"

	"os"
	"os/signal"
	"protoactor-go-sender-example/remote/messages"
	"time"
)

var counter atomic.Uint64

type pingActor struct {
	cnt     uint
	pongPid *actor.PID
}

func (p *pingActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case struct{}:
		ping := &messages.Ping{
			Cnt: counter.Add(1),
		}

		// Besides the fact that a payload is serialized and is sent to a remove environment,
		// the usage of ctx.RequestFuture is still the same as that of local messaging.
		future := ctx.RequestFuture(p.pongPid, ping, time.Second)
		result, err := future.Result()
		if err != nil {
			slog.Error("Failed to receive a result", "error", err)
			return
		}
		slog.Info("Received a message", "message", result)

	case *messages.Pong:
		// Never comes here.
		// When the pong actor responds to the sender,
		// the sender is not a ping actor but a future process.
		slog.Info("Received a pong message", "message", msg)

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

	// Set up actor system
	system := actor.NewActorSystem(
		actor.WithLoggerFactory(func(system *actor.ActorSystem) *slog.Logger {
			return logger.With("system", system.ID)
		}),
	)

	// Set up a remote env that listens to a randomly chosen port.
	// To specify a port number, pass a desired port number as a second argument of remote.Configure instead of 0.
	// Note that a ponger implementation at ../remote-pong/main.go specifies its actor system's port as 8080 so the pinger can refer to it.
	config := remote.Configure("127.0.0.1", 0)
	remoter := remote.NewRemote(system, config)
	remoter.Start()
	defer remoter.Shutdown(false)

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
			slog.Info("Finish")
			return

		}
	}
}
