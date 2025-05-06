package main

import (
	"testing"
)

func TestGuessServiceAndRegion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input          string
		expectedSvc    string
		expectedRegion string
	}{
		// Standard endpoints
		{"https://s3.us-west-2.amazonaws.com", "s3", "us-west-2"},
		{"https://ec2.eu-central-1.amazonaws.com", "ec2", "eu-central-1"},
		{"https://s3.amazonaws.com", "s3", ""},

		// Dualstack endpoints
		{"https://s3.dualstack.eu-west-1.amazonaws.com", "s3", "eu-west-1"},
		{"https://ec2-fips.us-east-1.api.aws", "ec2", "us-east-1"},
		{"https://ec2.us-east-1.api.aws", "ec2", "us-east-1"},

		// Chinese endpoints
		{"https://ec2.cn-north-1.amazonaws.com.cn", "ec2", "cn-north-1"},
		{"https://s3.dualstack.cn-northwest-1.amazonaws.com.cn", "s3", "cn-northwest-1"},
		{"https://ec2-fips.cn-north-1.api.amazonwebservices.com.cn", "ec2", "cn-north-1"},

		// Elasticsearch-style endpoint
		{"https://search-mydomain.us-east-1.es.amazonaws.com", "es", "us-east-1"},

		// Legacy and corner-case patterns
		{"s3.amazonaws.com", "s3", ""},
		{"s3.us-east-1.amazonaws.com", "s3", "us-east-1"},
		{"mybucket.s3.us-west-2.amazonaws.com", "s3", "us-west-2"},

		// Invalid or unsupported
		{"https://example.com", "", ""},
		{"not-a-url", "", ""},
	}

	for _, tt := range tests {
		svc, region := guessServiceAndRegion(tt.input)
		if svc != tt.expectedSvc || region != tt.expectedRegion {
			t.Errorf("guessServiceAndRegion(%q) = (%q, %q), want (%q, %q)",
				tt.input, svc, region, tt.expectedSvc, tt.expectedRegion)
		}
	}
}
