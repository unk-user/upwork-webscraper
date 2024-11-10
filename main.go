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
	QueryPrompt = "please provide a list of keywords: \n"
)

func ScanQuery(out io.Writer, in io.Reader) (string, error) {
	fmt.Fprint(out, QueryPrompt)
	reader := bufio.NewReader(in)

	query, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	query = strings.TrimSpace(query)

	return query, nil
}

func MakeParams(keywords string) string {
	array := strings.Fields(keywords)
	return "?q=%28" + strings.Join(array, "%20OR%20") + "%29"
}

func main() {
	keywords, err := ScanQuery(os.Stdout, os.Stdin)
	if err != nil {
		panic(err)
	}
	fmt.Printf("query: %q\n", keywords)
}
