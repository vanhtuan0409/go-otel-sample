package main

import (
	"context"
	"log"
	"sync"
)

func runConsumer(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		log.Println("Consumer shutting down")
		wg.Done()
	}()

	log.Println("Consumer started")
}
