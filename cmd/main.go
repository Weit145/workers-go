package main

import (
	"fmt"
	"sync"
	"time"
)

func workers(id int, ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range ch {
		fmt.Printf("Worker %d: %d*%d = %d\n", id, num, num, num*num)
		time.Sleep(2 * time.Second)
	}
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan int)
	for id := 1; id <= 5; id++ {
		wg.Add(1)
		go workers(id, ch, &wg)
	}
	for num := 1; num <= 30; num++ {
		ch <- num
	}
	close(ch)
	wg.Wait()
}
