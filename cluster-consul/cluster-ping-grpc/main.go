package main

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/clusterproviders/consul"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"
	"log"
	"os"
	"os/signal"
	"protoactor-go-sender-example/cluster/messages"
	"time"
)

var cnt uint64 = 0

type pingActor struct {
	system *actor.ActorSystem
	cnt    uint
}

func (p *pingActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case struct{}:
		cnt += 1
		ping := &messages.PingMessage{
			Cnt: cnt,
		}

		client := messages.GetPongerGrainClient(cluster.GetCluster(p.system), "ponger-1")
		pong, err := client.Ping(ping, cluster.WithTimeout(time.Second), cluster.WithRetry(5))
		if err != nil {
			log.Print(err.Error())
			return
		}
		log.Printf("Received %v", pong)

	case *messages.PongMessage:
		// Never comes here.
		// When the pong grain responds to the sender's gRPC call,
		// the sender is not a ping actor but a future process.
		log.Print("Received pong message")

	}
}

func main() {
	// Set up actor system
	system := actor.NewActorSystem()

	// Prepare remote env that listens to 8081
	remoteConfig := remote.Configure("127.0.0.1", 8081)

	// Configure cluster on top of the above remote env
	cp, err := consul.New()
	if err != nil {
		log.Fatal(err)
	}
	lookup := disthash.New()
	clusterConfig := cluster.Configure("cluster-example", cp, lookup, remoteConfig)
	c := cluster.New(system, clusterConfig)

	// Manage the cluster client's lifecycle
	c.StartClient() // Configure as a client
	defer c.Shutdown(false)

	// Start a ping actor that periodically send a "ping" payload to the "Ponger" cluster grain
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
			log.Print("Finish")
			return

		}
	}
}
