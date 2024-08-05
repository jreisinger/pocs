package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

type keepFunc[E any] func(E) bool

func Filter[S ~[]E, E any](s S, f keepFunc[E]) S {
	result := S{}
	for _, v := range s {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

func IsEven[T constraints.Integer](v T) bool {
	return v%2 == 0
}

func main() {
	s := []int{1, 2, 3, 4}
	fmt.Println(Filter(s, IsEven[int]))
}
