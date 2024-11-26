package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

func launchInLambda() *launcher.Launcher {
	return launcher.New().Bin("/opt/chromium").
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

func GetHTML(url string) (html string, err error) {
	const (
		navigationTimeout  = 5 * time.Second
		requestIdleTimeout = 10 * time.Second
		htmlTimeout        = 5 * time.Second
	)
	launcher := launchInLambda()

	defer launcher.Cleanup()

	defer launcher.Kill()

	u, err := launcher.Launch()
	if err != nil {
		return "", err
	}

	browser := rod.New().ControlURL(u).MustConnect()
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

	page.Timeout(navigationTimeout).MustNavigate(url)

	waitNavigation := page.Timeout(navigationTimeout).MustWaitNavigation()
	waitNavigation()

	err = page.WaitElementsMoreThan("section.card-list-container", 0)
	if err != nil {
		return "", fmt.Errorf("Job listing card not found: %w", err)
	}

	html = page.Timeout(htmlTimeout).MustElement("section.card-list-container").MustHTML()

	return html, err
}

type Job struct {
	UID             string   `json:"uid"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	JobType         string   `json:"jobType"`
	ExperienceLevel string   `json:"experienceLevel"`
	PublishedAt     string   `json:"publishedAt"`
	FixedPrice      string   `json:"fixedPrice"`
	Duration        string   `json:"duration"`
	Skills          []string `json:"skills"`
}

func ProcessHTML(html string) (jobs []Job, err error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	nodes := doc.Find("article.job-tile")
	if nodes.Length() == 0 {
		return nil, fmt.Errorf("no jobs found")
	}

	nodes.Each(func(i int, s *goquery.Selection) {
		uid, exists := s.Attr("data-ev-job-uid")
		if !exists {
			return
		}

		publishedAt := s.Find("[data-test=job-pubilshed-date]").Text()
		entries := strings.Fields(publishedAt)
		num, time := entries[1], entries[2]
		intNum, _ := strconv.Atoi(num)
		if time != "minutes" || intNum > 30 {
			return
		}

		title := s.Find("h2.job-tile-title").Text()
		description := s.Find("p.mb-0.text-body-sm").Text()

		detailsInfo := s.Find("ul.job-tile-info-list")
		jobType := detailsInfo.Find("[data-test='job-type-label']").Text()
		experienceLevel := detailsInfo.Find("[data-test='experience-level']").Text()
		duration := detailsInfo.Find("[data-test='duration-label']").Text()
		fixedPrice := detailsInfo.Find("[data-test='is-fixed-price']").Text()
		skills := make([]string, 0)

		skillsContainer := s.Find("div.air3-token-container")
		skillsContainer.Find("span").Each(func(i int, s *goquery.Selection) {
			skill := s.Text()
			skills = append(skills, skill)
		})

		jobs = append(jobs, Job{
			UID:             uid,
			Title:           title,
			Description:     description,
			JobType:         jobType,
			ExperienceLevel: experienceLevel,
			PublishedAt:     publishedAt,
			Duration:        duration,
			FixedPrice:      fixedPrice,
			Skills:          skills,
		})
	})

	return jobs, nil
}
