package whois

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

const defaultMaxDepth = 5

type Client struct {
	dialer net.Dialer
}

type Response struct {
	Server string
	Body   []byte
}

func NewClient() *Client {
	return &Client{
		dialer: net.Dialer{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) Lookup(
	ctx context.Context,
	server string,
	query string,
) ([]byte, error) {
	conn, err := c.dialer.DialContext(
		ctx,
		"tcp",
		net.JoinHostPort(server, "43"),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to WHOIS server %q: %w", server, err)
	}
	defer conn.Close()

	if deadline, ok := ctx.Deadline(); ok {
		if err := conn.SetDeadline(deadline); err != nil {
			return nil, fmt.Errorf("set WHOIS deadline: %w", err)
		}
	} else {
		if err := conn.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
			return nil, fmt.Errorf("set WHOIS deadline: %w", err)
		}
	}

	if _, err := fmt.Fprintf(conn, "%s\r\n", query); err != nil {
		return nil, fmt.Errorf("write WHOIS query: %w", err)
	}

	body, err := io.ReadAll(conn)
	if err != nil {
		return nil, fmt.Errorf("read WHOIS response: %w", err)
	}

	return body, nil
}

func (c *Client) LookupChain(
	ctx context.Context,
	server string,
	query string,
) ([]Response, error) {
	visited := make(map[string]bool)
	responses := make([]Response, 0, defaultMaxDepth)

	for len(responses) < defaultMaxDepth {
		server = strings.TrimSpace(server)

		if server == "" {
			break
		}

		key := strings.ToLower(server)

		if visited[key] {
			return responses, fmt.Errorf("WHOIS referral loop detected at %q", server)
		}

		visited[key] = true

		body, err := c.Lookup(ctx, server, query)
		if err != nil {
			return responses, err
		}

		responses = append(responses, Response{
			Server: server,
			Body:   body,
		})

		next := Referral(body)
		if next == "" {
			break
		}

		server = next
	}

	if len(responses) == defaultMaxDepth {
		last := responses[len(responses)-1]

		if Referral(last.Body) != "" {
			return responses, fmt.Errorf(
				"WHOIS referral depth exceeded %d",
				defaultMaxDepth,
			)
		}
	}

	return responses, nil
}

func Referral(response []byte) string {
	scanner := bufio.NewScanner(strings.NewReader(string(response)))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		if value == "" {
			continue
		}

		switch {
		case strings.EqualFold(key, "refer"):
			return normalizeServer(value)

		case strings.EqualFold(key, "whois"):
			return normalizeServer(value)

		case strings.EqualFold(key, "whois server"):
			return normalizeServer(value)

		case strings.EqualFold(key, "registrar whois server"):
			return normalizeServer(value)
		}
	}

	return ""
}

func normalizeServer(value string) string {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "whois://")
	value = strings.TrimSuffix(value, "/")

	if host, _, err := net.SplitHostPort(value); err == nil {
		return host
	}

	return value
}
