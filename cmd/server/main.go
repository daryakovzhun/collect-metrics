package main

import (
	"context"
	"github.com/daryakovzhun/collect-metrics/internal/app"
	"golang.org/x/sync/errgroup"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			cancel()
			return nil
		case <-egCtx.Done():
			return nil
		}
	})

	eg.Go(func() error {
		return app.Run(egCtx)
	})

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
		return
	}
}
