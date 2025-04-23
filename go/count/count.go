package count

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type Counter struct {
	files []io.Reader
	input io.Reader
}

type option func(c *Counter) error

func WithInput(r io.Reader) option {
	return func(c *Counter) error {
		if r == nil {
			return fmt.Errorf("nil input")
		}
		c.input = r
		return nil
	}
}

func WithInputFromArgs(args []string) option {
	return func(c *Counter) error {
		if len(args) < 1 {
			return nil
		}
		var files []io.Reader
		for _, arg := range args {
			f, err := os.Open(arg)
			if err != nil {
				return err
			}
			files = append(files, f)
			c.files = append(c.files, f)
		}
		c.input = io.MultiReader(files...)
		return nil
	}
}

func NewCounter(opts ...option) (*Counter, error) {
	c := &Counter{input: os.Stdin}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (c *Counter) CountLines() int {
	var lines int
	input := bufio.NewScanner(c.input)
	for input.Scan() {
		lines++
	}
	for _, f := range c.files {
		f.(io.Closer).Close()
	}
	return lines
}

func (c *Counter) CountWords() int {
	var words int
	input := bufio.NewScanner(c.input)
	input.Split(bufio.ScanWords)
	for input.Scan() {
		words++
	}
	for _, f := range c.files {
		f.(io.Closer).Close()
	}
	return words
}

func MainLines() {
	c, err := NewCounter(WithInputFromArgs(os.Args[1:]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c.CountLines())
}

func MainWords() {
	c, err := NewCounter(WithInputFromArgs(os.Args[1:]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c.CountWords())
}

func Main() {
	log.SetFlags(0)
	log.SetPrefix(os.Args[0] + ": ")

	flag.Usage = func() {
		fmt.Printf("Counts words (or lines) from stdin (or files).\n\n")
		fmt.Printf("%s [-lines] [files...]\n", os.Args[0])
		flag.PrintDefaults()
	}
	lines := flag.Bool("lines", false, "count lines, not words")
	flag.Parse()

	c, err := NewCounter(WithInputFromArgs(flag.Args()))
	if err != nil {
		log.Fatal(err)
	}

	if *lines {
		fmt.Println(c.CountLines())
	} else {
		fmt.Println(c.CountWords())
	}
}
