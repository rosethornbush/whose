package registry

import (
	"net/netip"
	"strings"
	"testing"
)

func TestLookupIPPrefersLongestPrefix(t *testing.T) {
	data := `{
		"services": [
			[
				["0.0.0.0/0"],
				["https://rdap.example/root/"]
			],
			[
				["1.0.0.0/8"],
				["https://rdap.example/one/"]
			],
			[
				["1.1.1.0/24"],
				["https://rdap.example/cloudflare/"]
			]
		]
	}`

	addr := netip.MustParseAddr("1.1.1.1")

	got, err := LookupIP(strings.NewReader(data), addr)
	if err != nil {
		t.Fatal(err)
	}

	want := "https://rdap.example/cloudflare/"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestLookupIPSupportsIPv6(t *testing.T) {
	data := `{
		"services": [
			[
				["2001:db8::/32"],
				["https://rdap.example/ipv6/"]
			]
		]
	}`

	addr := netip.MustParseAddr("2001:db8::1")

	got, err := LookupIP(strings.NewReader(data), addr)
	if err != nil {
		t.Fatal(err)
	}

	want := "https://rdap.example/ipv6/"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestLookupIPPrefersHTTPS(t *testing.T) {
	data := `{
		"services": [
			[
				["1.1.1.0/24"],
				[
					"http://rdap.example/",
					"https://rdap.example/"
				]
			]
		]
	}`

	addr := netip.MustParseAddr("1.1.1.1")

	got, err := LookupIP(strings.NewReader(data), addr)
	if err != nil {
		t.Fatal(err)
	}

	want := "https://rdap.example/"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestLookupIPNoMatch(t *testing.T) {
	data := `{
		"services": [
			[
				["192.0.2.0/24"],
				["https://rdap.example/"]
			]
		]
	}`

	addr := netip.MustParseAddr("198.51.100.1")

	_, err := LookupIP(strings.NewReader(data), addr)
	if err == nil {
		t.Fatal("expected no match")
	}
}

func TestLookupIPRejectsInvalidPrefixes(t *testing.T) {
	data := `{
		"services": [
			[
				["not-a-prefix"],
				["https://rdap.example/bad/"]
			],
			[
				["1.1.1.0/24"],
				["https://rdap.example/good/"]
			]
		]
	}`

	addr := netip.MustParseAddr("1.1.1.1")

	_, err := LookupIP(strings.NewReader(data), addr)
	if err == nil {
		t.Fatal("expected invalid prefix error")
	}
}
