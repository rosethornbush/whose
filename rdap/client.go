package rdap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	http *http.Client
}

func NewClient() *Client {
	return &Client{
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) LookupDomain(ctx context.Context, baseURL, domain string) (*DomainResult, error) {
	endpoint := strings.TrimRight(baseURL, "/") +
		"/domain/" +
		url.PathEscape(domain)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create RDAP request: %w", err)
	}

	req.Header.Set("Accept", "application/rdap+json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request RDAP data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("RDAP server returned %s", resp.Status)
	}

	// TODO: Limit the response size before reading it into memory.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read RDAP response: %w", err)
	}

	var domainResponse Domain

	if err := json.Unmarshal(body, &domainResponse); err != nil {
		return nil, fmt.Errorf("decode RDAP response: %w", err)
	}

	return &DomainResult{
		Domain: domainResponse,
		Raw:    body,
	}, nil
}

func (c *Client) LookupIP(
	ctx context.Context,
	baseURL string,
	address string,
) (*IPResult, error) {
	endpoint := strings.TrimRight(baseURL, "/") +
		"/ip/" +
		url.PathEscape(address)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create RDAP request: %w", err)
	}

	req.Header.Set("Accept", "application/rdap+json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request RDAP data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("RDAP server returned %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read RDAP response: %w", err)
	}

	var network IPNetwork

	if err := json.Unmarshal(body, &network); err != nil {
		return nil, fmt.Errorf("decode RDAP response: %w", err)
	}

	return &IPResult{
		Network: network,
		Raw:     body,
	}, nil
}
