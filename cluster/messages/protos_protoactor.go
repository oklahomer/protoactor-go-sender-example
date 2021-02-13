package messages

import (
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/gogo/protobuf/proto"
)

var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

var system *actor.ActorSystem

func SetSystem(sys *actor.ActorSystem) {
	system = sys
}

var c *cluster.Cluster

// SetCluster pass created cluster
func SetCluster(cluster *cluster.Cluster) {
	c = cluster
}

var xPongerFactory func() Ponger

// PongerFactory produces a Ponger
func PongerFactory(factory func() Ponger) {
	xPongerFactory = factory
}

// GetPongerGrain instantiates a new PongerGrain with given ID
func GetPongerGrain(id string) *PongerGrain {
	return &PongerGrain{ID: id}
}

// Ponger interfaces the services available to the Ponger
type Ponger interface {
	Init(id string)
	Terminate()

	SendPing(*Ping, cluster.GrainContext) (*Pong, error)
}

// PongerGrain holds the base data for the PongerGrain
type PongerGrain struct {
	ID string
}

// SendPing requests the execution on to the cluster using default options
func (g *PongerGrain) SendPing(r *Ping) (*Pong, error) {
	return g.SendPingWithOpts(r, cluster.DefaultGrainCallOptions(c))
}

// SendPingWithOpts requests the execution on to the cluster
func (g *PongerGrain) SendPingWithOpts(r *Ping, opts ...*cluster.GrainCallOptions) (*Pong, error) {
	bytes, err := proto.Marshal(r)
	if err != nil {
		return nil, err
	}

	request := &cluster.GrainRequest{MethodIndex: 0, MessageData: bytes}
	response, err := c.Call(g.ID, "Ponger", request, opts...)
	if err != nil {
		return nil, err
	}
	switch msg := response.(type) {
	case *cluster.GrainResponse:
		result := &Pong{}
		err = proto.Unmarshal(msg.MessageData, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	case *cluster.GrainErrorResponse:
		return nil, errors.New(msg.Err)
	default:
		return nil, errors.New("unknown response")
	}
}

// SendPingChan allows to use a channel to execute the method using default options
func (g *PongerGrain) SendPingChan(r *Ping) (<-chan *Pong, <-chan error) {
	return g.SendPingChanWithOpts(r, cluster.DefaultGrainCallOptions(c))
}

// SendPingChanWithOpts allows to use a channel to execute the method
func (g *PongerGrain) SendPingChanWithOpts(r *Ping, opts *cluster.GrainCallOptions) (<-chan *Pong, <-chan error) {
	c := make(chan *Pong)
	e := make(chan error)
	go func() {
		res, err := g.SendPingWithOpts(r, opts)
		if err != nil {
			e <- err
		} else {
			c <- res
		}
		close(c)
		close(e)
	}()
	return c, e
}

// PongerActor represents the actor structure
type PongerActor struct {
	inner   Ponger
	Timeout *time.Duration
}

// Receive ensures the lifecycle of the actor for the received message
func (a *PongerActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		a.inner = xPongerFactory()
		id := ctx.Self().Id
		a.inner.Init(id[7:]) // skip "remote$"
		if a.Timeout != nil {
			ctx.SetReceiveTimeout(*a.Timeout)
		}
	case *actor.ReceiveTimeout:
		a.inner.Terminate()
		system.Root.PoisonFuture(ctx.Self()).Wait()

	case actor.AutoReceiveMessage: // pass
	case actor.SystemMessage: // pass

	case *cluster.GrainRequest:
		switch msg.MethodIndex {

		case 0:
			req := &Ping{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				log.Fatalf("[GRAIN] proto.Unmarshal failed %v", err)
			}
			r0, err := a.inner.SendPing(req, ctx)
			if err == nil {
				bytes, errMarshal := proto.Marshal(r0)
				if errMarshal != nil {
					log.Fatalf("[GRAIN] proto.Marshal failed %v", errMarshal)
				}
				resp := &cluster.GrainResponse{MessageData: bytes}
				ctx.Respond(resp)
			} else {
				resp := &cluster.GrainErrorResponse{Err: err.Error()}
				ctx.Respond(resp)
			}

		}
	default:
		log.Printf("Unknown message %v", msg)
	}
}
