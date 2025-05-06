package main

import (
	"net/url"
	"regexp"
)

// https://raw.githubusercontent.com/aws/aws-sdk-go-v2/master/codegen/smithy-aws-go-codegen/src/main/resources/software/amazon/smithy/aws/go/codegen/endpoints.json
//
//nolint:lll
var endpointPatterns = []*regexp.Regexp{
	// Most Amazon Web Services offer a Regional endpoint that you can use to make your requests.
	// In general, these endpoints support IPv4 traffic and they use the following syntax:
	// protocol://service-code.region-code.amazonaws.com
	regexp.MustCompile(`^(?P<service>[a-z0-9-]+?)(-fips)?\.(?P<region>[a-z0-9-]+)\.amazonaws\.com(\.cn)?$`),

	// Some AWS services offer dual stack endpoints, so that you can access them using
	// either IPv4 or IPv6 requests. In general, the syntax of a dual stack endpoint is as follows:
	// protocol://service-code.region-code.api.aws
	regexp.MustCompile(`^(?P<service>[a-z0-9-]+?)(-fips)?\.(?P<region>[a-z0-9-]+)\.api\.aws$`),

	// However, Amazon S3 uses the following syntax for its dual stack endpoints:
	// protocol://service-code.dualstack.region-code.amazonaws.com
	regexp.MustCompile(`^(?P<service>[a-z0-9-]+)\.dualstack\.(?P<region>[a-z0-9-]+)\.amazonaws\.com(\.cn)?$`),

	// Chinese dualstack endpoint
	regexp.MustCompile(`^(?P<service>[a-z0-9-]+?)(-fips)?\.(?P<region>[a-z0-9-]+)\.api\.amazonwebservices\.com\.cn$`),

	// Well, some endpoints do not fit into standard patterns
	regexp.MustCompile(`^(?P<service>s3)\.amazonaws\.com?$`),
	regexp.MustCompile(`^(?P<service>s3)\.(?P<region>[a-z0-9-]+)\.amazonaws\.com(\.cn)?$`),
	regexp.MustCompile(`^[a-z0-9-]+\.(?P<service>s3)\.(?P<region>[a-z0-9-]+)\.amazonaws\.com(\.cn)?$`),

	regexp.MustCompile(`^[a-z0-9-]+\.(?P<region>[a-z0-9-]+)\.(?P<service>es)\.amazonaws\.com(\.cn)?$`),
}

func guessServiceAndRegion(rawURL string) (service, region string) { //nolint: nonamedreturns
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", ""
	}

	host := parsed.Host
	if host == "" {
		// Try to fall back to full input in case it was just a hostname
		host = rawURL
	}

	for _, pattern := range endpointPatterns {
		match := pattern.FindStringSubmatch(host)
		if match == nil {
			continue
		}

		for i, name := range pattern.SubexpNames() {
			if name == "service" {
				service = match[i]
			} else if name == "region" {
				region = match[i]
			}
		}
	}

	return
}
