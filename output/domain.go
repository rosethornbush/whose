package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/rosethornbush/whose/rdap"
)

func Domain(w io.Writer, domain rdap.Domain) error {
	if domain.LDHName != "" {
		fmt.Fprintf(w, "Domain:      %s\n", strings.ToLower(domain.LDHName))
	}

	if domain.UnicodeName != "" {
		fmt.Fprintf(w, "Unicode:     %s\n", domain.UnicodeName)
	}

	for _, event := range domain.Events {
		switch event.Action {
		case "registration":
			fmt.Fprintf(w, "Created:     %s\n", event.Date)
		case "last changed":
			fmt.Fprintf(w, "Updated:     %s\n", event.Date)
		case "expiration":
			fmt.Fprintf(w, "Expires:     %s\n", event.Date)
		}
	}

	if len(domain.Status) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Status:")

		for _, status := range domain.Status {
			fmt.Fprintf(w, "  %s\n", status)
		}
	}

	if len(domain.Nameservers) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Nameservers:")

		for _, nameserver := range domain.Nameservers {
			fmt.Fprintf(w, "  %s\n", strings.ToLower(nameserver.LDHName))
		}
	}

	if domain.SecureDNS != nil {
		fmt.Fprintln(w)

		if domain.SecureDNS.DelegationSigned {
			fmt.Fprintln(w, "DNSSEC:      signed")
		} else {
			fmt.Fprintln(w, "DNSSEC:      unsigned")
		}
	}

	return nil
}
