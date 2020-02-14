package wc

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

type Counter struct {
	lines  int64
	words  int64
	bytes  int64
	file   *os.File
	reader ReaderConfig
}

type ReaderConfig struct {
	Reader io.Reader
	Name   string
}

type CounterList []*Counter

var (
	pattern = "\t%d\t%d\t%d\t%s\t\n"
)

func NewCounter(reader ReaderConfig) Counter {
	return Counter{
		reader: reader,
	}
}

func (cl CounterList) String() string {
	var b strings.Builder
	w := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)
	sum := NewCounter(ReaderConfig{
		Name: "total",
	})
	for _, c := range cl {
		fmt.Fprintf(w, pattern, c.lines, c.words, c.bytes, c.reader.Name)

		sum.lines += c.lines
		sum.words += c.words
		sum.bytes += c.bytes
	}

	if len(cl) > 1 {
		fmt.Fprintf(w, pattern, sum.lines, sum.words, sum.bytes, sum.reader.Name)
	}

	err := w.Flush()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return b.String()
}

func (c Counter) String() string {
	b := strings.Builder{}
	w := tabwriter.NewWriter(&b, 0, 0, 4, ' ', 0)
	fmt.Fprintf(w, pattern, c.lines, c.words, c.bytes, c.reader.Name)
	w.Flush()
	return b.String()
}

func (c *Counter) ExecuteCount() {
	scanner := bufio.NewScanner(c.reader.Reader)
	for scanner.Scan() {
		text := scanner.Text()

		words := strings.Fields(text)
		bytes := []byte(text)

		c.bytes += int64(len(bytes))
		c.words += int64(len(words))
		c.lines++
	}
}

func CountWC(filename string) (*Counter, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	c := NewCounter(ReaderConfig{
		Reader: file,
		Name:   file.Name(),
	})
	c.ExecuteCount()
	return &c, nil
}
