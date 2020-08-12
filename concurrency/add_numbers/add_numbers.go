// This example tests a CPU bound process, where voluntary context switches are less.
// To maximize performance, work can be done concurrently with parallelism.
// Number of goroutines should be equal to the logical proccessors of the system.
package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateRandomNumbers(count int) []int {
	numbers := make([]int, count)
	for i := 0; i < count; i++ {
		numbers[i] = rand.Intn(count)
	}
	return numbers
}

func addSequential(numbers []int) int {
	var sum int
	for _, number := range numbers {
		sum = sum + number
	}
	return sum
}

func addConcurrent(goroutines int, numbers []int) int {
	var sum int64
	split := len(numbers) / goroutines

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			start := i * split
			end := start + split
			if i == (goroutines - 1) {
				end = len(numbers)
			}
			var temp int
			for _, number := range numbers[start:end] {
				temp += number
			}
			atomic.AddInt64(&sum, int64(temp))
		}(i)
	}

	wg.Wait()
	return int(sum)
}

func main() {
	numbers := generateRandomNumbers(1e7)
	fmt.Println(addSequential(numbers))
	fmt.Println(addConcurrent(runtime.NumCPU(), numbers))
}
