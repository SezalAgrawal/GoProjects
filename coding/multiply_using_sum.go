package main

import (
	"fmt"
)

type square struct {
	side int
}

func (s square) perimeter() int {
	return 2 * sum(s.side, s.side)
}

func sum(a, b int) int {
	return a + b
}

func multiply(a, b int) int {
	// save the sign
	sign := 1
	if a < 0 && b > 0 || a > 0 && b < 0 {
		sign = -1
	}
	if a < 0 {
		a *= -1
	}
	if b < 0 {
		b *= -1
	}
	// get the smaller element
	smallerValue := b
	largerValue := a
	if a < b {
		smallerValue = a
		largerValue = b
	}
	var multiplyValue int
	for x := smallerValue; x > 0; x-- {
		multiplyValue = sum(multiplyValue, largerValue)
	}
	return multiplyValue * sign
}

func multiplyConcurrent(a, b int) int {
	// get the smaller element
	smallerValue := b
	if a < b {
		smallerValue = a
	}
	var multiplyValue int
	for x := smallerValue; x > 0; x-- {
		multiplyValue = sum(multiplyValue, a)
	}
	return multiplyValue
}

func main() {
	// a := 1
	// sq := square{
	// 	side: a,
	// }
	// fmt.Println("perimeter", sq.perimeter())
	// fmt.Println("multiply", multiply(1, 8))
	fmt.Println("sum1", sum1(1, 2))
}

func sum1(a, n int) int {
	s := a
	for x := 1; x <= n; x++ {
		s = sum(s, s)
	}
	return s
}

// 1 1 1 1 1 1 1 1

// sum(1, 1) = 2

// 2*4

// 2 2 2 2

// sum(2, 2) = 4

// 4*2

// 4 4

// sum(4, 4) = 8

// 8
// 1000
// 3

// 9
// 1001
// 4

// 10
// 1010
// 4

// 11
// 1011
// 5

// 357
// 101100101

// 1 1 1 1 1 1 1 1 1 1

// sum(1, 1) = 2

// 2*5

// 2 2 2 2 2

// sum(2, 2) = 4

// 4*2 + 2

// 4 4

// sum(4, 4) = 8

// sum(8, 2)
