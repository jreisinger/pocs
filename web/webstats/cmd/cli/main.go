package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"webstats"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix(filepath.Base(os.Args[0]) + ": ")

	stats := webstats.Get(os.Args[1:])
	for _, stat := range stats {
		if stat.Err != nil {
			log.Print(stat.Err)
			continue
		}
		fmt.Println(stat)
	}
}
