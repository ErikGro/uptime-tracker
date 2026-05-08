package web

import (
	"net/http"

	"github.com/ErikGro/uptime-tracker/internal/web/templates"
)

func NewServer() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(StaticFS())))

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/poc", http.StatusFound)
	})

	mux.HandleFunc("GET /poc", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = templates.POC().Render(r.Context(), w)
	})
	mux.HandleFunc("POST /poc/ping", handlePOCPing)
	mux.HandleFunc("GET /poc/time", handlePOCTime)

	return mux
}
