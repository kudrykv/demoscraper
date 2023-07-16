package main

import (
	"context"
	"log"
	"os"
	"os/signal"
)

func main() {
	setupFlags()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	waitChan := make(chan struct{})

	signal.Notify(signalChan, os.Interrupt)

	go func() {
		<-signalChan

		log.Println("Received interrupt signal")
		cancel()
	}()

	go func() {
		run(ctx)

		close(waitChan)
	}()

	<-waitChan
}
