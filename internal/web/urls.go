package web

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/ErikGro/uptime-tracker/internal/web/templates"
	"gorm.io/gorm"
)

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	urls, err := s.store.ListURLs()
	if err != nil {
		http.Error(w, "list urls: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = templates.Dashboard(urls).Render(r.Context(), w)
}

func (s *Server) handleURLCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	label, target, errs := validateURLInput(r.FormValue("label"), r.FormValue("url"))
	if len(errs) == 0 {
		if _, err := s.store.CreateURL(label, target); err != nil {
			if isUniqueViolation(err) {
				errs["url"] = "This URL is already being monitored"
			} else {
				http.Error(w, "create url: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	if len(errs) > 0 {
		// Retarget the swap to the form itself so error markup replaces the form,
		// not the URL list.
		w.Header().Set("HX-Retarget", "#new-url-form")
		w.Header().Set("HX-Reswap", "outerHTML")
		_ = templates.URLNewForm(label, target, errs).Render(r.Context(), w)
		return
	}

	urls, err := s.store.ListURLs()
	if err != nil {
		http.Error(w, "list urls: "+err.Error(), http.StatusInternalServerError)
		return
	}
	_ = templates.URLList(urls).Render(r.Context(), w)
	_ = templates.URLNewFormOOB().Render(r.Context(), w)
}

func (s *Server) handleURLEditForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	u, err := s.store.GetURL(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "get url: "+err.Error(), http.StatusInternalServerError)
		return
	}
	_ = templates.URLEditRow(*u, u.Label, u.URL, nil).Render(r.Context(), w)
}

func (s *Server) handleURLRow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	u, err := s.store.GetURL(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "get url: "+err.Error(), http.StatusInternalServerError)
		return
	}
	_ = templates.URLRow(*u).Render(r.Context(), w)
}

func (s *Server) handleURLUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	existing, err := s.store.GetURL(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "get url: "+err.Error(), http.StatusInternalServerError)
		return
	}

	label, target, errs := validateURLInput(r.FormValue("label"), r.FormValue("url"))
	if len(errs) == 0 {
		updated, err := s.store.UpdateURL(id, label, target)
		if err != nil {
			if isUniqueViolation(err) {
				errs["url"] = "This URL is already being monitored"
			} else {
				http.Error(w, "update url: "+err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			_ = templates.URLRow(*updated).Render(r.Context(), w)
			return
		}
	}
	_ = templates.URLEditRow(*existing, label, target, errs).Render(r.Context(), w)
}

func (s *Server) handleURLDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	if err := s.store.DeleteURL(id); err != nil {
		http.Error(w, "delete url: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("HX-Trigger", "url-list-changed")
	// Empty body — outerHTML swap on the row removes it from the DOM.
}

func parseID(w http.ResponseWriter, r *http.Request) (uint, bool) {
	raw := r.PathValue("id")
	n, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || n == 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return 0, false
	}
	return uint(n), true
}

func validateURLInput(rawLabel, rawTarget string) (string, string, map[string]string) {
	errs := map[string]string{}
	label := strings.TrimSpace(rawLabel)
	target := strings.TrimSpace(rawTarget)

	switch {
	case label == "":
		errs["label"] = "Label is required"
	case len(label) > 100:
		errs["label"] = "Label must be 100 characters or fewer"
	}

	switch {
	case target == "":
		errs["url"] = "URL is required"
	default:
		u, err := url.Parse(target)
		if err != nil || u.Host == "" {
			errs["url"] = "URL must be a valid absolute URL"
		} else if u.Scheme != "http" && u.Scheme != "https" {
			errs["url"] = "URL must use http or https"
		}
	}

	return label, target, errs
}

func isUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "UNIQUE constraint failed")
}
