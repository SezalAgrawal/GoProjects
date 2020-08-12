package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

var numbers []int

func init() {
	rand.Seed(time.Now().UnixNano())
	numbers = generateRandomNumbers(1e7)
	fmt.Printf("Processing %d numbers using %d goroutines\n", len(numbers), runtime.NumCPU())
}

func BenchmarkSequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		addSequential(numbers)
	}
}

func BenchmarkConcurrent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		addConcurrent(runtime.NumCPU(), numbers)
	}
}

func BenchmarkSequentialAgain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		addSequential(numbers)
	}
}

func BenchmarkConcurrentAgain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		addConcurrent(runtime.NumCPU(), numbers)
	}
}

// Observations:

// 1. concurrency without parallelism
// --------------------------------
// ▶ GOGC=off go test -cpu 1 -run none -bench . -benchtime 3s
// Processing 10000000 numbers using 4 goroutines
// goos: darwin
// goarch: amd64
// pkg: github.com/goProjects/concurrency/add_numbers
// BenchmarkSequential                  564           7017853 ns/op
// BenchmarkConcurrent                  444           9672264 ns/op
// BenchmarkSequentialAgain             552           6398380 ns/op
// BenchmarkConcurrentAgain             492           7722396 ns/op
// PASS
// ok      github.com/goProjects/concurrency/add_numbers   25.367s
// ---------------------------------
// Concurrent takes more time as add_numbers is a cpu bound process
// and lot of context switch overhead takes place as there is only single thread.

// 2. concurrency with parallelism
// --------------------------------
// ▶ GOGC=off go test -run none -bench . -benchtime 3s       
// Processing 10000000 numbers using 4 goroutines
// goos: darwin
// goarch: amd64
// pkg: github.com/goProjects/concurrency/add_numbers
// BenchmarkSequential-4                561           6213913 ns/op
// BenchmarkConcurrent-4               1034           3231111 ns/op
// BenchmarkSequentialAgain-4           562           6507637 ns/op
// BenchmarkConcurrentAgain-4          1107           3410508 ns/op
// PASS
// ok      github.com/goProjects/concurrency/add_numbers   17.052s
// ---------------------------------
// Concurrent takes very less time than sequential as 4 goroutines are running
// parallely and thus completing their work concurrently.
