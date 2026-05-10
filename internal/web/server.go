package web

import (
	"net/http"

	"github.com/ErikGro/uptime-tracker/internal/config"
	"github.com/ErikGro/uptime-tracker/internal/store"
	"github.com/ErikGro/uptime-tracker/internal/web/templates"
)

type Server struct {
	cfg   *config.Config
	store *store.Store
}

func NewServer(cfg *config.Config, st *store.Store) http.Handler {
	s := &Server{cfg: cfg, store: st}
	auth := basicAuth(cfg.AdminUser, cfg.AdminPass)

	mux := http.NewServeMux()

	// Public routes — no auth.
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(StaticFS())))
	mux.HandleFunc("GET /healthz", s.handleHealthz)

	// Authenticated routes.
	mux.HandleFunc("GET /{$}", auth(s.handleDashboard))
	mux.HandleFunc("POST /urls", auth(s.handleURLCreate))
	mux.HandleFunc("GET /urls/{id}/edit", auth(s.handleURLEditForm))
	mux.HandleFunc("GET /urls/{id}/row", auth(s.handleURLRow))
	mux.HandleFunc("PUT /urls/{id}", auth(s.handleURLUpdate))
	mux.HandleFunc("DELETE /urls/{id}", auth(s.handleURLDelete))
	mux.HandleFunc("GET /poc", auth(s.handlePOC))
	mux.HandleFunc("POST /poc/ping", auth(s.handlePOCPing))
	mux.HandleFunc("GET /poc/time", auth(s.handlePOCTime))

	return mux
}

func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func (s *Server) handlePOC(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = templates.POC().Render(r.Context(), w)
}
