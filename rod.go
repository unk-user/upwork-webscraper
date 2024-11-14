package main

import (
	"fmt"
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

// func ScanQuery(out io.Writer, in io.Reader) (string, error) {
// 	fmt.Fprint(out, QueryPrompt)
// 	reader := bufio.NewReader(in)

// 	query, err := reader.ReadString('\n')
// 	if err != nil {
// 		return "", err
// 	}

// 	query = strings.TrimSpace(query)

// 	return query, nil
// }

func MakeParams(keywords string) string {
	array := strings.Fields(keywords)
	return "&q=%28" + strings.Join(array, "%20OR%20") + "%29"
}

func GetNewJobs(keyword string) (err error) {
	const (
		navigationTimeout  = 5 * time.Second
		requestIdleTimeout = 10 * time.Second
		htmlTimeout        = 5 * time.Second
	)

	err = rod.Try(func() {
		launcher := launchInLambda()

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
		fullUrl := URL + MakeParams(keyword)
		page.Timeout(navigationTimeout).MustNavigate(fullUrl).MustWaitElementsMoreThan(".job-tile", 9)

		waitRequestIdle := page.Timeout(requestIdleTimeout).MustWaitRequestIdle()
		waitRequestIdle()

		jobTiles := page.Timeout(htmlTimeout).MustElements(".job-tile")
		for i, jobTile := range jobTiles {
			title := jobTile.MustElement(".job-tile-title").MustText()
			uid := *jobTile.MustAttribute("data-ev-job-uid")

			fmt.Printf("%d - %s: %s\n", i, uid, title)
		}
	})

	if err != nil {
		return err
	}

	return nil
}

func launchInLambda() *launcher.Launcher {
	return launcher.New().
		Bin("/opt/chromium").
		Set("allow-running-insecure-content").
		Set("autoplay-policy", "user-gesture-required").
		Set("disable-component-update").
		Set("disable-domain-reliability").
		Set("disable-features", "AudioServiceOutOfProcess", "IsolateOrigins", "site-per-process").
		Set("disable-print-preview").
		Set("disable-setuid-sandbox").
		Set("disable-site-isolation-trials").
		Set("disable-speech-api").
		Set("disable-web-security").
		Set("disk-cache-size", "33554432").
		Set("enable-features", "SharedArrayBuffer").
		Set("hide-scrollbars").
		Set("ignore-gpu-blocklist").
		Set("in-process-gpu").
		Set("mute-audio").
		Set("no-default-browser-check").
		Set("no-pings").
		Set("no-sandbox").
		Set("no-zygote").
		Set("single-process").
		Set("use-gl", "swiftshader").
		Set("window-size", "1920", "1080")
}
