package main

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/cluster/automanaged"
	"github.com/AsynkronIT/protoactor-go/remote"
	"log"
	"os"
	"os/signal"
	"protoactor-go-sender-example/cluster/messages"
	"time"
)

var cnt uint64 = 0

type pingActor struct {
	cluster *cluster.Cluster
	cnt     uint
}

func (p *pingActor) Receive(ctx actor.Context) {
	switch ctx.Message().(type) {
	case struct{}:
		cnt += 1
		ping := &messages.PingMessage{
			Cnt: cnt,
		}

		grainPid, statusCode := p.cluster.Get("ponger-1", "Ponger")
		if statusCode != remote.ResponseStatusCodeOK {
			log.Printf("Get PID failed with StatusCode: %v", statusCode)
			return
		}

		future := ctx.RequestFuture(grainPid, ping, time.Second)
		result, err := future.Result()
		if err != nil {
			log.Print(err.Error())
			return
		}
		log.Printf("Received %v", result)

	case *messages.PongMessage:
		// Never comes here.
		// When the pong actor responds to the sender,
		// the sender is not a ping actor but a future process.
		log.Print("Received pong message")

	}
}

func main() {
	// Setup actor system
	system := actor.NewActorSystem()

	// Prepare remote env that listens to 8081
	remoteConfig := remote.Configure("127.0.0.1", 8081)

	// Configure cluster on top of the above remote env
	clusterProvider := automanaged.NewWithConfig(1*time.Second, 6330, "localhost:6331")
	clusterConfig := cluster.Configure("cluster-example", clusterProvider, remoteConfig)
	c := cluster.New(system, clusterConfig)
	c.StartClient()

	// Start ping actor that periodically send "ping" payload to "Ponger" cluster grain
	pingProps := actor.PropsFromProducer(func() actor.Actor {
		return &pingActor{
			cluster: c,
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
