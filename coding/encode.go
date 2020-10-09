// To execute Go code, please declare a func main() in a package "main"

package main

import (
	"fmt"
	"strconv"
	"time"
)

//https://en.wikipedia.org/wiki/Run-length_encoding

// Input: aaa
// Output: 3a

// Input: aaabaaacd
// Output: 3ab3acd

// encode encodes a given string.
func encode(input string) string {
	var output string
	n := len(input)
	for i := 0; i < n; i++ {
		count := 1
		for ; i < n-1 && input[i] == input[i+1]; i++ {
			count++
		}
		if count == 1 {
			output = output + fmt.Sprintf("%s", string(input[i]))
		} else {
			output = output + fmt.Sprintf("%d%s", count, string(input[i]))
		}
	}
	return output
}

func encodeCh(ch chan rune) string {
	ticker := time.NewTicker(2 * time.Minute)
	var input string
	for {
		select {
		case s := <-ch:
			input = input + string(s)
		case <-ticker.C:
			input = encode(input)
		}
	}
}

// decode decodes a given input string
func decode(input string) string {
	var output string
	n := len(input)
	for i := 0; i < n; i++ {
		if ok, num := isNumeric(string(input[i])); ok {
			nextChar := input[i+1]
			for j := 0; j < num; j++ {
				output = output + string(nextChar)
			}
			i = i + 2
		} else {
			output = output + string(input[i])
			i = i + 1
		}

	}
	return output
}

// #refactor
func isNumeric(i string) (bool, int) {
	switch i {
	case "2", "3", "4", "5", "6", "7", "8", "9":
		num, _ := strconv.Atoi(i)
		return true, num
	default:
		return false, 0
	}
}
