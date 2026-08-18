package query

import (
	"fmt"
	"net/netip"
	"strconv"
	"strings"
)

type Kind int

const (
	Domain Kind = iota
	IP
	ASN
)

type Query struct {
	Kind  Kind
	Value string
}

func Parse(input string) (Query, error) {
	if addr, err := netip.ParseAddr(input); err == nil {
		return Query{
			Kind:  IP,
			Value: addr.String(),
		}, nil
	}

	upper := strings.ToUpper(input)

	if strings.HasPrefix(upper, "AS") {
		value := strings.TrimPrefix(upper, "AS")

		if value != "" {
			allDigits := true

			for _, r := range value {
				if r < '0' || r > '9' {
					allDigits = false
					break
				}
			}

			if allDigits {
				n, err := strconv.ParseUint(value, 10, 32)
				if err != nil {
					return Query{}, fmt.Errorf("invalid ASN %q", input)
				}

				return Query{
					Kind:  ASN,
					Value: strconv.FormatUint(n, 10),
				}, nil
			}
		}
	}

	domain, err := parseDomain(input)
	if err != nil {
		return Query{}, err
	}

	return Query{
		Kind:  Domain,
		Value: domain,
	}, nil
}

func parseDomain(input string) (string, error) {
	domain := strings.ToLower(strings.TrimSuffix(input, "."))

	if domain == "" {
		return "", fmt.Errorf("invalid domain %q", input)
	}

	if len(domain) > 253 {
		return "", fmt.Errorf("invalid domain %q", input)
	}

	for _, label := range strings.Split(domain, ".") {
		if label == "" || len(label) > 63 {
			return "", fmt.Errorf("invalid domain %q", input)
		}

		if label[0] == '-' || label[len(label)-1] == '-' {
			return "", fmt.Errorf("invalid domain %q", input)
		}

		for _, r := range label {
			if (r >= 'a' && r <= 'z') ||
				(r >= '0' && r <= '9') ||
				r == '-' {
				continue
			}

			return "", fmt.Errorf("invalid domain %q", input)
		}
	}

	return domain, nil
}
