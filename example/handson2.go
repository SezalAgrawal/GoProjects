package main

import (
	"fmt"
)

//create typed and untyped constants
const (
	y int = 3
	z     = 4
)

const (
	a = 2017 + iota
	b
	c
	d
)

func main() {
	x := 3
	//Print decimal, binary, hex
	fmt.Printf("%d, %b, %#x\n", x, x, x)
	fmt.Println(a, b, c, d)

}

//Notes
//Use int, float64 generic
