package main

import (
	"context"
	"log"
	"sync"
)

func runBackend(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		log.Println("Backend shutting down")
		wg.Done()
	}()

	log.Println("Backend started")
}
