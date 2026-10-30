package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type RanRes struct {
	id  int
	val int
}

func merge(ctx context.Context, channels ...chan RanRes) <-chan RanRes {
	out := make(chan RanRes)
	var wg sync.WaitGroup
	c := func(ctx context.Context, ch <-chan RanRes) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case val, ok := <-ch:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case out <- val:
				}
			}
		}
	}

	for _, val := range channels {
		wg.Add(1)
		go c(ctx, val)

	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func ran(ctx context.Context, wg *sync.WaitGroup, id int, ch chan<- RanRes) {
	defer wg.Done()
	defer close(ch)
	for _ = range 10 {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second / 2):
		}
		select {
		case <-ctx.Done():
			return
		case ch <- RanRes{
			id:  id,
			val: rand.Intn(10),
		}:
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	ch1 := make(chan RanRes)
	ch2 := make(chan RanRes)
	ch3 := make(chan RanRes)

	chanals := []chan RanRes{ch1, ch2, ch3}
	var wg sync.WaitGroup

	for i, val := range chanals {
		wg.Add(1)
		go ran(ctx, &wg, i+1, val)
	}

	out := merge(ctx, chanals...)

	for _ = range 5 {
		val := <-out
		fmt.Printf("Worker %d: %d\n", val.id, val.val)
	}
	cancel()

	wg.Wait()
}
