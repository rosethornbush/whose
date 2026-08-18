package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rosethornbush/whose/output"
	"github.com/rosethornbush/whose/rdap"
	"github.com/rosethornbush/whose/registry"
)

const dnsRegistryURL = "https://data.iana.org/rdap/dns.json"

func main() {
	var (
		query      string
		jsonOutput bool
		rawOutput  bool
	)

	for _, arg := range os.Args[1:] {
		switch arg {
		case "--json":
			jsonOutput = true
		case "--raw":
			rawOutput = true
		default:
			if strings.HasPrefix(arg, "-") {
				usage(fmt.Sprintf("unknown option %q", arg))
			}

			if query != "" {
				usage("expected exactly one query")
			}

			query = arg
		}
	}

	if query == "" {
		usage("missing query")
	}

	if jsonOutput && rawOutput {
		usage("--json and --raw are mutually exclusive")
	}

	ctx := context.Background()

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

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

	server, err := registry.LookupDomain(resp.Body, query)
	if err != nil {
		fail(err)
	}

	client := rdap.NewClient()

	result, err := client.LookupDomain(ctx, server, query)
	if err != nil {
		fail(err)
	}

	if rawOutput {
		if _, err := os.Stdout.Write(result.Raw); err != nil {
			fail(err)
		}
		return
	}

	if jsonOutput {
		var value any

		if err := json.Unmarshal(result.Raw, &value); err != nil {
			fail(fmt.Errorf("decode RDAP response: %w", err))
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.SetEscapeHTML(false)

		if err := enc.Encode(value); err != nil {
			fail(fmt.Errorf("encode RDAP response: %w", err))
		}

		return
	}

	if err := output.Domain(os.Stdout, result.Domain); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "whose:", err)
	os.Exit(1)
}

func usage(message string) {
	fmt.Fprintln(os.Stderr, "whose:", message)
	fmt.Fprintln(os.Stderr, "usage: whose <query> [--json | --raw]")
	os.Exit(2)
}
