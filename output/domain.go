package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/rosethornbush/whose/rdap"
)

func Domain(w io.Writer, domain rdap.Domain) error {
	var b strings.Builder

	if domain.LDHName != "" {
		fmt.Fprintf(&b, "Domain:      %s\n", strings.ToLower(domain.LDHName))
	}

	if domain.UnicodeName != "" {
		fmt.Fprintf(&b, "Unicode:     %s\n", domain.UnicodeName)
	}

	for _, entity := range domain.Entities {
		if !entity.HasRole("registrar") {
			continue
		}

		name := entity.Name()
		if name == "" {
			name = entity.Organization()
		}

		if name != "" {
			fmt.Fprintf(&b, "Registrar:   %s\n", name)
			break
		}
	}

	for _, action := range []string{"registration", "last changed", "expiration"} {
		for _, event := range domain.Events {
			if event.Action != action {
				continue
			}

			date := formatDate(event.Date)

			switch action {
			case "registration":
				fmt.Fprintf(&b, "Created:     %s\n", date)
			case "last changed":
				fmt.Fprintf(&b, "Updated:     %s\n", date)
			case "expiration":
				fmt.Fprintf(&b, "Expires:     %s\n", date)
			}

			break
		}
	}

	if len(domain.Status) > 0 {
		fmt.Fprintln(&b)
		fmt.Fprintln(&b, "Status:")

		for _, status := range domain.Status {
			fmt.Fprintf(&b, "  %s\n", status)
		}
	}

	if len(domain.Nameservers) > 0 {
		fmt.Fprintln(&b)
		fmt.Fprintln(&b, "Nameservers:")

		for _, nameserver := range domain.Nameservers {
			fmt.Fprintf(&b, "  %s\n", strings.ToLower(nameserver.LDHName))
		}
	}

	if domain.SecureDNS != nil {
		fmt.Fprintln(&b)

		if domain.SecureDNS.DelegationSigned {
			fmt.Fprintln(&b, "DNSSEC:      signed")
		} else {
			fmt.Fprintln(&b, "DNSSEC:      unsigned")
		}
	}

	_, err := io.WriteString(w, b.String())
	return err
}
