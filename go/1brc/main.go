package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

var maxGoroutines = runtime.NumCPU()

func main() {
	log.SetPrefix(os.Args[0] + ": ")
	log.SetFlags(0)

	var p = flag.Bool("p", false, "run the parallelised version")
	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatal("please supply measurements file")
	}
	file := flag.Args()[0]

	var process func(string, io.Writer) error = simple
	if *p {
		process = parallel
	}
	if err := process(file, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

type stats struct {
	min, max, sum float64
	count         int64 // mean (or average) = sum / count
}

// simple is a simple and idiomatic Go code to process the measurements file.
func simple(filename string, output io.Writer) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	stationStats := make(map[string]stats)

	input := bufio.NewScanner(f)
	for input.Scan() {
		line := input.Text()
		station, tempStr, hasSemi := strings.Cut(line, ";")
		if !hasSemi {
			continue
		}

		temp, err := strconv.ParseFloat(tempStr, 64)
		if err != nil {
			return err
		}

		s, seen := stationStats[station]
		if !seen {
			s.min = temp
			s.max = temp
			s.sum = temp
			s.count = 1
		} else {
			s.min = min(s.min, temp)
			s.max = max(s.max, temp)
			s.sum += temp
			s.count++
		}
		stationStats[station] = s
	}

	stations := make([]string, 0, len(stationStats))
	for station := range stationStats {
		stations = append(stations, station)
	}
	sort.Strings(stations)

	fmt.Fprint(output, "{")
	for i, station := range stations {
		if i > 0 {
			fmt.Fprint(output, ", ")
		}
		s := stationStats[station]
		mean := s.sum / float64(s.count)
		fmt.Fprintf(output, "%s=%.1f/%.1f/%.1f", station, s.min, mean, s.max)
	}
	fmt.Fprint(output, "}\n")

	return nil
}

type part struct {
	offset, size int64
}

func processPart(fileName string, fileOffset, fileSize int64, resultsCh chan map[string]stats) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Seek(fileOffset, io.SeekStart)
	if err != nil {
		panic(err)
	}
	f := io.LimitedReader{R: file, N: fileSize}

	stationStats := make(map[string]stats)

	scanner := bufio.NewScanner(&f)
	for scanner.Scan() {
		line := scanner.Text()
		station, tempStr, hasSemi := strings.Cut(line, ";")
		if !hasSemi {
			continue
		}

		temp, err := strconv.ParseFloat(tempStr, 64)
		if err != nil {
			panic(err)
		}

		s, ok := stationStats[station]
		if !ok {
			s.min = temp
			s.max = temp
			s.sum = temp
			s.count = 1
		} else {
			s.min = min(s.min, temp)
			s.max = max(s.max, temp)
			s.sum += temp
			s.count++
		}
		stationStats[station] = s
	}

	resultsCh <- stationStats
}

// parallel parallelises the simple function. This is a map-reduce problem so we
// split the file into similar-sized chunks (one for each CPU core) and process
// each part in a single thread (in Go, a goroutine). Then we merge the results
// together at the end.
func parallel(filename string, output io.Writer) error {
	parts, err := splitFile(filename, maxGoroutines)
	if err != nil {
		return err
	}

	resultsCh := make(chan map[string]stats)
	for _, part := range parts {
		go processPart(filename, part.offset, part.size, resultsCh)
	}

	totals := make(map[string]stats)
	for i := 0; i < len(parts); i++ {
		result := <-resultsCh
		for station, s := range result {
			ts, seen := totals[station]
			if !seen {
				totals[station] = stats{
					min:   s.min,
					max:   s.max,
					sum:   s.sum,
					count: s.count,
				}
				continue
			}
			ts.min = min(ts.min, s.min)
			ts.max = min(ts.max, s.max)
			ts.sum += s.sum
			ts.count += s.count
			totals[station] = ts
		}
	}

	stations := make([]string, 0, len(totals))
	for station := range totals {
		stations = append(stations, station)
	}
	sort.Strings(stations)

	fmt.Fprint(output, "{")
	for i, station := range stations {
		if i > 0 {
			fmt.Fprint(output, ", ")
		}
		s := totals[station]
		mean := s.sum / float64(s.count)
		fmt.Fprintf(output, "%s=%.1f/%.1f/%.1f", station, s.min, mean, s.max)
	}
	fmt.Fprint(output, "}\n")

	return nil
}

func splitFile(file string, numParts int) ([]part, error) {
	const maxLineLength = 100

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	st, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := st.Size()
	splitSize := size / int64(numParts)

	buf := make([]byte, maxLineLength)

	parts := make([]part, 0, numParts)
	offset := int64(0)
	for offset < size {
		seekOffset := max(offset+splitSize-maxLineLength, 0)
		if seekOffset > size {
			break
		}
		_, err := f.Seek(seekOffset, io.SeekStart)
		if err != nil {
			return nil, err
		}
		n, _ := io.ReadFull(f, buf)
		chunk := buf[:n]
		newline := bytes.LastIndexByte(chunk, '\n')
		if newline < 0 {
			return nil, fmt.Errorf("newline not found at offset %d", offset+splitSize-maxLineLength)
		}
		remaining := len(chunk) - newline - 1
		nextOffset := seekOffset + int64(len(chunk)) - int64(remaining)
		parts = append(parts, part{offset, nextOffset - offset})
		offset = nextOffset
	}
	return parts, nil
}
