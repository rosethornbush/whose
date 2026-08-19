package registry

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const cacheTTL = 24 * time.Hour

func Fetch(ctx context.Context, name, url string) ([]byte, error) {
	cachePath, err := registryCachePath(name)
	if err != nil {
		return nil, err
	}

	if data, ok := readFreshCache(cachePath); ok {
		return data, nil
	}

	stale, _ := os.ReadFile(cachePath)

	data, err := fetch(ctx, url)
	if err != nil {
		if stale != nil {
			return stale, nil
		}

		return nil, err
	}

	if err := writeCache(cachePath, data); err != nil {
		return data, nil
	}

	return data, nil
}

func registryCachePath(name string) (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("find cache directory: %w", err)
	}

	return filepath.Join(dir, "whose", "rdap", name+".json"), nil
}

func readFreshCache(path string) ([]byte, bool) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, false
	}

	if time.Since(info.ModTime()) > cacheTTL {
		return nil, false
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	return data, true
}

func writeCache(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create cache directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write registry cache: %w", err)
	}

	return nil
}

func fetch(ctx context.Context, url string) ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create IANA registry request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request IANA registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("IANA registry returned %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read IANA registry: %w", err)
	}

	return data, nil
}
