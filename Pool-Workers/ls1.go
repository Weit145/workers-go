package poolworkers

// package main

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// type Satats struct {
// 	Mx           sync.Mutex
// 	Progress     int
// 	EvenCount    int
// 	OddCount     int
// 	WorkerCounts map[int]int
// }

// func NewSatats() *Satats {
// 	return &Satats{
// 		WorkerCounts: make(map[int]int),
// 	}
// }

// func (s *Satats) Add(workerID int, number int) {
// 	s.Mx.Lock()
// 	defer s.Mx.Unlock()

// 	s.Progress++
// 	s.WorkerCounts[workerID]++

// 	if number%2 == 0 {
// 		s.EvenCount++
// 	} else {
// 		s.OddCount++
// 	}
// }

// func workers(id int, jobs <-chan int, stats *Satats, wg *sync.WaitGroup) {
// 	defer wg.Done()

// 	for number := range jobs {
// 		time.Sleep(1 * time.Second)

// 		stats.Add(id, number)

// 		fmt.Printf(
// 			"Worker %d: %d\n",
// 			id,
// 			number,
// 		)
// 	}
// }

// func main() {
// 	jobs := make(chan int)

// 	var wg sync.WaitGroup
// 	stats := NewSatats()

// 	for id := 1; id <= 4; id++ {
// 		wg.Add(1)
// 		go workers(id, jobs, stats, &wg)
// 	}

// 	for number := 1; number < 101; number++ {
// 		jobs <- number
// 	}

// 	close(jobs)

// 	wg.Wait()
// 	fmt.Println()
// 	fmt.Println("Обработано:", stats.Progress)
// 	fmt.Println("Чётных:", stats.EvenCount)
// 	fmt.Println("Нечётных:", stats.OddCount)

// 	for workerID, count := range stats.WorkerCounts {
// 		fmt.Printf("Worker %d: %d задач\n", workerID, count)
// 	}
// }
