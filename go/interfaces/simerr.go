package main

import (
	"errors"
	"math/rand"
)

func simulateError() error {
	if rand.Intn(10) < 3 {
		return errors.New("some error")
	}
	return nil
}
