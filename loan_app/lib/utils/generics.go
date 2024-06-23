package utils

import "strconv"

func Ptr[T any](v T) *T {
	return &v
}

type Number interface {
	int | int64 | float64
}

func SumNumbers[T Number](inputs []T) T {
	var sum T
	for _, input := range inputs {
		sum += input
	}

	return sum
}

func SliceConvert[T, V any](inputs []T, converter func(T) V) []V {
	result := make([]V, len(inputs))
	for idx, input := range inputs {
		result[idx] = converter(input)
	}
	return result
}

func StringToIntConverter(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return val
}
