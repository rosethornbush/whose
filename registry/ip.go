package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/netip"
	"strings"
)

func LookupIP(r io.Reader, addr netip.Addr) (string, error) {
	var reg registry

	if err := json.NewDecoder(r).Decode(&reg); err != nil {
		return "", fmt.Errorf("decode IP bootstrap registry: %w", err)
	}

	bestBits := -1
	bestURL := ""

	for _, svc := range reg.Services {
		prefixes := svc[0]
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

		for _, rawPrefix := range prefixes {
			prefix, err := netip.ParsePrefix(rawPrefix)
			if err != nil {
				return "", fmt.Errorf("parse IP bootstrap prefix %q: %w", rawPrefix, err)
			}

			if !prefix.Contains(addr) {
				continue
			}

			if prefix.Bits() > bestBits {
				bestBits = prefix.Bits()
				bestURL = serviceURL
			}
		}
	}

	if bestURL == "" {
		return "", fmt.Errorf("no RDAP server found for %q", addr)
	}

	return bestURL, nil
}
