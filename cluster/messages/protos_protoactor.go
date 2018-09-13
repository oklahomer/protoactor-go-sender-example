
package messages


import errors "errors"
import log "log"
import actor "github.com/AsynkronIT/protoactor-go/actor"
import remote "github.com/AsynkronIT/protoactor-go/remote"
import cluster "github.com/AsynkronIT/protoactor-go/cluster"

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

	
var xPongerFactory func() Ponger

func PongerFactory(factory func() Ponger) {
	xPongerFactory = factory
}

func GetPongerGrain(id string) *PongerGrain {
	return &PongerGrain{ID: id}
}

type Ponger interface {
	Init(id string)
		
	SendPing(*Ping, cluster.GrainContext) (*Pong, error)
		
}
type PongerGrain struct {
	ID string
}

	
func (g *PongerGrain) SendPing(r *Ping) (*Pong, error) {
	return g.SendPingWithOpts(r, cluster.DefaultGrainCallOptions())
}

func (g *PongerGrain) SendPingWithOpts(r *Ping, opts *cluster.GrainCallOptions) (*Pong, error) {
	fun := func() (*Pong, error) {
			pid, statusCode := cluster.Get(g.ID, "Ponger")
			if statusCode != remote.ResponseStatusCodeOK {
				return nil, fmt.Errorf("Get PID failed with StatusCode: %v", statusCode)
			}
			bytes, err := proto.Marshal(r)
			if err != nil {
				return nil, err
			}
			request := &cluster.GrainRequest{Method: "SendPing", MessageData: bytes}
			response, err := pid.RequestFuture(request, opts.Timeout).Result()
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
				return nil, errors.New("Unknown response")
			}
		}
	
	var res *Pong
	var err error
	for i := 0; i < opts.RetryCount; i++ {
		res, err = fun()
		if err == nil {
			return res, nil
		} else {
			if opts.RetryAction != nil {
				opts.RetryAction(i)
			}
		}
	}
	return nil, err
}

func (g *PongerGrain) SendPingChan(r *Ping) (<-chan *Pong, <-chan error) {
	return g.SendPingChanWithOpts(r, cluster.DefaultGrainCallOptions())
}

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
	

type PongerActor struct {
	inner Ponger
}

func (a *PongerActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		a.inner = xPongerFactory()
		id := ctx.Self().Id
		a.inner.Init(id[7:]) //skip "remote$"

	case actor.AutoReceiveMessage: //pass
	case actor.SystemMessage: //pass

	case *cluster.GrainRequest:
		switch msg.Method {
			
		case "SendPing":
			req := &Ping{}
			err := proto.Unmarshal(msg.MessageData, req)
			if err != nil {
				log.Fatalf("[GRAIN] proto.Unmarshal failed %v", err)
			}
			r0, err := a.inner.SendPing(req, ctx)
			if err == nil {
				bytes, err := proto.Marshal(r0)
				if err != nil {
					log.Fatalf("[GRAIN] proto.Marshal failed %v", err)
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

	


//Why has this been removed?
//This should only be done on servers of the below Kinds
//Clients should not be forced to also be servers

//func init() {
//	
//	remote.Register("Ponger", actor.FromProducer(func() actor.Actor {
//		return &PongerActor {}
//		})		)
//	
//}



// type ponger struct {
//	cluster.Grain
// }

// func (*ponger) SendPing(r *Ping, cluster.GrainContext) (*Pong, error) {
// 	return &Pong{}, nil
// }



// func init() {
// 	//apply DI and setup logic

// 	PongerFactory(func() Ponger { return &ponger{} })

// }





