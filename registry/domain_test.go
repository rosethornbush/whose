package registry

import (
	"strings"
	"testing"
)

func TestLookupDomainPrefersLongestMatch(t *testing.T) {
	data := `{
		"services": [
			[
				["com"],
				["https://rdap.example/com/"]
			],
			[
				["example.com"],
				["https://rdap.example/example/"]
			]
		]
	}`

	got, err := LookupDomain(strings.NewReader(data), "foo.example.com")
	if err != nil {
		t.Fatal(err)
	}

	want := "https://rdap.example/example/"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestLookupDomainDoesNotMatchPartialLabel(t *testing.T) {
	data := `{
		"services": [
			[
				["example.com"],
				["https://rdap.example/"]
			]
		]
	}`

	_, err := LookupDomain(strings.NewReader(data), "notexample.com")
	if err == nil {
		t.Fatal("expected no match")
	}
}

func TestLookupDomainUsesRootFallback(t *testing.T) {
	data := `{
		"services": [
			[
				[""],
				["https://rdap.example/root/"]
			]
		]
	}`

	got, err := LookupDomain(strings.NewReader(data), "foo.example")
	if err != nil {
		t.Fatal(err)
	}

	want := "https://rdap.example/root/"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestLookupDomainPrefersHTTPS(t *testing.T) {
	data := `{
		"services": [
			[
				["com"],
				[
					"http://rdap.example/",
					"https://rdap.example/"
				]
			]
		]
	}`

	got, err := LookupDomain(strings.NewReader(data), "example.com")
	if err != nil {
		t.Fatal(err)
	}

	want := "https://rdap.example/"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
