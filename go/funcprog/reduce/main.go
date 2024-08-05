package main

import "fmt"

type reduceFunc[E any] func(cur, next E) E

func Reduce[E any](s []E, init E, f reduceFunc[E]) E {
	cur := init
	for _, v := range s {
		cur = f(cur, v)
	}
	return cur
}

func main() {
	s := []int{1, 2, 3, 4}
	sum := Reduce(s, 0, func(cur, next int) int {
		return cur + next
	})
	fmt.Println(sum)
}
