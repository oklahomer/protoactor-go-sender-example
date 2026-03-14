package main

import (
	"log/slog"
	"sync/atomic"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/clusterproviders/consul"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/lmittmann/tint"

	"os"
	"os/signal"
	"protoactor-go-sender-example/cluster/messages"
	"time"
)

var counter atomic.Uint64

type pingActor struct {
	system *actor.ActorSystem
	cnt    uint
}

func (p *pingActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case struct{}:
		ping := &messages.PingMessage{
			Cnt: counter.Add(1),
		}

		grainPid := cluster.GetCluster(p.system).Get("ponger-1", "Ponger")
		// Below do not set ctx.Self() as the sender, and hence the recipient has no knowledge of the sender
		// even though the message is sent from one actor to another.
		//
		ctx.Send(grainPid, ping)

	case *messages.PongMessage:
		// Never comes here because the recipient can not refer to the sender.
		// Instead, the cluster grain leaves logs as below:
		// hh:mm:ss INF Received ping message
		// hh:mm:ss INF [DeadLetter] system=NeyzJ6ZJ3DV7GNMFYz7R8H pid=<nil> message=cnt:84 sender=<nil>
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

	// Prepare a remote env that listens to 8081
	remoteConfig := remote.Configure("127.0.0.1", 8081)

	// Configure a cluster on top of the above remote env
	cp, err := consul.New()
	if err != nil {
		slog.Error("Failed to create a consul provider", "error", err)
		os.Exit(1)
	}
	lookup := disthash.New()
	clusterConfig := cluster.Configure("cluster-example", cp, lookup, remoteConfig)
	c := cluster.New(system, clusterConfig)

	// Manage the cluster client's lifecycle
	c.StartClient() // Configure as a client
	defer c.Shutdown(false)

	// Start a ping actor that periodically sends a "ping" payload to the "Ponger" cluster grain
	pingProps := actor.PropsFromProducer(func() actor.Actor {
		return &pingActor{
			system: system,
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
