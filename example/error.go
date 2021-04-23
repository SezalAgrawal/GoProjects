package main

import (
	"errors"
	"fmt"
)

var (
	errRecordNotFound = errors.New("record not found")
	errConflict       = errors.New("duplicate record")
)

type splitPayError struct {
	description string
	error
}

func (e *splitPayError) Description() string {
	return e.description
}

func (e *splitPayError) Error() string {
	return fmt.Sprintf("%s: %v", e.description, e.error)
}

func (e *splitPayError) Unwrap() error {
	return e.error
}

func test1() *splitPayError {
	return &splitPayError{"error from test1 code", errRecordNotFound}
}

func test2() *splitPayError {
	if err := test1(); err != nil {
		return &splitPayError{fmt.Sprintf("error from test2 code, %v", err.Error()), errConflict}
	}
	return nil
}

func test3() *splitPayError {
	if err := test1(); err != nil {
		return &splitPayError{"error from test3 code", err}
	}
	return nil
}

func main() {
	// err := test3()
	// matched := errors.Is(err, errRecordNotFound)
	// fmt.Println(err.Error())
	// fmt.Println(matched)

	// err = test2()
	// matched = errors.Is(err, errConflict)
	// fmt.Println(err.Error())
	// fmt.Println(matched)
	type a struct {
		name string
	}
	type b struct {
		id    string
		names []a
	}

	bo := b{id: "123"}
	x := len(bo.names)
	fmt.Println(bo)
	fmt.Println(x)
	xo :=make([]a, x)
	bo.names=xo
	if bo.names == nil {
		fmt.Println("herex")
	}
	fmt.Println(x)
	var abc []int
	if len(abc) == 0 {
		fmt.Println("here")
	}
	abc = append(abc, 2)
	fmt.Println(x)
}
