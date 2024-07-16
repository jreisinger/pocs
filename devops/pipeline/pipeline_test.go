package pipeline_test

import (
	"bytes"
	"errors"
	"io"
	"log"
	"pipeline"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFromFile_ReadsAllDataFromFile(t *testing.T) {
	t.Parallel()
	want := []byte("Hello, world.\n")
	p := pipeline.FromFile("testdata/hello.txt")
	if p.Error != nil {
		t.Fatal(p.Error)
	}
	got, err := io.ReadAll(p.Input)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestFromFile_SetsErrorGivenNonexistentFile(t *testing.T) {
	t.Parallel()
	p := pipeline.FromFile("doesnt-exist.txt")
	if p.Error == nil {
		t.Fatal("want error opening non-existent file, got nil")
	}
}

func TestStdoutPrintsMessageToOutput(t *testing.T) {
	t.Parallel()
	want := "Hello, world.\n"
	p := pipeline.FromString(want)
	buf := new(bytes.Buffer)
	p.Output = buf
	p.Stdout()
	if p.Error != nil {
		t.Fatal(p.Error)
	}
	got := buf.String()
	if !cmp.Equal(want, got) {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestStdoutPrintsNothingOnError(t *testing.T) {
	t.Parallel()
	p := pipeline.FromString("Hello, world.\n")
	p.Error = errors.New("oh no")
	buf := new(bytes.Buffer)
	p.Output = buf
	p.Stdout()
	got := buf.String()
	if got != "" {
		t.Errorf("want no output from Stdout afer error, but got %q", got)
	}
}

func TestColumnExtractsSecondColumn(t *testing.T) {
	t.Parallel()
	input := "1 2 3\n1 2 3\n1 2 3\n"
	p := pipeline.FromString(input)
	want := "2\n2\n2\n"
	got, err := p.Column(2).String()
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Errorf(cmp.Diff(want, got))
	}
}

func TestColumnProducesNothingWhenPipeErrorSet(t *testing.T) {
	t.Parallel()
	p := pipeline.FromString("1 2 3\n")
	p.Error = errors.New("oh no")
	data, err := io.ReadAll(p.Column(2).Input)
	if err != nil {
		log.Fatal(err)
	}
	if len(data) > 0 {
		t.Errorf("want no output from Column after error, but got %q", data)
	}
}

func TestStringReturnsPipeContents(t *testing.T) {
	t.Parallel()
	want := "Hello, world.\n"
	p := pipeline.FromString(want)
	got, err := p.String()
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestStringReturnsErrorWhenPipeErrorSet(t *testing.T) {
	t.Parallel()
	p := pipeline.FromString("Hello, world.\n")
	p.Error = errors.New("oh no")
	_, err := p.String()
	if err == nil {
		t.Error("want error from String when pipeline has error, but got nil")
	}
}

func TestColumnSetsErrorAndProducesNothingGivenInvalidArg(t *testing.T) {
	t.Parallel()
	p := pipeline.FromString("1\n2\n3\n")
	p.Column(-1)
	if p.Error == nil {
		t.Error("want error on non-positive Column arg, but got nil")
	}
	data, err := io.ReadAll(p.Column(1).Input)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) > 0 {
		t.Errorf("want no output from Column with invalid arg, but got %q", data)
	}
}
