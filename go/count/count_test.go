package count_test

import (
	"bytes"
	"count"
	"io"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestCountLines_CountsLinesCorrectly(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input io.Reader
		want  int
	}{
		{bytes.NewBufferString(""), 0},
		{bytes.NewBufferString("\n"), 1},
		{bytes.NewBufferString("1\n2\n3"), 3},
		{bytes.NewBufferString("1\n2\n3\n"), 3},
	}

	for i, test := range tests {
		c, err := count.NewCounter(count.WithInput(test.input))
		if err != nil {
			t.Fatal(err)
		}
		got := c.CountLines()
		if test.want != got {
			t.Errorf("test %d: line count should be %d but is %d", i+1, test.want, got)
		}
	}
}

func TestCountWords_CountsWordsCorrectly(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input io.Reader
		want  int
	}{
		{bytes.NewBufferString(""), 0},
		{bytes.NewBufferString("1"), 1},
		{bytes.NewBufferString("1 2 3"), 3},
		{bytes.NewBufferString("1\n2\n3\n"), 3},
	}
	for _, test := range tests {
		c, err := count.NewCounter(count.WithInput(test.input))
		if err != nil {
			t.Fatal(err)
		}
		got := c.CountWords()
		if test.want != got {
			t.Errorf("word count should be %d but is %d", test.want, got)
		}
	}
}

func TestCountLines_WithInputFromArgs(t *testing.T) {
	t.Parallel()
	c, err := count.NewCounter(count.WithInputFromArgs([]string{"testdata/three_lines.txt"}))
	if err != nil {
		t.Fatal(err)
	}
	want := 3
	got := c.CountLines()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"lines": count.MainLines,
		"words": count.MainWords,
		"count": count.Main,
	})
}

func Test(t *testing.T) {
	t.Parallel()
	testscript.Run(t, testscript.Params{
		Dir: "testdata/script",
	})
}
