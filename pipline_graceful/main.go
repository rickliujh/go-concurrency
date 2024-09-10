package main

import (
	"fmt"
	"math/rand"
	"time"
)

func generator(done chan bool) <-chan int {
	ch := make(chan int)
	randomn := func() int { return rand.Intn(9000000000000) }
	go func() {
		defer close(ch)
		for {
			select {
			case <-done:
				fmt.Println("generator gracefully shutdown")
				return
			default:
				ch <- randomn()
			}
		}
	}()
	return ch
}

func multiple(steam <-chan int, n int, done chan bool) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for {
			select {
			case <-done:
				fmt.Println("multiple gracefully shutdown")
				return
			case i := <-steam:
				ch <- n * i
			}
		}
	}()
	return ch
}

func printout(steam <-chan int, done chan bool) {
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("printout gracefully shutdown")
				return
			case i := <-steam:
				fmt.Println(i)
			}
		}
	}()
}

func main() {
	done := make(chan bool)

	step1 := generator(done)

	step2 := multiple(step1, 2, done)

	printout(step2, done)

	time.Sleep(time.Second * 5)

	close(done)
}
