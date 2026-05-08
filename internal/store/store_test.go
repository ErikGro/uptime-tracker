package store

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	st, err := Open(dbPath)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return st
}

func TestStoreSeedsSettings(t *testing.T) {
	st := newTestStore(t)

	if got := st.PollInterval(); got != 300*time.Second {
		t.Errorf("PollInterval = %v, want 300s", got)
	}
	if got := st.FailureThreshold(); got != 3 {
		t.Errorf("FailureThreshold = %d, want 3", got)
	}
	if got := st.WebhookEnabled(); got != false {
		t.Errorf("WebhookEnabled = %v, want false", got)
	}
}

func TestSetSettingPersists(t *testing.T) {
	st := newTestStore(t)

	if err := st.SetSetting(KeyPollInterval, "60"); err != nil {
		t.Fatalf("SetSetting: %v", err)
	}
	if got := st.PollInterval(); got != 60*time.Second {
		t.Errorf("PollInterval after set = %v, want 60s", got)
	}
}

func TestURLCRUD(t *testing.T) {
	st := newTestStore(t)

	created, err := st.CreateURL("Example", "https://example.com")
	if err != nil {
		t.Fatalf("CreateURL: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("CreateURL returned zero ID")
	}
	if created.CurrentStatus != StatusUnknown {
		t.Errorf("status = %q, want %q", created.CurrentStatus, StatusUnknown)
	}

	got, err := st.GetURL(created.ID)
	if err != nil {
		t.Fatalf("GetURL: %v", err)
	}
	if got.URL != "https://example.com" {
		t.Errorf("got URL %q", got.URL)
	}

	if _, err := st.UpdateURL(created.ID, "Example v2", "https://example.org"); err != nil {
		t.Fatalf("UpdateURL: %v", err)
	}
	got, _ = st.GetURL(created.ID)
	if got.Label != "Example v2" || got.URL != "https://example.org" {
		t.Errorf("after update: %+v", got)
	}

	all, err := st.ListURLs()
	if err != nil {
		t.Fatalf("ListURLs: %v", err)
	}
	if len(all) != 1 {
		t.Errorf("ListURLs len = %d, want 1", len(all))
	}

	if err := st.DeleteURL(created.ID); err != nil {
		t.Fatalf("DeleteURL: %v", err)
	}
	if _, err := st.GetURL(created.ID); err == nil {
		t.Error("GetURL after delete should error")
	}
}

func TestChecksAppendAndPrune(t *testing.T) {
	st := newTestStore(t)
	url, _ := st.CreateURL("x", "https://x.example")

	old := time.Now().Add(-48 * time.Hour)
	fresh := time.Now()

	if err := st.AppendCheck(&Check{URLID: url.ID, CheckedAt: old, OK: true, StatusCode: 200}); err != nil {
		t.Fatal(err)
	}
	if err := st.AppendCheck(&Check{URLID: url.ID, CheckedAt: fresh, OK: false, StatusCode: 500}); err != nil {
		t.Fatal(err)
	}

	checks, err := st.ListChecksFor(url.ID, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(checks) != 2 {
		t.Errorf("len = %d, want 2", len(checks))
	}

	cutoff := time.Now().Add(-24 * time.Hour)
	pruned, err := st.PruneChecksOlderThan(cutoff)
	if err != nil {
		t.Fatal(err)
	}
	if pruned != 1 {
		t.Errorf("pruned = %d, want 1", pruned)
	}
}

func TestUpdateURLStatus(t *testing.T) {
	st := newTestStore(t)
	url, _ := st.CreateURL("x", "https://x.example")

	now := sql.NullTime{Time: time.Now(), Valid: true}
	if err := st.UpdateURLStatus(url.ID, StatusDown, 3, now); err != nil {
		t.Fatal(err)
	}

	got, _ := st.GetURL(url.ID)
	if got.CurrentStatus != StatusDown || got.ConsecutiveFailures != 3 {
		t.Errorf("got %+v", got)
	}
}
