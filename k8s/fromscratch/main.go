package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	var numbers []int
	for _, arg := range os.Args[1:] {
		n, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fibspin: %v\n", err)
			os.Exit(1)
		}
		numbers = append(numbers, n)
	}
	go spinner(time.Millisecond * 100)
	c := make(chan string)
	for _, n := range numbers {
		go func(n int) {
			start := time.Now()
			f := fib(n)
			duration := time.Since(start)
			c <- fmt.Sprintf("\rfib(%d) = %d (%s)", n, f, duration)
		}(n)
	}
	for range numbers {
		fmt.Println(<-c)
	}
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `\|/-` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-2) + fib(n-1)
}
