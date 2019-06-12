package main

import "fmt"

//global scope, can't use short hand declaration here
var y = 2

//assign 0 for int, 0.0 for float, "" for string, nil for pointers, false for bool
var z int

func main() {
	x := 3
	m := `She said,
	 "Hello"`
	fmt.Println("Hello", y, x)
	foo()
	fmt.Println("Hello", y, x)
	fmt.Printf("%T\n", z)
	fmt.Println(m)

	// can't change type
	//x = "sdsd"
}

func foo() {
	y = 21
	fmt.Println("Hello", y)
}