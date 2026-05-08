package web

import (
	"crypto/subtle"
	"net/http"
)

func basicAuth(user, pass string) func(http.HandlerFunc) http.HandlerFunc {
	userBytes := []byte(user)
	passBytes := []byte(pass)
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			u, p, ok := r.BasicAuth()
			userOK := ok && subtle.ConstantTimeCompare([]byte(u), userBytes) == 1
			passOK := ok && subtle.ConstantTimeCompare([]byte(p), passBytes) == 1
			if !userOK || !passOK {
				w.Header().Set("WWW-Authenticate", `Basic realm="uptime-tracker", charset="UTF-8"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next(w, r)
		}
	}
}
