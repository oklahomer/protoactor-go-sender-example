package main

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"log"
	"os"
	"os/signal"
	"protoactor-go-sender-example/remote/messages"
)

func main() {
	// Setup actor system
	system := actor.NewActorSystem()

	// Setup remote env that listens to 8080
	config := remote.Configure("127.0.0.1", 8080)
	remoting := remote.NewRemote(system, config)
	remoting.Start()

	// Run pong actor that receives ping payload, and then send back pong payload
	pongProps := actor.PropsFromFunc(func(ctx actor.Context) {
		switch msg := ctx.Message().(type) {
		case *messages.Ping:
			pong := &messages.Pong{Cnt: msg.Cnt}
			log.Print("Received ping message")
			ctx.Respond(pong)

		default:

		}
	})
	pongPid, err := system.Root.SpawnNamed(pongProps, "pongActorID")
	if err != nil {
		log.Fatalf("Failed to spawn actor: %s.", err.Error())
	}
	log.Printf("Actor is running. Address: %s. ID: %s.", pongPid.Address, pongPid.Id)

	// Run till signal comes
	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, os.Kill)
	<-finish

	remoting.Shutdown(false)
}
