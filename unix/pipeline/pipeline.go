package pipeline

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

type Pipeline struct {
	Input  io.Reader
	Output io.Writer
	Error  error
}

func newPipeline() *Pipeline {
	return &Pipeline{
		Output: os.Stdout,
	}
}

func FromFile(name string) *Pipeline {
	file, err := os.Open(name)
	if err != nil {
		return &Pipeline{Error: err}
	}
	p := newPipeline()
	p.Input = file
	return p
}

func FromString(s string) *Pipeline {
	p := newPipeline()
	p.Input = strings.NewReader(s)
	return p
}

func (p *Pipeline) Column(n int) *Pipeline {
	if p.Error != nil {
		p.Input = strings.NewReader("")
		return p
	}
	if n < 1 {
		p.Error = fmt.Errorf("bad column %d: must be positive", n)
		return p
	}
	result := new(bytes.Buffer)
	input := bufio.NewScanner(p.Input)
	for input.Scan() {
		fields := strings.Fields(input.Text())
		if len(fields) < n {
			continue
		}
		fmt.Fprintln(result, fields[n-1])
	}
	return &Pipeline{
		Input: result,
	}
}

func (p *Pipeline) Stdout() {
	if p.Error != nil {
		return
	}
	_, err := io.Copy(p.Output, p.Input)
	if err != nil {
		p.Error = err
	}
}
func (p *Pipeline) String() (string, error) {
	if p.Error != nil {
		return "", p.Error
	}
	data, err := io.ReadAll(p.Input)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
