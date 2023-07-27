package main

import (
	"context"
	"fmt"
	"github.com/ybalcin/wallet-service/cmd"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-shutdown
		fmt.Printf("%s signal received, server is shutting down...\n", s.String())
		cancel()
	}()

	if err := cmd.RunApi(ctx); err != nil {
		panic(err)
	}
}
