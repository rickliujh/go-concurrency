package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func is_prime(n int) bool {
	for i := n - 1; i > 1; i-- {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func or_done[T any, K any](done <-chan K, stream <-chan T) <-chan T {
	relay_stream := make(chan T)
	go func() {
		defer close(relay_stream)
		for {
			select {
			case <-done:
				close(relay_stream)
			case data := <-stream:
				select {
				case <-done:
					return
				default:
					relay_stream <- data
				}
			}
		}
	}()
	return relay_stream
}

func generator[T any, K any](done chan K, fn func() T) <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for {
			select {
			case <-done:
				return
			default:
				for {
					ch <- fn()
				}
			}
		}
	}()
	return ch
}

func prime_finder(done chan bool, steam <-chan int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for n := range or_done(done, steam) {
			if is_prime(n) {
				ch <- n
			}
		}
	}()
	return ch
}

func fanin[T any, K any](done chan K, streams ...<-chan T) <-chan T {
	relaych := make(chan T)
	wg := &sync.WaitGroup{}
	for _, s := range streams {
		wg.Add(1)
		go func(s <-chan T) {
			defer wg.Done()
			for n := range or_done(done, s) {
				relaych <- n
			}
		}(s)
	}
	go func() {
		defer close(relaych)
		for {
			select {
			case <-done:
				return
			default:
				wg.Wait()
			}
		}
	}()
	return relaych
}

func take[T any, K any](done chan K, stream <-chan T, n int) <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		count := 0
		for data := range or_done(done, stream) {
			if count < n {
				count++
				ch <- data
			} else {
				return
			}
		}
	}()
	return ch
}

func main() {
	fmt.Println("CPU Cores: ", runtime.NumCPU())
	randomn := func() int { return rand.Intn(900000000) }
	done := make(chan bool)
	defer close(done)
	start := time.Now()

	step1 := generator(done, randomn)

	// step2: normal Elapsed:  48.436571153s
	// step2 := prime_finder(done, step1)

	// step2: fan out Elapsed:  4.871425703s
	prime_streams := make([]<-chan int, runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		prime_streams = append(prime_streams, prime_finder(done, step1))
	}
	step2 := fanin(done, prime_streams...)

	step3 := take(done, step2, 5)

	for i := range step3 {
		fmt.Println(i)
	}

	fmt.Println("time elapsed: ", time.Since(start))
}
