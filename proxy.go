package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func runProxy(ctx context.Context, wg *sync.WaitGroup) {
	defer func() {
		log.Println("Proxy shutting down")
		wg.Done()
	}()

	tp, err := makeTraceProvider("proxy")
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.HideBanner = true
	go func() {
		<-ctx.Done()
		e.Shutdown(context.Background())
	}()

	e.Use(TracingMiddleware(e, tp))
	e.Use(LogTraceIDMiddleware())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "proxy")
	})

	targets := []*middleware.ProxyTarget{
		{
			URL: mustParseUrl(fmt.Sprintf("http://localhost:%d", backendPort)),
		},
	}
	g := e.Group("/proxy")
	g.Use(middleware.ProxyWithConfig(middleware.ProxyConfig{
		Balancer: middleware.NewRoundRobinBalancer(targets),
		Rewrite: map[string]string{
			"/proxy/*": "/$1",
		},
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithTracerProvider(tp),
		),
	}))

	log.Println("Proxy starting up")
	addr := fmt.Sprintf(":%d", proxyPort)
	e.Start(addr)
}

func mustParseUrl(s string) *url.URL {
	ret, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return ret
}
