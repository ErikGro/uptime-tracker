package main

import (
	"log"
	"net/http"

	"github.com/ErikGro/uptime-tracker/internal/config"
	"github.com/ErikGro/uptime-tracker/internal/store"
	"github.com/ErikGro/uptime-tracker/internal/web"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	st, err := store.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("store: %v", err)
	}
	defer st.Close()

	log.Printf("uptime-tracker listening on %s (db=%s)", cfg.ListenAddr, cfg.DBPath)
	if err := http.ListenAndServe(cfg.ListenAddr, web.NewServer(cfg, st)); err != nil {
		log.Fatal(err)
	}
}
