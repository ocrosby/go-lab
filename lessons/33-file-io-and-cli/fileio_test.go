package fileio_test

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	fileio "github.com/ocrosby/go-lab/lessons/33-file-io-and-cli"
)

func TestCountLines_CountsFromReader(t *testing.T) {
	// io.Reader-taking functions test against strings.Reader — no
	// filesystem needed, fast, deterministic.
	input := strings.NewReader("one\ntwo\nthree\n")

	n, err := fileio.CountLines(input)
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if n != 3 {
		t.Errorf("n = %d, want 3", n)
	}
}

func TestCountLines_EmptyReaderReturnsZero(t *testing.T) {
	n, err := fileio.CountLines(strings.NewReader(""))
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if n != 0 {
		t.Errorf("n = %d, want 0", n)
	}
}

func TestCountWords_UsesScanWordsSplit(t *testing.T) {
	input := strings.NewReader("the quick brown fox jumps")

	n, err := fileio.CountWords(input)
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if n != 5 {
		t.Errorf("n = %d, want 5", n)
	}
}

func TestWriteLinesToFile_RoundTrips(t *testing.T) {
	// t.TempDir() returns a fresh directory that gets cleaned up
	// automatically after the test. Filesystem isolation without
	// per-test bookkeeping.
	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")

	err := fileio.WriteLinesToFile(path, []string{"alpha", "beta", "gamma"})
	if err != nil {
		t.Fatalf("WriteLinesToFile err = %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile err = %v", err)
	}

	want := "alpha\nbeta\ngamma\n"
	if string(got) != want {
		t.Errorf("file contents = %q, want %q", string(got), want)
	}
}

func TestReadWholeFile_ReadsBackWritten(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "greeting.txt")
	_ = os.WriteFile(path, []byte("hello, world"), 0o644)

	got, err := fileio.ReadWholeFile(path)
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if string(got) != "hello, world" {
		t.Errorf("got %q, want %q", string(got), "hello, world")
	}
}

func TestReadWholeFile_ReportsNotExist(t *testing.T) {
	_, err := fileio.ReadWholeFile("/tmp/definitely-does-not-exist-xyz-999")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	// The underlying os error is os.ErrNotExist (an fs.ErrNotExist).
	// errors.Is walks the wrap chain and finds it.
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("err = %v, want to wrap fs.ErrNotExist", err)
	}
}

func TestCountLinesInFile_MatchesReaderCount(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "many.txt")
	_ = os.WriteFile(path, []byte("a\nb\nc\nd\ne\n"), 0o644)

	n, err := fileio.CountLinesInFile(path)
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if n != 5 {
		t.Errorf("n = %d, want 5", n)
	}
}
