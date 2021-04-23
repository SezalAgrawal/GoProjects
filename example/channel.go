package main

import (
	"context"
	"fmt"
	"time"
)

func test() bool {
	time.Sleep(5 * time.Second)
	return true
}

func main() {
	timeoutContext, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	checkResponse := make(chan bool, 1)
	
	go func(checkResponse chan<- bool) {
		checkResponse <- test()
	}(checkResponse)


	select {
	case approved := <-checkResponse:
		if approved {
			fmt.Println("here")
		}
	case <-timeoutContext.Done():
		fmt.Println("closed")
	}
}
