// The example tests an IO bound process, where voluntary context switches are more
// as process goes into waiting stage.
package main

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func generateRandomFiles(count int) []string {
	files := make([]string, count)
	for i := 0; i < count; i++ {
		files[i] = "abc.txt"
	}
	return files
}

func readFile(file string) ([]item, error) {
	// assuming reading from the file
	time.Sleep(time.Millisecond) // simulate blocking disk read
	d := new(doc)
	if err := json.Unmarshal([]byte(file), d); err != nil {
		return nil, err
	}
	return d.Items, nil
}

func findSequential(text string, files []string) int {
	var count int
	for _, file := range files {
		items, err := readFile(file)
		if err != nil {
			continue
		}
		for _, item := range items {
			if strings.Contains(item.Description, text) {
				count++
			}
		}
	}
	return count
}

func findConcurrentWithoutChannel(goroutines int, text string, files []string) int {
	var count int64
	split := len(files) / goroutines

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			start := i * split
			end := start + split
			if i == goroutines-1 {
				end = len(files)
			}
			for _, file := range files[start:end] {
				items, err := readFile(file)
				if err != nil {
					continue
				}
				for _, item := range items {
					if strings.Contains(item.Description, text) {
						atomic.AddInt64(&count, 1)
					}
				}
			}
		}(i)
	}
	wg.Wait()
	return int(count)
}

// This is a better way because some files can be larger and
// thus dividing files equally won't be efficient.
func findConcurrentWithChannel(goroutines int, text string, files []string) int {
	var count int64

	ch := make(chan string, len(files))
	for _, file := range files {
		ch <- file
	}
	close(ch)

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			var lcount int64
			for file := range ch {
				items, err := readFile(file)
				if err != nil {
					continue
				}
				for _, item := range items {
					if strings.Contains(item.Description, text) {
						lcount++
					}
				}
			}
			atomic.AddInt64(&count, lcount)
		}()
	}
	wg.Wait()
	return int(count)
}

func main() {
	files := generateRandomFiles(1e3)
	fmt.Println(findSequential("go", files))
	fmt.Println(findConcurrentWithoutChannel(runtime.NumCPU(), "go", files))
	fmt.Println(findConcurrentWithChannel(runtime.NumCPU(), "go", files))
}

type item struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type doc struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Items       []item `json:"items"`
}

var file = `{
	"title": "GO programming",
	"description": "test description",
	"items": [
		{
			"title": "concurrency",
			"description": "go routines are used"
		}
	]
}`
