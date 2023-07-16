package main

import (
	"context"
	"log"
	"time"
)

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
