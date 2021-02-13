package main

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/oklahomer/protoactor-go-sender-example/remote/messages"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	system := actor.NewActorSystem()

	config := remote.Configure("127.0.0.1", 8080)
	remoting := remote.NewRemote(system, config)
	remoting.Start()

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

	finish := make(chan os.Signal, 1)
	signal.Notify(finish, syscall.SIGINT)
	signal.Notify(finish, syscall.SIGTERM)
	<-finish

	remoting.Shutdown(false)
}
