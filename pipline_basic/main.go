package main

import (
	"fmt"
	"math/rand"
	"sync"
)

func generator(wg *sync.WaitGroup) <-chan int {
	ch := make(chan int)
	randomn := func() int { return rand.Intn(9000000000000) }
	go func() {
		defer wg.Done()
		defer close(ch)
		for ;; {
			ch <- randomn()
		}
	}()
	return ch
}

func multiple(wg *sync.WaitGroup, steam <-chan int, n int) <-chan int {
	ch := make(chan int)
	go func() {
		defer wg.Done()
		defer close(ch)
		for i := range steam {
			ch <- n * i
		}
	}()
	return ch
}

func printout(wg *sync.WaitGroup, steam <-chan int) {
	go func() {
		defer wg.Done()
		for i := range steam {
			fmt.Println(i)
		}
	}()
}

func main() {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	step1 := generator(wg)

	wg.Add(1)
	step2 := multiple(wg, step1, 2)

	printout(wg, step2)

	wg.Wait()
}
