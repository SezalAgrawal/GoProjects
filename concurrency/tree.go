package main

// Usecase: check if 2 binary trees have same in-order traversal

import (
	"fmt"

	"golang.org/x/tour/tree"
)

// recWalk walks recursively through the tree and push
// values to the channel at each recursion
func recWalk(t *tree.Tree, ch chan int) {
	if t != nil {
		// send the left part of the tree to be iterated over first
		recWalk(t.Left, ch)
		// push the value to the channel
		ch <- t.Value
		// send the right part of the tree to be iterated over last
		recWalk(t.Right, ch)
	}
}

// walk walks the tree t sending all values
// from the tree to the channel ch
func walk(t *tree.Tree, ch chan int) {
	recWalk(t, ch)
	// close the channel so that range can finish
	close(ch)
}

// same determines whether the trees
// t1 and t2 contain the same values
func same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go walk(t1, ch1)
	go walk(t2, ch2)

	for {
		x1, ok1 := <-ch1
		x2, ok2 := <-ch2
		switch {
		case ok1 != ok2: 
			// different tree size
			return false
		case !ok1:
			// both channels are empty
			return true
		case x1 != x2:
			// values are different
			return false
		}
	}
	
}

func main() {
	// print the tree
	ch := make(chan int)
	go walk(tree.New(1), ch)
	for v := range ch {
		fmt.Println(v)
	}

	// test if same
	fmt.Println(same(tree.New(1), tree.New(1)))
	fmt.Println(same(tree.New(1), tree.New(2)))
}
