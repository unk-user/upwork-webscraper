package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	URL         = "https://www.upwork.com/nx/search/jobs/"
	QueryPrompt = "please provide a query: \n"
)

func ScanQuery(out io.Writer, in io.Reader) (string, error) {
	fmt.Fprint(out, QueryPrompt)
	reader := bufio.NewReader(in)

	query, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	query = strings.TrimSpace(query)
	fmt.Printf("query: %q\n", query)
	return query, nil
}

func main() {
	ScanQuery(os.Stdout, os.Stdin)
}
