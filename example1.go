package main

import "fmt"
import "runtime"

//global scope, can't use short hand declaration here
var y = 2

//assign 0 for int, 0.0 for float, "" for string, nil for pointers, false for bool
var z int

type t int

var s t

//untyped constant
const (
	b = "Hello"
)

//typed constant
const c int = 34

//byte = unit8
//rune = int32

const (
	d = iota //0
	e        //1
	f        //2
)

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

	//runtime package
	fmt.Println(runtime.GOOS)
	fmt.Println(runtime.GOARCH)

	// can't change type, go is static, not dynamic
	//x = "sdsd"

	//conversion: convert value of one type to another
	x = int(s)
	//can't do x = s, even though underlying type is same
	//no casting in go

	//string data type
	fmt.Printf("%v\t%T\n", []byte("Hello"), []byte("Hello"))
	//print the UTF-8 codepoint
	fmt.Printf("%#U\n", m[0])

	for i, v := range m {
		fmt.Printf("at index %v, hex value %#x\n", i, v)
	}
}

func foo() {
	y = 21
	fmt.Println("Hello", y)
}
