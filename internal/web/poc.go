package web

import (
	"context"
	"net/http"
	"time"

	"github.com/ErikGro/uptime-tracker/internal/web/templates"
)

func handlePOCPing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	url := r.FormValue("url")
	if url == "" {
		_ = templates.PingResultFragment(templates.PingResult{
			URL:   "(empty)",
			OK:    false,
			Error: "no URL provided",
		}).Render(r.Context(), w)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		_ = templates.PingResultFragment(templates.PingResult{URL: url, Error: err.Error()}).Render(r.Context(), w)
		return
	}

	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	latency := time.Since(start)
	if err != nil {
		_ = templates.PingResultFragment(templates.PingResult{URL: url, Error: err.Error()}).Render(r.Context(), w)
		return
	}
	defer resp.Body.Close()

	_ = templates.PingResultFragment(templates.PingResult{
		URL:     url,
		OK:      resp.StatusCode < 400,
		Status:  resp.StatusCode,
		Latency: latency.Round(time.Millisecond).String(),
	}).Render(r.Context(), w)
}

func handlePOCTime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = templates.ClockFragment(time.Now().Format("15:04:05")).Render(r.Context(), w)
}
