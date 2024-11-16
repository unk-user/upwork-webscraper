package main

import (
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

const (
	URL         = "https://www.upwork.com/nx/search/jobs/?category2_uid=531770282580668418&sort=recency&per_page=20"
	QueryPrompt = "please provide a list of keywords: \n"
)

func MakeParams(keywords string) string {
	array := strings.Fields(keywords)
	return "&q=%28" + strings.Join(array, "%20OR%20") + "%29"
}

func GetNewJobs(keywords string, launcher *launcher.Launcher) (JobMap map[string]string, err error) {
	const (
		navigationTimeout  = 5 * time.Second
		requestIdleTimeout = 10 * time.Second
		htmlTimeout        = 5 * time.Second
	)

	err = rod.Try(func() {
		defer launcher.Cleanup()

		defer launcher.Kill()

		u := launcher.MustLaunch()

		browser := rod.New().ControlURL(u).MustConnect()
		defer browser.MustClose()

		page := browser.MustPage()

		router := page.HijackRequests()

		router.MustAdd("*", func(ctx *rod.Hijack) {
			if ctx.Request.Type() != proto.NetworkResourceTypeDocument {
				ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
				return
			}

			ctx.ContinueRequest(&proto.FetchContinueRequest{})
		})

		go router.Run()

		if err != nil {
			panic(err)
		}
		fullUrl := URL + MakeParams(keywords)
		page.Timeout(navigationTimeout).MustNavigate(fullUrl).MustWaitElementsMoreThan(".job-tile", 9)

		waitRequestIdle := page.Timeout(requestIdleTimeout).MustWaitRequestIdle()
		waitRequestIdle()

		jobTiles := page.Timeout(htmlTimeout).MustElements(".job-tile")
		for _, jobTile := range jobTiles {
			title := jobTile.MustElement(".job-tile-title").MustText()
			uid := *jobTile.MustAttribute("data-ev-job-uid")

			JobMap[uid] = title
		}
	})

	return JobMap, err
}
