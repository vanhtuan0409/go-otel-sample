package main

import (
	"context"
	"log"
	"sync"
)

func runProxy(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		log.Println("Proxy shutting down")
		wg.Done()
	}()

	log.Println("Proxy started")
}
