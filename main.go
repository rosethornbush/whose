package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/netip"
	"os"
	"strconv"
	"strings"

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
		lookupASN(ctx, q.Value, jsonOutput, rawOutput)

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
	data, err := registry.Fetch(
		ctx,
		"dns",
		"https://data.iana.org/rdap/dns.json",
	)
	if err != nil {
		fail(err)
	}

	server, err := registry.LookupDomain(bytes.NewReader(data), domain)
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

	registryName := "ipv4"
	registryURL := "https://data.iana.org/rdap/ipv4.json"

	if addr.Is6() {
		registryName = "ipv6"
		registryURL = "https://data.iana.org/rdap/ipv6.json"
	}

	data, err := registry.Fetch(ctx, registryName, registryURL)
	if err != nil {
		fail(err)
	}

	server, err := registry.LookupIP(bytes.NewReader(data), addr)
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

func lookupASN(
	ctx context.Context,
	value string,
	jsonOutput bool,
	rawOutput bool,
) {
	n, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		fail(fmt.Errorf("invalid ASN %q: %w", value, err))
	}

	asn := uint32(n)

	data, err := registry.Fetch(
		ctx,
		"asn",
		"https://data.iana.org/rdap/asn.json",
	)
	if err != nil {
		fail(err)
	}

	server, err := registry.LookupASN(bytes.NewReader(data), asn)
	if err != nil {
		fail(err)
	}

	client := rdap.NewClient()

	result, err := client.LookupASN(ctx, server, asn)
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

	if err := output.ASN(os.Stdout, result.Autnum); err != nil {
		fail(err)
	}
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
