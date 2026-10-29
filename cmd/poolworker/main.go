package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	wg     *sync.WaitGroup
	result chan<- Result
	jobs   <-chan int
	id     int
}

type Result struct {
	idWorker int
	after    int
	befor    int
}

func (w *Worker) worker(ctx context.Context) {
	defer w.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case val, ok := <-w.jobs:
			if !ok {
				return
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(500 * time.Millisecond):
			}
			select {
			case <-ctx.Done():
				return
			case w.result <- Result{
				idWorker: w.id,
				after:    val,
				befor:    val * val,
			}:
			}

		}
	}
}

func main() {
	var wg sync.WaitGroup

	jobs := make(chan int)

	result := make(chan Result)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	go func(ctx context.Context) {
		defer close(jobs)
		for i := range 1000 {
			select {
			case <-ctx.Done():
				return
			case jobs <- i + 1:
			}
		}
	}(ctx)

	for i := range 5 {
		work := Worker{
			wg:     &wg,
			result: result,
			jobs:   jobs,
			id:     i + 1,
		}
		wg.Add(1)
		go work.worker(ctx)
	}

	go func() {
		wg.Wait()
		close(result)
	}()
	count := 0

	mp := make(map[int]int)

	for res := range result {
		count++
		mp[res.idWorker]++
		fmt.Printf("Воркер %d: %d ^ 2 = %d\n", res.idWorker, res.after, res.befor)
	}
	for k, v := range mp {
		fmt.Printf("Worker %d done %d task\n", k, v)
	}
	fmt.Printf("Отработало %d task\n", count)
}
