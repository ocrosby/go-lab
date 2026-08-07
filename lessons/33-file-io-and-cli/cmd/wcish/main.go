// Command wcish is a small `wc`-like tool that counts lines, words, and
// bytes in one or more files. Demonstrates flag parsing, os.Args
// positional args, and delegating to the fileio helpers.
//
//	go run ./lessons/33-file-io-and-cli/cmd/wcish -lines /etc/hosts
//	go run ./lessons/33-file-io-and-cli/cmd/wcish -words -bytes README.md go.mod
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	fileio "github.com/ocrosby/go-lab/lessons/33-file-io-and-cli"
)

func main() {
	var (
		showLines = flag.Bool("lines", false, "count lines")
		showWords = flag.Bool("words", false, "count words")
		showBytes = flag.Bool("bytes", false, "count bytes")
	)
	flag.Parse()

	// If no flag is set, mimic real wc's default: show all three.
	if !*showLines && !*showWords && !*showBytes {
		*showLines, *showWords, *showBytes = true, true, true
	}

	paths := flag.Args() // everything left over after flag parsing
	if len(paths) == 0 {
		fmt.Fprintln(os.Stderr, "usage: wcish [-lines] [-words] [-bytes] file...")
		os.Exit(2)
	}

	for _, path := range paths {
		if err := countOne(path, *showLines, *showWords, *showBytes); err != nil {
			fmt.Fprintf(os.Stderr, "wcish: %v\n", err)
			os.Exit(1)
		}
	}
}

func countOne(path string, showLines, showWords, showBytes bool) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Read the whole file into memory so we can count all three metrics
	// without opening the file three times. Fine for the sizes wc-shaped
	// tools handle in practice.
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var out string
	if showLines {
		lines, _ := fileio.CountLines(byteReader(data))
		out += fmt.Sprintf(" %6d", lines)
	}
	if showWords {
		words, _ := fileio.CountWords(byteReader(data))
		out += fmt.Sprintf(" %6d", words)
	}
	if showBytes {
		out += fmt.Sprintf(" %6d", len(data))
	}
	fmt.Printf("%s %s\n", out, path)
	return nil
}

// byteReader wraps a []byte in an io.Reader without importing bytes.
// (Kept small — real code would just use bytes.NewReader.)
type byteReader []byte

func (b byteReader) Read(p []byte) (int, error) {
	if len(b) == 0 {
		return 0, io.EOF
	}
	n := copy(p, b)
	return n, io.EOF
}
