package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type registry struct {
	Services []service `json:"services"`
}

type service [2][]string

func LookupDomain(r io.Reader, domain string) (string, error) {
	var reg registry

	if err := json.NewDecoder(r).Decode(&reg); err != nil {
		return "", fmt.Errorf("decode bootstrap registry: %w", err)
	}

	domain = strings.TrimSuffix(strings.ToLower(domain), ".")

	bestMatchLength := -1
	bestURL := ""

	for _, svc := range reg.Services {
		names := svc[0]
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

		for _, name := range names {
			name = strings.ToLower(name)

			if name != "" && domain != name && !strings.HasSuffix(domain, "."+name) {
				continue
			}

			if len(name) > bestMatchLength {
				bestMatchLength = len(name)
				bestURL = serviceURL
			}
		}
	}

	if bestURL == "" {
		return "", fmt.Errorf("no RDAP server found for %q", domain)
	}

	return bestURL, nil
}
