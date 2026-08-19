# whose

A fast, human-friendly replacement for `whois`, powered by RDAP.

`whose` looks up domains, IP addresses, and ASNs using RDAP where available, with WHOIS fallback for domains that need it.

```sh
whose github.com
whose 1.1.1.1
whose AS13335
```

## Features

- domain, IPv4, IPv6, and ASN lookups
- human-readable output by default
- formatted RDAP JSON with `--json`
- raw protocol responses with `--raw`
- WHOIS fallback for domains without RDAP support
- WHOIS referral following
- IDN support
- cached IANA RDAP bootstrap registries

## Installation

Requires Go 1.24 or later.

```sh
go install github.com/rosethornbush/whose@latest
```

Or build from source:

```sh
git clone https://github.com/rosethornbush/whose
cd whose
go build .
```

## Usage

```text
whose <query> [--json | --raw]
```

### Domains

```sh
whose github.com
```

```text
Domain:      github.com
Registrar:   MarkMonitor Inc.
Created:     2007-10-09 18:20:50 UTC
Updated:     2024-09-07 09:16:32 UTC
Expires:     2026-10-09 18:20:50 UTC

Status:
  client delete prohibited
  client transfer prohibited
  client update prohibited

Nameservers:
  dns1.p08.nsone.net
  dns2.p08.nsone.net
  dns3.p08.nsone.net
  dns4.p08.nsone.net
  ns-1283.awsdns-32.org
  ns-1707.awsdns-21.co.uk
  ns-421.awsdns-52.com
  ns-520.awsdns-01.net

DNSSEC:      unsigned
```

Internationalized domain names are supported:

```sh
whose münchen.de
```

### IP addresses

```sh
whose 1.1.1.1
```

```text
Name:        APNIC-LABS
Org:         APNIC Research and Development
Handle:      1.1.1.0 - 1.1.1.255
Start:       1.1.1.0
End:         1.1.1.255
Version:     v4
Type:        ASSIGNED PORTABLE
Country:     AU
Created:     2011-08-10 23:12:35 UTC
Updated:     2023-04-26 22:57:58 UTC

Status:
  active
```

IPv6 works the same way:

```sh
whose 2606:4700:4700::1111
```

### ASNs

```sh
whose AS13335
```

```text
ASN:         AS13335
Name:        CLOUDFLARENET
Org:         Cloudflare, Inc.
Created:     2010-07-14 22:35:57 UTC
Updated:     2017-02-17 23:04:32 UTC

Status:
  active
```

ASN input is case-insensitive:

```sh
whose as13335
```

### JSON output

Use `--json` to pretty-print the response from the authoritative RDAP server:

```sh
whose github.com --json
```

The response is standard RDAP JSON. `whose` does not transform it into a custom schema.

### Raw output

Use `--raw` to print the protocol response without reformatting it:

```sh
whose github.com --raw
```

For RDAP lookups, this is the original JSON response. For WHOIS fallback, this is the WHOIS text response.

## How it works

`whose` uses IANA's RDAP bootstrap registries to find the authoritative server for domains, IP addresses, and ASNs.

```text
query
  -> IANA RDAP bootstrap
  -> authoritative RDAP server
  -> response
```

The bootstrap registries are cached locally for 24 hours. If a refresh fails, `whose` can use a valid stale cache instead.

Some TLDs do not publish an RDAP service through IANA. For those domains, `whose` falls back to WHOIS:

```text
domain
  -> whois.iana.org
  -> follow referrals
  -> authoritative WHOIS server
```

WHOIS referral responses are preserved in order so the fallback stays close to traditional `whois` behavior.

`--json` is unavailable for WHOIS fallback because WHOIS has no structured JSON response.

## Development

Run the test suite:

```sh
go test ./...
```

Run the race detector:

```sh
go test -race ./...
```

Run vet:

```sh
go vet ./...
```

Build:

```sh
go build .
```

## License

[GPL-3.0](./LICENSE)
