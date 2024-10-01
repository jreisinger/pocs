package main

import (
	"reflect"
	"testing"
)

func TestChunkSlice(t *testing.T) {
	testcases := []struct {
		slice   []string
		maxSize int
		want    [][]string
	}{
		{[]string{"a", "b", "c", "d", "e"}, 2, [][]string{{"a", "b"}, {"c", "d"}, {"e"}}},
		{},
		{[]string{""}, 0, nil},
		{[]string{"a"}, 0, nil},
	}
	for _, tc := range testcases {
		chunks := chunkSlice(tc.slice, tc.maxSize)
		for _, chunk := range chunks {
			for _, chunkWant := range tc.want {
				if !reflect.DeepEqual(chunk, chunkWant) {
					t.Errorf("got %v, want %v", chunk, chunkWant)
				}
			}
		}
	}
}
