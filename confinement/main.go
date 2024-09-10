package main

import (
	"fmt"
	"sync"
	"time"
)

func multiple(wg *sync.WaitGroup, n int, res_arr *[]int, multiplier int) {
	defer wg.Done()
	time.Sleep(time.Second)
	*res_arr = append(*res_arr, n*multiplier)
}

func multiple_comfinement(wg *sync.WaitGroup, n int, resAddr *int, multiplier int) {
	defer wg.Done()
	time.Sleep(time.Second)
	*resAddr = n * multiplier
}

func main() {
	wg := &sync.WaitGroup{}
	inputs := []int{2, 4, 6, 8, 3}
	outputs := make([]int, 5)

	for i := range len(outputs) {
		wg.Add(1)

		// changing array - resulted in race condition
		// go multiple(wg, inputs[i], &outputs, 2)

		// confinement - concurrency safe
		go multiple_comfinement(wg, inputs[i], &outputs[i], 2)
	}

	wg.Wait()
	fmt.Println(outputs)
}
