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

// ponger handles the incoming messages.
// This supports gRPC Ponger service and plain message handling.
type ponger struct {
}

var _ messages.Ponger = (*ponger)(nil)

// Init takes care of the initialization.
func (p *ponger) Init(ctx cluster.GrainContext) {
	log.Printf("Initializing ponger: %s", ctx.Self().GetId())
}

// Terminate takes care of the finalization.
func (p *ponger) Terminate(ctx cluster.GrainContext) {
	// Do finalization if required. e.g. Store the current state to storage and switch its behavior to reject further messages.
	// This method is called when a pre-configured idle interval passes from the last message reception.
	// The actor will be re-initialized when a message comes for the next time.
	// Terminating the idle actor is effective to free unused server resource.
	//
	// A poison pill message is enqueued right after this method execution and the actor eventually stops.
	log.Printf("Terminating ponger: %s", ctx.Self().GetId())
}

// ReceiveDefault is a default method to receive and handle incoming messages.
func (p *ponger) ReceiveDefault(ctx cluster.GrainContext) {
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
	sender := ctx.Sender()
	log.Printf("Received Ping call from sender. Address: %s. ID: %s.", sender.GetAddress(), sender.GetId())

	pong := &messages.PongMessage{
		Cnt: ping.Cnt,
	}
	return pong, nil
}

func main() {
	// Set up actor system
	system := actor.NewActorSystem()

	// Register a ponger constructor.
	// This is called when the wrapping PongerActor is initialized.
	// PongerActor proxies messages to ponger's corresponding methods.
	messages.PongerFactory(func() messages.Ponger {
		return &ponger{}
	})

	// Prepare a remote env that listens to 8080.
	// Messages are sent to this port.
	remoteConfig := remote.Configure("127.0.0.1", 8080)

	// Configure a cluster provider so the Ponger grain can work as a cluster member.
	cp, err := consul.New()
	if err != nil {
		log.Fatal(err)
	}

	// Register an actor constructor for the Ponger kind.
	// With this registration, the message sender and other cluster nodes know this node is capable of providing a Ponger.
	// PongerActor will implicitly be initialized when the first message comes in.
	clusterKind := cluster.NewKind(
		"Ponger",
		actor.PropsFromProducer(func() actor.Actor {
			return &messages.PongerActor{
				// The actor stops when 10 seconds is passed after the last message reception.
				// Ponger.Terminate() is called on the finalization.
				// When the next message comes, the actor is revitalized so Ponger.Init() is called again.
				Timeout: 10 * time.Second,
			}
		}))
	lookup := disthash.New()
	clusterConfig := cluster.Configure("cluster-example", cp, lookup, remoteConfig, cluster.WithKinds(clusterKind))
	c := cluster.New(system, clusterConfig)

	// Manage the cluster node's lifecycle
	// Use StartClient() when this process is not a member of the cluster nodes but required to send messages to cluster grains.
	c.StartMember()
	defer c.Shutdown(false)

	// Run till a signal comes
	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, os.Kill)
	<-finish
}
