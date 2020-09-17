package main

// testing profiling
// go tool pprof http://localhost:6060/debug/pprof/allocs
// curl -sK -v http://localhost:6060/debug/pprof/allocs > allocs.out
// go tool pprof -http=:6061 allocs.out

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
	_ "net/http/pprof"
)

func main() {
	// we need a webserver to get the pprof webserver
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", http.DefaultServeMux))
	}()
	fmt.Println("hello world")
	var wg sync.WaitGroup
	wg.Add(1)
	go leakyFunction(wg)
	wg.Wait()
}

func leakyFunction(wg sync.WaitGroup) {
	defer wg.Done()
	s := make([]string, 3)
	for i := 0; i < 10000000; i++ {
		s = append(s, "magical pandas")
		if (i % 100000) == 0 {
			time.Sleep(500 * time.Millisecond)
		}
	}
}
