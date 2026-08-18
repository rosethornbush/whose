package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/netip"
	"os"
	"strings"
	"time"

	"github.com/rosethornbush/whose/output"
	"github.com/rosethornbush/whose/query"
	"github.com/rosethornbush/whose/rdap"
	"github.com/rosethornbush/whose/registry"
)

func main() {
	var (
		queryInput string
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

			if queryInput != "" {
				usage("expected exactly one query")
			}

			queryInput = arg
		}
	}

	if queryInput == "" {
		usage("missing query")
	}

	if jsonOutput && rawOutput {
		usage("--json and --raw are mutually exclusive")
	}

	q, err := query.Parse(queryInput)
	if err != nil {
		fail(err)
	}

	ctx := context.Background()

	switch q.Kind {
	case query.Domain:
		lookupDomain(ctx, q.Value, jsonOutput, rawOutput)

	case query.IP:
		lookupIP(ctx, q.Value, jsonOutput, rawOutput)

	case query.ASN:
		fail(fmt.Errorf("ASN queries are not supported yet"))

	default:
		fail(fmt.Errorf("unsupported query type"))
	}
}

func lookupDomain(
	ctx context.Context,
	domain string,
	jsonOutput bool,
	rawOutput bool,
) {
	resp := fetchRegistry(ctx, "https://data.iana.org/rdap/dns.json")
	defer resp.Body.Close()

	server, err := registry.LookupDomain(resp.Body, domain)
	if err != nil {
		fail(err)
	}

	client := rdap.NewClient()

	result, err := client.LookupDomain(ctx, server, domain)
	if err != nil {
		fail(err)
	}

	if rawOutput {
		printRaw(result.Raw)
		return
	}

	if jsonOutput {
		printJSON(result.Raw)
		return
	}

	if err := output.Domain(os.Stdout, result.Domain); err != nil {
		fail(err)
	}
}

func lookupIP(
	ctx context.Context,
	address string,
	jsonOutput bool,
	rawOutput bool,
) {
	addr, err := netip.ParseAddr(address)
	if err != nil {
		fail(fmt.Errorf("invalid IP address %q: %w", address, err))
	}

	registryURL := "https://data.iana.org/rdap/ipv4.json"

	if addr.Is6() {
		registryURL = "https://data.iana.org/rdap/ipv6.json"
	}

	resp := fetchRegistry(ctx, registryURL)
	defer resp.Body.Close()

	server, err := registry.LookupIP(resp.Body, addr)
	if err != nil {
		fail(err)
	}

	client := rdap.NewClient()

	result, err := client.LookupIP(ctx, server, address)
	if err != nil {
		fail(err)
	}

	if rawOutput {
		printRaw(result.Raw)
		return
	}

	if jsonOutput {
		printJSON(result.Raw)
		return
	}

	if err := output.IP(os.Stdout, result.Network); err != nil {
		fail(err)
	}
}

func fetchRegistry(ctx context.Context, registryURL string) *http.Response {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		registryURL,
		nil,
	)
	if err != nil {
		fail(fmt.Errorf("create IANA registry request: %w", err))
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		fail(fmt.Errorf("request IANA registry: %w", err))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resp.Body.Close()
		fail(fmt.Errorf("IANA registry returned %s", resp.Status))
	}

	return resp
}

func printJSON(raw []byte) {
	var value any

	if err := json.Unmarshal(raw, &value); err != nil {
		fail(fmt.Errorf("decode RDAP response: %w", err))
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)

	if err := enc.Encode(value); err != nil {
		fail(fmt.Errorf("encode RDAP response: %w", err))
	}
}

func printRaw(raw []byte) {
	if _, err := os.Stdout.Write(raw); err != nil {
		fail(fmt.Errorf("write RDAP response: %w", err))
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
