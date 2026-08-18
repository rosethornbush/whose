package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func LookupASN(r io.Reader, asn uint32) (string, error) {
	var reg registry

	if err := json.NewDecoder(r).Decode(&reg); err != nil {
		return "", fmt.Errorf("decode ASN bootstrap registry: %w", err)
	}

	matchedURL := ""

	for _, svc := range reg.Services {
		ranges := svc[0]
		urls := svc[1]

		if len(urls) == 0 {
			continue
		}

		serviceURL := urls[0]

		for _, candidate := range urls {
			if strings.HasPrefix(strings.ToLower(candidate), "https://") {
				serviceURL = candidate
				break
			}
		}

		for _, rawRange := range ranges {
			start, end, err := parseASNRange(rawRange)
			if err != nil {
				return "", fmt.Errorf("parse ASN bootstrap range %q: %w", rawRange, err)
			}

			if asn < start || asn > end {
				continue
			}

			if matchedURL != "" {
				return "", fmt.Errorf("multiple RDAP servers found for AS%d", asn)
			}

			matchedURL = serviceURL
		}
	}

	if matchedURL == "" {
		return "", fmt.Errorf("no RDAP server found for AS%d", asn)
	}

	return matchedURL, nil
}

func parseASNRange(value string) (uint32, uint32, error) {
	parts := strings.Split(value, "-")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return 0, 0, fmt.Errorf("expected start-end range")
	}

	start, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return 0, 0, err
	}

	end, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return 0, 0, err
	}

	if end < start {
		return 0, 0, fmt.Errorf("range end precedes start")
	}

	return uint32(start), uint32(end), nil
}
