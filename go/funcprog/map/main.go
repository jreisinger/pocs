package main

import (
	"fmt"
	"strings"
)

type mapFunc[E any] func(E) E

func Map[S ~[]E, E any](s S, f mapFunc[E]) S {
	result := make(S, len(s))
	for i := range s {
		result[i] = f(s[i])
	}
	return result
}

func main() {
	s := []string{"a", "b", "c"}
	fmt.Println(Map(s, strings.ToUpper))
}
