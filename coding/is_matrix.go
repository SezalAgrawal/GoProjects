// program to check if a given square is magic square

package main

import (
	"golang.org/x/exp/errors/fmt"
)

func main() {
	// given square matrix
	matrix := [4][4]int{
		{16, 2, 3, 13},
		{5, 11, 10, 8},
		{9, 7, 6, 12},
		{4, 14, 15, 1},
	}

	// n is the length of the square matrix
	n := len(matrix)

	// take sum of the first row as a reference sum
	expectedSum := 0
	for j := 0; j < n; j++ {
		expectedSum += matrix[0][j]
	}

	// flag to check if the given sum is equal to expected sum
	isMagic := true
	// colSum is used to store column sum
	colSum := make([]int, n)

	// iterate over rows to check if ths row sum is same
	// and also store intermediary column sum
	for i := 0; i < n; i++ {
		rowSum := 0
		for j := 0; j < n; j++ {
			rowSum += matrix[i][j]
			colSum[j] += matrix[i][j]
		}
		if rowSum != expectedSum {
			isMagic = false
			break
		}
	}
	if !isMagic {
		handleNonMagicSquare()
		return
	}

	// loop through colSum to check if all the values are same
	for i := 0; i < n; i++ {
		if colSum[i] != expectedSum {
			handleNonMagicSquare()
			return
		}
	}

	// iterate over diagonals to check if ths sum is same
	leftDiagonalSum := 0
	rightDiagonalSum := 0
	for i, j := 0, 0; i < n && j < n; {
		leftDiagonalSum += matrix[i][j]
		rightDiagonalSum += matrix[i][n-1-j]
		i++
		j++
	}
	if rightDiagonalSum != expectedSum || leftDiagonalSum != expectedSum {
		handleNonMagicSquare()
		return
	}

	fmt.Println("A magic square!")
}

// handleNonMagicSquare can be customized to handle non magic case scenario
func handleNonMagicSquare() {
	fmt.Println("Not a magic square!")
}
