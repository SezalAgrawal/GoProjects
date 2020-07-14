package main

// Generates the execution steps in assembly language
// go tool compile -S object_dump.go > object_dump.s

func example() {
	panic("Panicking!!")
}

func main() {
	example()
}