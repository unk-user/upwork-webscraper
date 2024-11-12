package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-rod/rod"
)

const (
	URL         = "https://www.upwork.com/nx/search/jobs/?category2_uid=531770282580668418&sort=recency"
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
	return "&q=%28" + strings.Join(array, "%20OR%20") + "%29"
}

func GetNewJobs() {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage("")

	keywords, err := ScanQuery(os.Stdout, os.Stdin)
	if err != nil {
		panic(err)
	}
	fullUrl := URL + MakeParams(keywords)
	page.MustNavigate(fullUrl).MustWaitElementsMoreThan(".job-tile", 9)

	jobTiles := page.MustElements(".job-tile")
	for i, jobTile := range jobTiles {
		title := jobTile.MustElement(".job-tile-title").MustText()
		uid := *jobTile.MustAttribute("data-ev-job-uid")

		fmt.Printf("%d - %s: %s\n", i, uid, title)
	}
}
