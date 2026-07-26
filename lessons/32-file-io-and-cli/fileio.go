// Package fileio demonstrates the standard file-I/O patterns: whole-
// file read/write shortcuts, line-by-line streaming with bufio.Scanner,
// and the io.Reader interface as a testable function parameter.
package fileio

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// CountLines counts the newline-terminated lines in r. Takes an
// io.Reader (not a filename or *os.File) so callers can pass anything
// that reads bytes — an open file, a strings.Reader in tests, an
// HTTP response body, a gzip.Reader, etc.
func CountLines(r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)
	count := 0
	for scanner.Scan() {
		count++
	}
	// Errors surface here, not from Scan(). Easy to miss.
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("scan: %w", err)
	}
	return count, nil
}

// CountWords counts the whitespace-separated words in r. Uses
// bufio.ScanWords to override the default line-splitting behaviour.
func CountWords(r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	count := 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("scan: %w", err)
	}
	return count, nil
}

// WriteLinesToFile writes each string in lines to path, one per line.
// Uses os.WriteFile — the whole-file shortcut. Not appropriate for
// large slices; for those, use os.Create + a bufio.Writer.
func WriteLinesToFile(path string, lines []string) error {
	content := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// ReadWholeFile is the "small file" shortcut. Reads everything into
// memory at once. Prefer streaming for anything that could be big.
func ReadWholeFile(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return b, nil
}

// CountLinesInFile is the streaming version — opens a file and counts
// its lines without loading everything into memory. Suitable for a log
// file that could be gigabytes.
func CountLinesInFile(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	return CountLines(f)
}
