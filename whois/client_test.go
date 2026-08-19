package whois

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestReferral(t *testing.T) {
	response := []byte(`
domain:       GG
organisation: Island Networks Ltd.

refer:        whois.gg
`)

	got := Referral(response)
	want := "whois.gg"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestReferralIsCaseInsensitive(t *testing.T) {
	response := []byte(`
REFER: whois.example
`)

	got := Referral(response)
	want := "whois.example"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestReferralMissing(t *testing.T) {
	response := []byte(`
domain: example
`)

	if got := Referral(response); got != "" {
		t.Fatalf("got %q, want empty string", got)
	}
}

func TestLookupChainSingleResponse(t *testing.T) {
	client := &Client{}

	client.lookup = func(
		ctx context.Context,
		server string,
		query string,
	) ([]byte, error) {
		return []byte("no referral"), nil
	}

	responses, err := client.LookupChain(
		context.Background(),
		"whois.example",
		"example.com",
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(responses) != 1 {
		t.Fatalf("got %d responses, want 1", len(responses))
	}

	if responses[0].Server != "whois.example" {
		t.Fatalf("got server %q, want %q", responses[0].Server, "whois.example")
	}
}

func TestLookupChainFollowsReferrals(t *testing.T) {
	client := &Client{}

	responsesByServer := map[string][]byte{
		"whois.one":   []byte("refer: whois.two"),
		"whois.two":   []byte("refer: whois.three"),
		"whois.three": []byte("done"),
	}

	client.lookup = func(
		ctx context.Context,
		server string,
		query string,
	) ([]byte, error) {
		body, ok := responsesByServer[server]
		if !ok {
			t.Fatalf("unexpected server %q", server)
		}

		return body, nil
	}

	responses, err := client.LookupChain(
		context.Background(),
		"whois.one",
		"example.com",
	)
	if err != nil {
		t.Fatal(err)
	}

	if len(responses) != 3 {
		t.Fatalf("got %d responses, want 3", len(responses))
	}

	want := []string{
		"whois.one",
		"whois.two",
		"whois.three",
	}

	for i, server := range want {
		if responses[i].Server != server {
			t.Fatalf(
				"response %d server = %q, want %q",
				i,
				responses[i].Server,
				server,
			)
		}
	}
}

func TestLookupChainDetectsLoop(t *testing.T) {
	client := &Client{}

	responsesByServer := map[string][]byte{
		"whois.one": []byte("refer: whois.two"),
		"whois.two": []byte("refer: whois.one"),
	}

	client.lookup = func(
		ctx context.Context,
		server string,
		query string,
	) ([]byte, error) {
		return responsesByServer[server], nil
	}

	responses, err := client.LookupChain(
		context.Background(),
		"whois.one",
		"example.com",
	)

	if err == nil {
		t.Fatal("expected referral loop error")
	}

	if len(responses) != 2 {
		t.Fatalf("got %d responses, want 2", len(responses))
	}
}

func TestLookupChainRejectsReferralDepthExceeded(t *testing.T) {
	client := &Client{}

	client.lookup = func(
		ctx context.Context,
		server string,
		query string,
	) ([]byte, error) {
		var n int

		if _, err := fmt.Sscanf(server, "whois.%d", &n); err != nil {
			t.Fatalf("unexpected server %q", server)
		}

		return fmt.Appendf(nil, "refer: whois.%d", n+1), nil
	}

	responses, err := client.LookupChain(
		context.Background(),
		"whois.1",
		"example.com",
	)

	if err == nil {
		t.Fatal("expected referral depth error")
	}

	if len(responses) != defaultMaxDepth {
		t.Fatalf(
			"got %d responses, want %d",
			len(responses),
			defaultMaxDepth,
		)
	}
}

func TestLookupChainPropagatesLookupError(t *testing.T) {
	client := &Client{}

	expected := errors.New("lookup failed")

	client.lookup = func(
		ctx context.Context,
		server string,
		query string,
	) ([]byte, error) {
		return nil, expected
	}

	responses, err := client.LookupChain(
		context.Background(),
		"whois.example",
		"example.com",
	)

	if !errors.Is(err, expected) {
		t.Fatalf("got error %v, want %v", err, expected)
	}

	if len(responses) != 0 {
		t.Fatalf("got %d responses, want 0", len(responses))
	}
}

func TestLookupChainPreservesResponsesBeforeLaterError(t *testing.T) {
	client := &Client{}

	expected := errors.New("lookup failed")

	client.lookup = func(
		ctx context.Context,
		server string,
		query string,
	) ([]byte, error) {
		switch server {
		case "whois.one":
			return []byte("refer: whois.two"), nil
		case "whois.two":
			return nil, expected
		default:
			t.Fatalf("unexpected server %q", server)
			return nil, nil
		}
	}

	responses, err := client.LookupChain(
		context.Background(),
		"whois.one",
		"example.com",
	)

	if !errors.Is(err, expected) {
		t.Fatalf("got error %v, want %v", err, expected)
	}

	if len(responses) != 1 {
		t.Fatalf("got %d responses, want 1", len(responses))
	}

	if responses[0].Server != "whois.one" {
		t.Fatalf("got server %q, want %q", responses[0].Server, "whois.one")
	}
}
