package whois

import "testing"

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
