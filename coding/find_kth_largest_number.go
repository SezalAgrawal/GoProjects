// find a number k such that k and -k both exists and k is the greatest number

package main

import (
	"fmt"
	"sort"
)

func main() {
	arr := []int{1, 1, 2, -1, 2, -1}
	fmt.Println(diff(arr))
}

func diff(arr []int) int {
	sort.Ints(arr)
	for i, n := 0, len(arr)-1; i <= len(arr)/2 && n >= len(arr)/2; {
		left := mod(arr[i])
		right := mod(arr[n])
		if left == right {
			return left
		} else if left > right {
			i++
		} else if right > left {
			n--
		}
	}
	return 0
}

func mod(i int) int {
	if i < 0 {
		return -1 * i
	}
	return i
}
