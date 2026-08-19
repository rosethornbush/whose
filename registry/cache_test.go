package registry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFetchCachedUsesFreshCache(t *testing.T) {
	cacheDir := t.TempDir()
	cachePath := filepath.Join(cacheDir, "dns.json")

	cached := []byte(`{"services":[]}`)
	remote := []byte(`{"services":[[["com"],["https://rdap.example/"]]]}`)

	if err := os.WriteFile(cachePath, cached, 0o644); err != nil {
		t.Fatal(err)
	}

	before, err := os.Stat(cachePath)
	if err != nil {
		t.Fatal(err)
	}

	requests := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		_, _ = w.Write(remote)
	}))
	defer server.Close()

	got, err := fetchCached(
		context.Background(),
		cacheDir,
		server.Client(),
		"dns",
		server.URL,
	)
	if err != nil {
		t.Fatal(err)
	}

	after, err := os.Stat(cachePath)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != string(cached) {
		t.Fatalf("got %q, want %q", got, cached)
	}

	if requests != 0 {
		t.Fatalf("got %d HTTP requests, want 0", requests)
	}

	if !before.ModTime().Equal(after.ModTime()) {
		t.Fatalf(
			"cache mtime changed from %v to %v",
			before.ModTime(),
			after.ModTime(),
		)
	}
}

func TestFetchCachedRefreshesStaleCache(t *testing.T) {
	cacheDir := t.TempDir()
	cachePath := filepath.Join(cacheDir, "dns.json")

	stale := []byte(`{"services":[]}`)
	fresh := []byte(`{"services":[[["com"],["https://rdap.example/"]]]}`)

	if err := os.WriteFile(cachePath, stale, 0o644); err != nil {
		t.Fatal(err)
	}

	staleTime := time.Now().Add(-cacheTTL - time.Hour)

	if err := os.Chtimes(cachePath, staleTime, staleTime); err != nil {
		t.Fatal(err)
	}

	requests := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		_, _ = w.Write(fresh)
	}))
	defer server.Close()

	got, err := fetchCached(
		context.Background(),
		cacheDir,
		server.Client(),
		"dns",
		server.URL,
	)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != string(fresh) {
		t.Fatalf("got %q, want %q", got, fresh)
	}

	if requests != 1 {
		t.Fatalf("got %d HTTP requests, want 1", requests)
	}

	cached, err := os.ReadFile(cachePath)
	if err != nil {
		t.Fatal(err)
	}

	if string(cached) != string(fresh) {
		t.Fatalf("cached %q, want %q", cached, fresh)
	}
}

func TestFetchCachedFallsBackToStaleCache(t *testing.T) {
	cacheDir := t.TempDir()
	cachePath := filepath.Join(cacheDir, "dns.json")

	stale := []byte(`{"services":[]}`)

	if err := os.WriteFile(cachePath, stale, 0o644); err != nil {
		t.Fatal(err)
	}

	staleTime := time.Now().Add(-cacheTTL - time.Hour)

	if err := os.Chtimes(cachePath, staleTime, staleTime); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "broken", http.StatusInternalServerError)
	}))
	defer server.Close()

	got, err := fetchCached(
		context.Background(),
		cacheDir,
		server.Client(),
		"dns",
		server.URL,
	)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != string(stale) {
		t.Fatalf("got %q, want %q", got, stale)
	}
}

func TestFetchCachedFallsBackToStaleCacheOnInvalidJSON(t *testing.T) {
	cacheDir := t.TempDir()
	cachePath := filepath.Join(cacheDir, "dns.json")

	stale := []byte(`{"services":[]}`)

	if err := os.WriteFile(cachePath, stale, 0o644); err != nil {
		t.Fatal(err)
	}

	staleTime := time.Now().Add(-cacheTTL - time.Hour)

	if err := os.Chtimes(cachePath, staleTime, staleTime); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("not-json"))
	}))
	defer server.Close()

	got, err := fetchCached(
		context.Background(),
		cacheDir,
		server.Client(),
		"dns",
		server.URL,
	)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != string(stale) {
		t.Fatalf("got %q, want %q", got, stale)
	}
}

func TestFetchCachedRefetchesCorruptFreshCache(t *testing.T) {
	cacheDir := t.TempDir()
	cachePath := filepath.Join(cacheDir, "dns.json")

	if err := os.WriteFile(cachePath, []byte("not-json"), 0o644); err != nil {
		t.Fatal(err)
	}

	fresh := []byte(`{"services":[]}`)
	requests := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		_, _ = w.Write(fresh)
	}))
	defer server.Close()

	got, err := fetchCached(
		context.Background(),
		cacheDir,
		server.Client(),
		"dns",
		server.URL,
	)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != string(fresh) {
		t.Fatalf("got %q, want %q", got, fresh)
	}

	if requests != 1 {
		t.Fatalf("got %d HTTP requests, want 1", requests)
	}

	cached, err := os.ReadFile(cachePath)
	if err != nil {
		t.Fatal(err)
	}

	if string(cached) != string(fresh) {
		t.Fatalf("cached %q, want %q", cached, fresh)
	}
}

func TestFetchCachedFailsWithoutCache(t *testing.T) {
	cacheDir := t.TempDir()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "broken", http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := fetchCached(
		context.Background(),
		cacheDir,
		server.Client(),
		"dns",
		server.URL,
	)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFetchCachedFailsWithoutCacheOnInvalidJSON(t *testing.T) {
	cacheDir := t.TempDir()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("not-json"))
	}))
	defer server.Close()

	_, err := fetchCached(
		context.Background(),
		cacheDir,
		server.Client(),
		"dns",
		server.URL,
	)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestFetchCachedDoesNotFallBackToCorruptCache(t *testing.T) {
	cacheDir := t.TempDir()
	cachePath := filepath.Join(cacheDir, "dns.json")

	if err := os.WriteFile(cachePath, []byte("not-json"), 0o644); err != nil {
		t.Fatal(err)
	}

	staleTime := time.Now().Add(-cacheTTL - time.Hour)

	if err := os.Chtimes(cachePath, staleTime, staleTime); err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "broken", http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := fetchCached(
		context.Background(),
		cacheDir,
		server.Client(),
		"dns",
		server.URL,
	)
	if err == nil {
		t.Fatal("expected error")
	}
}
