package main

import(
	"fmt"
)

func main() {
	//scope of x is limited to if
	if x := 42; x == 2{
		fmt.Printf("Here")
	}

	switch {
	case true: fmt.Println("Case 1")
				fallthrough
	case true: fmt.Println("Case 2")
	case true: fmt.Println("Case 3")
	}
}