package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
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

func run(ctx context.Context) {
	select {
	case <-time.After(5 * time.Second):
		log.Println("Done")

		return
	case <-ctx.Done():
		log.Println("Cancelled")

		return
	}
}
