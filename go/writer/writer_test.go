package writer_test

import (
	"io/fs"
	"os"
	"testing"

	"writer"

	"github.com/google/go-cmp/cmp"
)

func TestWriteToFile_WritesGivenDataToFile(t *testing.T) {
	t.Parallel()

	path := t.TempDir() + "/write_test.txt"

	want := []byte{1, 2, 3}
	err := writer.WriteToFile(path, want)
	if err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestWriteToFile_ReturnsErrorForUnwritableFile(t *testing.T) {
	t.Parallel()

	path := "nonexistentdir/write_test.txt"
	want := []byte{1, 2, 3}
	err := writer.WriteToFile(path, want)
	if err == nil {
		t.Fatal("want error when file not writable")
	}
}

func TestWriteToFile_ClobbersExistingFile(t *testing.T) {
	t.Parallel()

	path := t.TempDir() + "/write_test.txt"

	if err := os.WriteFile(path, []byte{4, 5, 6}, 0o600); err != nil {
		t.Fatal(err)
	}

	want := []byte{1, 2, 3}
	if err := writer.WriteToFile(path, want); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestWriteToFile_WritesNewFileWithCorrectPermissions(t *testing.T) {
	t.Parallel()

	want := fs.FileMode(0o600)

	path := t.TempDir() + "/write_test.txt"
	if err := writer.WriteToFile(path, []byte{}); err != nil {
		t.Fatal(err)
	}

	stat, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	got := stat.Mode().Perm()

	if want != got {
		t.Errorf("want permissions 0%o, got 0%o", want, got)
	}
}

func TestWriteToFile_WritesExistingFileWithCorrectPermissions(t *testing.T) {
	t.Parallel()

	want := fs.FileMode(0o600)

	path := t.TempDir() + "/write_test.txt"
	if err := os.WriteFile(path, []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}

	if err := writer.WriteToFile(path, []byte{}); err != nil {
		t.Fatal(err)
	}

	stat, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	got := stat.Mode().Perm()

	if want != got {
		t.Errorf("want permissions 0%o, got 0%o", want, got)
	}
}
