package main

import(
	"fmt"
)

func main() {
	//for loop to print ascii
	for i:=33; i<122; i++ {
		fmt.Printf("%v\t%#U", i, i)
	}
}