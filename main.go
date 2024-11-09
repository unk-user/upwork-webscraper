package main

import (
	"fmt"

	"github.com/go-rod/rod"
)

func main() {
	browser := rod.New().MustConnect()

	defer browser.MustClose()

	page := browser.MustPage("https://www.upwork.com/nx/search/jobs/?nbs=1&q=html").MustWaitStable()

	list := page.MustElements(".job-tile")

	for _, el := range list {
		fmt.Printf("Title: %q, uid: %q\n", el.MustElement("h2").MustText(), *el.MustAttribute("data-ev-job-uid"))
	}
}
