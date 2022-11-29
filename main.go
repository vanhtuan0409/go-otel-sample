package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
)

const (
	proxyPort      = 8081
	backendPort    = 8080
	jaegerEndpoint = "127.0.0.1"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(3)

	go runBackend(ctx, &wg)
	go runProxy(ctx, &wg)
	go runConsumer(ctx, &wg)

	wg.Wait()
	log.Println("All services shutdown")
}
