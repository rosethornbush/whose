package registry

import (
	"strings"
	"testing"
)

func TestLookupASNMatchesRange(t *testing.T) {
	data := `{
		"services": [
			[
				["13300-13400"],
				["https://rdap.example/"]
			]
		]
	}`

	got, err := LookupASN(strings.NewReader(data), 13335)
	if err != nil {
		t.Fatal(err)
	}

	want := "https://rdap.example/"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestLookupASNSupportsSingleASN(t *testing.T) {
	data := `{
		"services": [
			[
				["13335"],
				["https://rdap.example/"]
			]
		]
	}`

	got, err := LookupASN(strings.NewReader(data), 13335)
	if err != nil {
		t.Fatal(err)
	}

	want := "https://rdap.example/"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestLookupASNPrefersHTTPS(t *testing.T) {
	data := `{
		"services": [
			[
				["13300-13400"],
				[
					"http://rdap.example/",
					"https://rdap.example/"
				]
			]
		]
	}`

	got, err := LookupASN(strings.NewReader(data), 13335)
	if err != nil {
		t.Fatal(err)
	}

	want := "https://rdap.example/"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestLookupASNNoMatch(t *testing.T) {
	data := `{
		"services": [
			[
				["13300-13400"],
				["https://rdap.example/"]
			]
		]
	}`

	_, err := LookupASN(strings.NewReader(data), 64512)
	if err == nil {
		t.Fatal("expected no match")
	}
}

func TestLookupASNRejectsInvalidRange(t *testing.T) {
	data := `{
		"services": [
			[
				["not-a-range"],
				["https://rdap.example/"]
			]
		]
	}`

	_, err := LookupASN(strings.NewReader(data), 13335)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLookupASNRejectsReversedRange(t *testing.T) {
	data := `{
		"services": [
			[
				["13400-13300"],
				["https://rdap.example/"]
			]
		]
	}`

	_, err := LookupASN(strings.NewReader(data), 13335)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLookupASNRejectsOverlappingRanges(t *testing.T) {
	data := `{
		"services": [
			[
				["13000-14000"],
				["https://rdap.example/one/"]
			],
			[
				["13300-13400"],
				["https://rdap.example/two/"]
			]
		]
	}`

	_, err := LookupASN(strings.NewReader(data), 13335)
	if err == nil {
		t.Fatal("expected error")
	}
}
