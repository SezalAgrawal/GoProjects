package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	test()
}

func test() {
	timeoutContext, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	test1(timeoutContext)
}

func test1(ctx context.Context) {
	timeoutContext, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	time.Sleep(1 * time.Second)

	select {
	case <-timeoutContext.Done():
		fmt.Println("here")
	case <-ctx.Done():
		fmt.Println("here main")
	}
}
