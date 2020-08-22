package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

var files []string

func init() {
	rand.Seed(time.Now().UnixNano())
	files = generateRandomFiles(1e3)
	fmt.Printf("Processing %d files using %d goroutines\n", len(files), runtime.NumCPU())
}

func BenchmarkSequential(b *testing.B) {
	for i := 0; i < b.N; i++ {
		findSequential("test", files)
	}
}

func BenchmarkConcurrent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		findConcurrentWithChannel(runtime.NumCPU(), "test", files)
	}
}

func BenchmarkSequentialAgain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		findSequential("test", files)
	}
}

func BenchmarkConcurrentAgain(b *testing.B) {
	for i := 0; i < b.N; i++ {
		findConcurrentWithChannel(runtime.NumCPU(), "test", files)
	}
}
