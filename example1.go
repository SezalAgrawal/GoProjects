package main

import "fmt"

//global scope, can't use short hand declaration here
var y = 2

//assign 0 for int, 0.0 for float, "" for string, nil for pointers, false for bool
var z int
type t int
var s t

func main() {
	x := 3
	s = 5
	m := `She said,
	 "Hello"`
	fmt.Println("Hello", y, x)
	foo()
	fmt.Println("Hello", y, x)
	fmt.Printf("%T\n", z)
	fmt.Println(m)
	fmt.Printf("%d\n%x\n%#x", y, y, y)
	p := fmt.Sprintf("%d\n%x\n%#x", y, y, y)
	fmt.Println(p)
	fmt.Printf("%T\n", s)

	// can't change type, go is static, not dynamic
	//x = "sdsd"

	//conversion: convert value of one type to another
	x = int(s)
	//can't do x = s, even though underlying type is same
	//no casting in go
}

func foo() {
	y = 21
	fmt.Println("Hello", y)
}