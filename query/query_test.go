package query

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		input string
		kind  Kind
		value string
	}{
		{"example.com", Domain, "example.com"},
		{"EXAMPLE.COM", Domain, "example.com"},
		{"example.com.", Domain, "example.com"},
		{"com", Domain, "com"},
		{"foo-bar.example", Domain, "foo-bar.example"},
		{"asia.com", Domain, "asia.com"},
		{"ask.com", Domain, "ask.com"},
		{"as", Domain, "as"},

		{"1.1.1.1", IP, "1.1.1.1"},
		{"2606:4700:4700::1111", IP, "2606:4700:4700::1111"},

		{"AS13335", ASN, "13335"},
		{"as13335", ASN, "13335"},

		{"münchen.de", Domain, "xn--mnchen-3ya.de"},
		{"bücher.example", Domain, "xn--bcher-kva.example"},
		{"MÜNCHEN.DE", Domain, "xn--mnchen-3ya.de"},
		{"💩.com", Domain, "xn--ls8h.com"},
		{"münchen。de。", Domain, "xn--mnchen-3ya.de"},
		{"aſ13335", Domain, "as13335"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := Parse(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			if got.Kind != tt.kind {
				t.Fatalf("kind = %v, want %v", got.Kind, tt.kind)
			}

			if got.Value != tt.value {
				t.Fatalf("value = %q, want %q", got.Value, tt.value)
			}
		})
	}
}

func TestParseInvalidDomains(t *testing.T) {
	tests := []string{
		".",
		"foo..bar",
		"hello world.com",
		"-foo.com",
		"foo-.com",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			if _, err := Parse(input); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestParseInvalidASN(t *testing.T) {
	_, err := Parse("AS4294967296")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseRejectsScopedIPv6(t *testing.T) {
	_, err := Parse("fe80::1%en0")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestParseInvalidIDN(t *testing.T) {
	tests := []string{
		"\u200d.com",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			if _, err := Parse(input); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
