package main

import (
	"net/url"
	"regexp"
)

var endpointPatterns = []*regexp.Regexp{
	// Most Amazon Web Services offer a Regional endpoint that you can use to make your requests.
	// In general, these endpoints support IPv4 traffic and they use the following syntax:
	// protocol://service-code.region-code.amazonaws.com
	regexp.MustCompile(`^[a-z]+://([a-z0-9-]+)\.[a-z0-9-]+\.amazonaws\.com(\.cn)?$`),

	// Some AWS services offer dual stack endpoints, so that you can access them using
	// either IPv4 or IPv6 requests. In general, the syntax of a dual stack endpoint is as follows:
	// protocol://service-code.region-code.api.aws
	regexp.MustCompile(`^[a-z]+://([a-z0-9-]+)\.[a-z0-9-]+\.api\.aws$`),

	// However, Amazon S3 uses the following syntax for its dual stack endpoints:
	// protocol://service-code.dualstack.region-code.amazonaws.com
	regexp.MustCompile(`^[a-z]+://([a-z0-9-]+)\.dualstack\.[a-z0-9-]+\.amazonaws\.com(\.cn)?$`),

	// Well, some endpoints do not fit into standard patterns
	regexp.MustCompile(`(s3|es)\.amazonaws\.com(\.cn)?$`),
}

func guessService(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	host := parsed.Host
	if host == "" {
		// Try to fall back to full input in case it was just a hostname
		host = rawURL
	}

	for _, p := range endpointPatterns {
		m := p.FindStringSubmatch(host)
		if len(m) >= 2 { //nolint:mnd
			return m[1]
		}
	}

	return ""
}
