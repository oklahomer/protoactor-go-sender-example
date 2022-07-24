// Package messages is generated by protoactor-go/protoc-gen-gograin@0.1.0
package messages

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	logmod "github.com/asynkron/protoactor-go/log"
	"google.golang.org/protobuf/proto"
)

var (
	plog = logmod.New(logmod.InfoLevel, "[GRAIN][messages]")
	_    = proto.Marshal
	_    = fmt.Errorf
	_    = math.Inf
)

// SetLogLevel sets the log level.
func SetLogLevel(level logmod.Level) {
	plog.SetLevel(level)
}

var xPongerFactory func() Ponger

// PongerFactory produces a Ponger
func PongerFactory(factory func() Ponger) {
	xPongerFactory = factory
}

// GetPongerGrainClient instantiates a new PongerGrainClient with given Identity
func GetPongerGrainClient(c *cluster.Cluster, id string) *PongerGrainClient {
	if c == nil {
		panic(fmt.Errorf("nil cluster instance"))
	}
	if id == "" {
		panic(fmt.Errorf("empty id"))
	}
	return &PongerGrainClient{Identity: id, cluster: c}
}

// GetPongerKind instantiates a new cluster.Kind for Ponger
func GetPongerKind(opts ...actor.PropsOption) *cluster.Kind {
	props := actor.PropsFromProducer(func() actor.Actor {
		return &PongerActor{
			Timeout: 60 * time.Second,
		}
	}, opts...)
	kind := cluster.NewKind("Ponger", props)
	return kind
}

// GetPongerKind instantiates a new cluster.Kind for Ponger
func NewPongerKind(factory func() Ponger, timeout time.Duration ,opts ...actor.PropsOption) *cluster.Kind {
	xPongerFactory = factory
	props := actor.PropsFromProducer(func() actor.Actor {
		return &PongerActor{
			Timeout: timeout,
		}
	}, opts...)
	kind := cluster.NewKind("Ponger", props)
	return kind
}

// Ponger interfaces the services available to the Ponger
type Ponger interface {
	Init(ctx cluster.GrainContext)
	Terminate(ctx cluster.GrainContext)
	ReceiveDefault(ctx cluster.GrainContext)
	Ping(*PingMessage, cluster.GrainContext) (*PongMessage, error)
	
}

// PongerGrainClient holds the base data for the PongerGrain
type PongerGrainClient struct {
	Identity      string
	cluster *cluster.Cluster
}

// Ping requests the execution on to the cluster with CallOptions
func (g *PongerGrainClient) Ping(r *PingMessage, opts ...cluster.GrainCallOption) (*PongMessage, error) {
	bytes, err := proto.Marshal(r)
	if err != nil {
		return nil, err
	}
	reqMsg := &cluster.GrainRequest{MethodIndex: 0, MessageData: bytes}
	resp, err := g.cluster.Call(g.Identity, "Ponger", reqMsg, opts...)
	if err != nil {
		return nil, err
	}
	switch msg := resp.(type) {
	case *cluster.GrainResponse:
		result := &PongMessage{}
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


// PongerActor represents the actor structure
type PongerActor struct {
	ctx     cluster.GrainContext
	inner   Ponger
	Timeout time.Duration
}

// Receive ensures the lifecycle of the actor for the received message
func (a *PongerActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started: //pass
	case *cluster.ClusterInit:
		a.ctx = cluster.NewGrainContext(ctx, msg.Identity, msg.Cluster)
		a.inner = xPongerFactory()
		a.inner.Init(a.ctx)

		if a.Timeout > 0 {
			ctx.SetReceiveTimeout(a.Timeout)
		}
	case *actor.ReceiveTimeout:		
		ctx.Poison(ctx.Self())
	case *actor.Stopped:
		a.inner.Terminate(a.ctx)
	case actor.AutoReceiveMessage: // pass
	case actor.SystemMessage: // pass

	case *cluster.GrainRequest:
		switch msg.MethodIndex {
		case 0:
			req := &PingMessage{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				plog.Error("Ping(PingMessage) proto.Unmarshal failed.", logmod.Error(err))
				resp := &cluster.GrainErrorResponse{Err: err.Error()}
				ctx.Respond(resp)
				return
			}
			r0, err := a.inner.Ping(req, a.ctx)
			if err != nil {
				resp := &cluster.GrainErrorResponse{Err: err.Error()}
				ctx.Respond(resp)
				return
			}
			bytes, err := proto.Marshal(r0)
			if err != nil {
				plog.Error("Ping(PingMessage) proto.Marshal failed", logmod.Error(err))
				resp := &cluster.GrainErrorResponse{Err: err.Error()}
				ctx.Respond(resp)
				return
			}
			resp := &cluster.GrainResponse{MessageData: bytes}
			ctx.Respond(resp)
		
		}
	default:
		a.inner.ReceiveDefault(a.ctx)
	}
}
