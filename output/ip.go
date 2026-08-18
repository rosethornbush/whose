package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/rosethornbush/whose/rdap"
)

func IP(w io.Writer, network rdap.IPNetwork) error {
	var b strings.Builder

	if network.Name != "" {
		fmt.Fprintf(&b, "Name:        %s\n", network.Name)
	}

	if network.Handle != "" {
		fmt.Fprintf(&b, "Handle:      %s\n", network.Handle)
	}

	if network.StartAddress != "" {
		fmt.Fprintf(&b, "Start:       %s\n", network.StartAddress)
	}

	if network.EndAddress != "" {
		fmt.Fprintf(&b, "End:         %s\n", network.EndAddress)
	}

	if network.IPVersion != "" {
		fmt.Fprintf(&b, "Version:     %s\n", network.IPVersion)
	}

	if network.Type != "" {
		fmt.Fprintf(&b, "Type:        %s\n", network.Type)
	}

	if network.Country != "" {
		fmt.Fprintf(&b, "Country:     %s\n", network.Country)
	}

	for _, event := range network.Events {
		switch event.Action {
		case "registration":
			fmt.Fprintf(&b, "Created:     %s\n", event.Date)
		case "last changed":
			fmt.Fprintf(&b, "Updated:     %s\n", event.Date)
		}
	}

	if len(network.Status) > 0 {
		fmt.Fprintln(&b)
		fmt.Fprintln(&b, "Status:")

		for _, status := range network.Status {
			fmt.Fprintf(&b, "  %s\n", status)
		}
	}

	_, err := io.WriteString(w, b.String())
	return err
}
