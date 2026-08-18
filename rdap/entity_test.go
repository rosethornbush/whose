package rdap

import (
	"encoding/json"
	"testing"
)

func TestEntityHasRole(t *testing.T) {
	entity := Entity{
		Roles: []string{"registrar", "technical"},
	}

	if !entity.HasRole("registrar") {
		t.Fatal("expected registrar role")
	}

	if entity.HasRole("abuse") {
		t.Fatal("did not expect abuse role")
	}
}

func TestEntityName(t *testing.T) {
	entity := Entity{
		VCardArray: json.RawMessage(`[
			"vcard",
			[
				["version", {}, "text", "4.0"],
				["fn", {}, "text", "Example Registrar"]
			]
		]`),
	}

	got := entity.Name()
	want := "Example Registrar"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestEntityOrganization(t *testing.T) {
	entity := Entity{
		VCardArray: json.RawMessage(`[
			"vcard",
			[
				["version", {}, "text", "4.0"],
				["org", {}, "text", "Example Networks"]
			]
		]`),
	}

	got := entity.Organization()
	want := "Example Networks"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestEntityOrganizationArray(t *testing.T) {
	entity := Entity{
		VCardArray: json.RawMessage(`[
			"vcard",
			[
				["org", {}, "text", ["Example Corp", "Network Operations"]]
			]
		]`),
	}

	got := entity.Organization()
	want := "Example Corp Network Operations"

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestEntityNameMissingVCard(t *testing.T) {
	entity := Entity{}

	if got := entity.Name(); got != "" {
		t.Fatalf("got %q, want empty string", got)
	}
}

func TestEntityNameIgnoresMalformedVCard(t *testing.T) {
	entity := Entity{
		VCardArray: json.RawMessage(`["not-vcard", []]`),
	}

	if got := entity.Name(); got != "" {
		t.Fatalf("got %q, want empty string", got)
	}
}

func TestEntityNameIgnoresNonTextValue(t *testing.T) {
	entity := Entity{
		VCardArray: json.RawMessage(`[
			"vcard",
			[
				["fn", {}, "uri", "https://example.com"]
			]
		]`),
	}

	if got := entity.Name(); got != "" {
		t.Fatalf("got %q, want empty string", got)
	}
}
