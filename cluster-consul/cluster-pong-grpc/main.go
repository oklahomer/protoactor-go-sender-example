package main

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/cluster/consul"
	"github.com/AsynkronIT/protoactor-go/remote"
	"log"
	"os"
	"os/signal"
	"protoactor-go-sender-example/cluster/messages"
	"time"
)

// ponger handles the incoming messages.
// This supports gRPC Ponger service and plain message handling.
type ponger struct {
	cluster.Grain
}

var _ messages.Ponger = (*ponger)(nil)

// Terminate takes care of the finalization.
func (p *ponger) Terminate() {
	// Do finalization if required. e.g. Store the current state to storage and switch its behavior to reject further messages.
	// This method is called when a pre-configured idle interval passes from the last message reception.
	// The actor will be re-initialized when a message comes for the next time.
	// Terminating the idle actor is effective to free unused server resource.
	//
	// A poison pill message is enqueued right after this method execution and the actor eventually stops.
	log.Printf("Terminating ponger: %s", p.ID())
}

// ReceiveDefault is a default method to receive and handle incoming messages.
func (p *ponger) ReceiveDefault(ctx actor.Context) {
	log.Printf("A plain message is sent from sender: %+v", ctx.Sender())

	switch msg := ctx.Message().(type) {
	case *messages.PingMessage:
		log.Print("Received ping message")
		pong := &messages.PongMessage{Cnt: msg.Cnt}
		ctx.Respond(pong)

	default:

	}
}

// Ping is called when gRPC-based request is sent against Ponger service.
func (p *ponger) Ping(ping *messages.PingMessage, ctx cluster.GrainContext) (*messages.PongMessage, error) {
	// The sender process is not a sending actor, but a future process
	log.Printf("Received Ping call from sender: %+v", ctx.Sender())

	pong := &messages.PongMessage{
		Cnt: ping.Cnt,
	}
	return pong, nil
}

func main() {
	// Setup actor system
	system := actor.NewActorSystem()

	// Register ponger constructor.
	// This is called when the wrapping PongerActor is initialized.
	// PongerActor proxies messages to ponger's corresponding methods.
	messages.PongerFactory(func() messages.Ponger {
		return &ponger{}
	})

	// Prepare remote env that listens to 8080
	// Messages are sent to this port.
	remoteConfig := remote.Configure("127.0.0.1", 8080)

	// Configure cluster provider to work as a cluster member.
	cp, err := consul.New()
	if err != nil {
		log.Fatal(err)
	}

	// Register an actor constructor for the Ponger kind.
	// With this registration, the message sender and other cluster nodes know this node is capable of providing Ponger.
	// PongerActor will implicitly be initialized when the first message comes.
	clusterKind := cluster.NewKind(
		"Ponger",
		actor.PropsFromProducer(func() actor.Actor {
			return &messages.PongerActor{
				// The actor stops when 10 seconds passed since the last message reception.
				// When the next
				Timeout: 10 * time.Second,
			}
		}))
	clusterConfig := cluster.Configure("cluster-example", cp, remoteConfig, clusterKind)
	c := cluster.New(system, clusterConfig)

	// Start as a cluster member.
	// Use StartClient() when this process is not a member of cluster nodes but required to send messages to cluster grains.
	c.Start()

	// Run till signal comes
	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, os.Kill)
	<-finish
}
