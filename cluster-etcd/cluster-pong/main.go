package main

import (
	"log/slog"
	"os"
	"os/signal"
	"protoactor-go-sender-example/cluster/messages"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/clusterproviders/etcd"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/lmittmann/tint"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ponger struct {
}

func (*ponger) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *messages.PingMessage:
		pong := &messages.PongMessage{Cnt: msg.Cnt}
		slog.Info("Received a ping message")
		ctx.Respond(pong)

	default:

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

	// Prepare a remote env that listens to 8080
	remoteConfig := remote.Configure("127.0.0.1", 8080)

	// Configure a cluster on top of the above remote env
	clusterKind := cluster.NewKind(
		"Ponger",
		actor.PropsFromProducer(func() actor.Actor {
			return &ponger{}
		}))
	// Configure cluster on top of the above remote env
	cp, err := etcd.NewWithConfig("/protoactor", clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		slog.Error("Failed to create an etcd provider", "error", err)
		os.Exit(1)
	}
	lookup := disthash.New()
	clusterConfig := cluster.Configure("cluster-example", cp, lookup, remoteConfig, cluster.WithKinds(clusterKind))
	c := cluster.New(system, clusterConfig)

	// Manage the cluster node's lifecycle
	c.StartMember()
	defer c.Shutdown(false)

	// Run till a signal comes
	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, os.Kill)
	<-finish
}
