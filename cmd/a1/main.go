package main

import "fmt"

func num(ch chan<- int, turnNum <-chan struct{}, turnLet chan<- struct{}) {
	for i := range 5 {
		<-turnNum
		ch <- i + 1
		turnLet <- struct{}{}
	}
	close(ch)
}

func run(ch chan<- rune, turnLet <-chan struct{}, turnNum chan<- struct{}) {
	for i := range 5 {
		<-turnLet
		ch <- rune('a' + i)
		if i < 4 {
			turnNum <- struct{}{}
		}
	}
	close(ch)
}

func main() {
	numCh := make(chan int)
	letterCh := make(chan rune)
	numTurn := make(chan struct{})
	letterTurn := make(chan struct{})
	go num(numCh, numTurn, letterTurn)
	go run(letterCh, letterTurn, numTurn)
	numTurn <- struct{}{}
	numOpen, letterOpen := true, true
	for numOpen || letterOpen {
		select {
		case val, ok := <-numCh:
			if !ok {
				numCh = nil
				numOpen = false
			} else {
				fmt.Printf("Число: %d\n", val)
			}
		case val, ok := <-letterCh:
			if !ok {
				letterCh = nil
				letterOpen = false
			} else {
				fmt.Printf("Буква: %c\n", val)
			}
		}
	}
}
