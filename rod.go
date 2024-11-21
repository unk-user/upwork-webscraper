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
	URL = "https://www.upwork.com/nx/search/jobs/"
)

type Payload struct {
	CategoryId string
	Keywords   string
}

func MakeParams(p Payload) string {
	keywordsArr := strings.Fields(p.Keywords)
	return "?category2_uid=" + p.CategoryId + "&per_page=20" + "&q=%28" + strings.Join(keywordsArr, "%20OR%20") + "%29"
}

type Job struct {
	UID             string `json:"uid"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	JobType         string `json:"jobType"`
	ExperienceLevel string `json:"experienceLevel"`
}

func GetNewJobs(p Payload) (jobs []Job, err error) {
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

		fullUrl := URL + MakeParams(p)
		page.Timeout(navigationTimeout).MustNavigate(fullUrl).MustWaitElementsMoreThan(".job-tile", 9)

		jobTiles := page.MustElements(".job-tile")

		if len(jobTiles) == 0 {
			panic("no jobs found")
		}

		for _, jobTile := range jobTiles {
			uid := *jobTile.MustAttribute("data-ev-job-uid")
			title := jobTile.MustElement(".job-tile-title").MustText()
			description := jobTile.MustElement(".text-body-sm").MustText()
			jobType := jobTile.MustElement("[data-test=job-type-label]").MustText()
			experienceLevel := jobTile.MustElement("[data-test=experience-level]").MustText()

			job := Job{
				UID:             uid,
				Title:           title,
				Description:     description,
				JobType:         jobType,
				ExperienceLevel: experienceLevel,
			}

			jobs = append(jobs, job)
		}
	})

	if err != nil {
		return nil, fmt.Errorf("failed to launch: %w", err)
	}

	return jobs, nil
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
