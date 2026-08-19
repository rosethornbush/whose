package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/rosethornbush/whose/rdap"
)

func ASN(w io.Writer, autnum rdap.Autnum) error {
	var b strings.Builder

	if autnum.Handle != "" {
		fmt.Fprintf(&b, "ASN:         %s\n", autnum.Handle)
	} else if autnum.StartAutnum == autnum.EndAutnum && autnum.StartAutnum != 0 {
		fmt.Fprintf(&b, "ASN:         AS%d\n", autnum.StartAutnum)
	}

	if autnum.Name != "" {
		fmt.Fprintf(&b, "Name:        %s\n", autnum.Name)
	}

	if name := entityName(autnum.Entities); name != "" {
		fmt.Fprintf(&b, "Org:         %s\n", name)
	}

	if autnum.StartAutnum != 0 &&
		autnum.EndAutnum != 0 &&
		autnum.StartAutnum != autnum.EndAutnum {
		fmt.Fprintf(
			&b,
			"Range:       AS%d-AS%d\n",
			autnum.StartAutnum,
			autnum.EndAutnum,
		)
	}

	if autnum.Type != "" {
		fmt.Fprintf(&b, "Type:        %s\n", autnum.Type)
	}

	if autnum.Country != "" {
		fmt.Fprintf(&b, "Country:     %s\n", autnum.Country)
	}

	for _, action := range []string{"registration", "last changed"} {
		for _, event := range autnum.Events {
			if event.Action != action {
				continue
			}

			switch action {
			case "registration":
				fmt.Fprintf(&b, "Created:     %s\n", event.Date)
			case "last changed":
				fmt.Fprintf(&b, "Updated:     %s\n", event.Date)
			}

			break
		}
	}

	if len(autnum.Status) > 0 {
		fmt.Fprintln(&b)
		fmt.Fprintln(&b, "Status:")

		for _, status := range autnum.Status {
			fmt.Fprintf(&b, "  %s\n", status)
		}
	}

	_, err := io.WriteString(w, b.String())
	return err
}
