package main

import "fmt"

// node is an entity in a linked list
type node struct {
	val  int
	next *node
}

// reverse reverses the given linked list
func (n *node) reverse() {
	var previous, next *node
	current := n
	for current != nil {
		next = current.next
		current.next = previous
		previous = current
		current = next
	}
	n = previous
}

// prints the elements in the linked list
func (n *node) print() {
	current := n
	for current != nil {
		fmt.Println(current.val)
		current = current.next
	}
}

func main() {
	n := new(node)
	n.val = 1
	n.next = &node{val: 2, next: &node{val: 3}}
	n.reverse()
	n.print()
}
