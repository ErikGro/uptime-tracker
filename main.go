package main

import (
	"log"
	"net/http"

	"github.com/ErikGro/uptime-tracker/internal/web"
)

func main() {
	addr := ":8080"
	log.Printf("uptime-tracker listening on %s", addr)
	if err := http.ListenAndServe(addr, web.NewServer()); err != nil {
		log.Fatal(err)
	}
}
