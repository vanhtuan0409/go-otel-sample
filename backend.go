package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
)

func runBackend(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		log.Println("Backend shutting down")
		wg.Done()
	}()

	e := echo.New()
	e.HideBanner = true
	go func() {
		<-ctx.Done()
		e.Shutdown(context.Background())
	}()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "backend")
	})
	e.GET("/echo", echoHandler)

	log.Println("Backend starting up")
	addr := fmt.Sprintf(":%d", backendPort)
	e.Start(addr)
}

func echoHandler(c echo.Context) error {
	req := c.Request()
	resp := struct {
		Request string
		Header  http.Header
	}{
		Request: req.URL.String(),
		Header:  req.Header,
	}
	return c.JSON(http.StatusOK, &resp)
}
