package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const cacheTTL = 24 * time.Hour

func Fetch(ctx context.Context, name, url string) ([]byte, error) {
	cacheRoot, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("find cache directory: %w", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	return fetchCached(
		ctx,
		filepath.Join(cacheRoot, "whose", "rdap"),
		client,
		name,
		url,
	)
}

func fetchCached(
	ctx context.Context,
	cacheDir string,
	client *http.Client,
	name string,
	url string,
) ([]byte, error) {
	cachePath := filepath.Join(cacheDir, name+".json")

	if data, ok := readFreshCache(cachePath); ok {
		return data, nil
	}

	stale, _ := os.ReadFile(cachePath)
	if !json.Valid(stale) {
		stale = nil
	}

	data, err := fetch(ctx, client, url)
	if err != nil {
		if stale != nil {
			return stale, nil
		}

		return nil, err
	}

	if !json.Valid(data) {
		if stale != nil {
			return stale, nil
		}

		return nil, fmt.Errorf("IANA registry returned invalid JSON")
	}

	if err := writeCache(cachePath, data); err != nil {
		return data, nil
	}

	return data, nil
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

	if !json.Valid(data) {
		return nil, false
	}

	return data, true
}

func writeCache(path string, data []byte) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create cache directory: %w", err)
	}

	file, err := os.CreateTemp(dir, ".whose-*.tmp")
	if err != nil {
		return fmt.Errorf("create temporary cache file: %w", err)
	}

	tempPath := file.Name()

	defer func() {
		_ = file.Close()
		_ = os.Remove(tempPath)
	}()

	if err := file.Chmod(0o644); err != nil {
		return fmt.Errorf("set cache permissions: %w", err)
	}

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("write temporary cache file: %w", err)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("close temporary cache file: %w", err)
	}

	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("replace registry cache: %w", err)
	}

	return nil
}

func fetch(
	ctx context.Context,
	client *http.Client,
	url string,
) ([]byte, error) {
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
