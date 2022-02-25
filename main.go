package main

import (
	"fmt"
	"github.com/DictumMortuum/edgemax-exporter/internal/metrics"
	"github.com/xonvanetta/shutdown/pkg/shutdown"
)

func main() {
	metrics.Init()

	serverDead := make(chan struct{})
	s := metrics.NewServer("9191", metrics.NewClient())

	go func() {
		s.ListenAndServe()
		close(serverDead)
	}()

	ctx := shutdown.Context()

	go func() {
		<-ctx.Done()
		s.Stop()
	}()

	select {
	case <-ctx.Done():
	case <-serverDead:
	}

	version := "0.0.1"
	fmt.Printf("edgemax-exporter v%s HTTP server stopped\n", version)
}
