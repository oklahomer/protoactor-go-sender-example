package main

import (
	"log/slog"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/lmittmann/tint"

	"os"
	"os/signal"
	"protoactor-go-sender-example/remote/messages"
)

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

	// Set up a remote env that listens to 8080
	config := remote.Configure("127.0.0.1", 8080)
	remoter := remote.NewRemote(system, config)
	remoter.Start()
	defer remoter.Shutdown(false)

	// Run pong actor that receives ping payload, and then send back pong payload
	pongProps := actor.PropsFromFunc(func(ctx actor.Context) {
		switch msg := ctx.Message().(type) {
		case *messages.Ping:
			pong := &messages.Pong{Cnt: msg.Cnt}
			slog.Info("Received a ping message", "message", msg)
			ctx.Respond(pong)

		default:

		}
	})
	pongPid, err := system.Root.SpawnNamed(pongProps, "pongActorID")
	if err != nil {
		slog.Error("Failed to spawn an actor", "error", err)
		os.Exit(1)
	}
	slog.Info("An actor is running", "address", pongPid.Address, "id", pongPid.Id)

	// Run till a signal comes
	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, os.Kill)
	<-finish
}
