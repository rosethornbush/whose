package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rosethornbush/whose/output"
	"github.com/rosethornbush/whose/rdap"
	"github.com/rosethornbush/whose/registry"
)

const dnsRegistryURL = "https://data.iana.org/rdap/dns.json"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: whose <query>")
		os.Exit(2)
	}

	domain := os.Args[1]

	ctx := context.Background()

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	// TODO: Cache the IANA registry using its HTTP expiration metadata.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, dnsRegistryURL, nil)
	if err != nil {
		fail(err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		fail(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fail(fmt.Errorf("IANA registry returned %s", resp.Status))
	}

	server, err := registry.LookupDomain(resp.Body, domain)
	if err != nil {
		fail(err)
	}

	client := rdap.NewClient()

	result, err := client.LookupDomain(ctx, server, domain)
	if err != nil {
		fail(err)
	}

	if err := output.Domain(os.Stdout, result.Domain); err != nil {
		fail(err)
	}

	// var response any
	// if err := json.Unmarshal(result.Raw, &response); err != nil {
	// 	fail(fmt.Errorf("decode RDAP response: %w", err))
	// }

	// encoder := json.NewEncoder(os.Stdout)
	// encoder.SetIndent("", "  ")
	// encoder.SetEscapeHTML(false)

	// if err := encoder.Encode(response); err != nil {
	// 	fail(fmt.Errorf("encode RDAP response: %w", err))
	// }
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "whose:", err)
	os.Exit(1)
}
