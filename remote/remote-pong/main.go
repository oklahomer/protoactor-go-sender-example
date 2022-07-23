package main

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/remote"
	"log"
	"os"
	"os/signal"
	"protoactor-go-sender-example/remote/messages"
)

func main() {
	// Set up the actor system
	system := actor.NewActorSystem()

	// Set up a remote env that listens to 8080
	config := remote.Configure("127.0.0.1", 8080)
	remoter := remote.NewRemote(system, config)
	remoter.Start()
	defer remoter.Shutdown(false)

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

	// Run till a signal comes
	finish := make(chan os.Signal, 1)
	signal.Notify(finish, os.Interrupt, os.Kill)
	<-finish
}
