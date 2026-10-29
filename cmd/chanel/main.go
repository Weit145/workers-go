package main

import (
	"context"
	"log/slog"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// ctx, stop := context.WithTimeout(context.Background(), 10*time.Second)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	ch := make(chan int, 100)

	for i := 1; i < 6; i++ {
		wg.Add(1)
		go Worker(ctx, i, &wg, ch)
	}

	for i := range 100 {
		ch <- i
		// select {
		// case <-ctx.Done():
		// 	continue
		// default:
		// 	ch <- i
		// }
	}
	wg.Wait()
}

func Worker(ctx context.Context, id int, wg *sync.WaitGroup, ch <-chan int) {
	for {
		select {
		case <-ctx.Done():
			slog.InfoContext(ctx, "Cancel worker: ", slog.Int("id", id))
			wg.Done()
			return
		case n := <-ch:
			slog.InfoContext(ctx, "Worker: ", slog.Int("id", id), slog.Int("n", n))
			time.Sleep(2 * time.Second)
		}
	}
}
