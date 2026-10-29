package main

import (
	"fmt"
	"math/rand"
)

func ra(ch chan<- int) {
	for _ = range 10 {
		ch <- rand.Intn(10)
	}
	close(ch)
}

func main() {
	ch := make(chan int)
	go ra(ch)
	for val := range ch {
		fmt.Println(val * val)
	}
}
