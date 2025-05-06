package main

import (
	"testing"
)

func TestGuessService(t *testing.T) {
	tests := []struct {
		name     string
		rawURL   string
		expected string
	}{
		{
			name:     "Standard AWS endpoint",
			rawURL:   "https://sqs.us-west-2.amazonaws.com",
			expected: "sqs",
		},
		{
			name:     "Standard AWS endpoint with .cn",
			rawURL:   "https://s3.cn-north-1.amazonaws.com.cn",
			expected: "s3",
		},
		{
			name:     "Dual stack endpoint",
			rawURL:   "https://sqs.us-west-2.api.aws",
			expected: "sqs",
		},
		{
			name:     "S3 dualstack endpoint",
			rawURL:   "https://s3.dualstack.us-west-2.amazonaws.com",
			expected: "s3",
		},
		{
			name:     "ES endpoint (non-standard)",
			rawURL:   "https://es.amazonaws.com",
			expected: "es",
		},
		{
			name:     "S3 endpoint (non-standard) with .cn",
			rawURL:   "https://s3.amazonaws.com.cn",
			expected: "s3",
		},
		{
			name:     "Malformed URL",
			rawURL:   "http://[::1]:namedport",
			expected: "",
		},
		{
			name:     "Plain hostname without scheme",
			rawURL:   "s3.dualstack.us-west-2.amazonaws.com",
			expected: "s3",
		},
		{
			name:     "Unknown pattern",
			rawURL:   "https://example.com",
			expected: "",
		},
		{
			name:     "Empty input",
			rawURL:   "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := guessService(tt.rawURL)
			if got != tt.expected {
				t.Errorf("guessService(%q) = %q; want %q", tt.rawURL, got, tt.expected)
			}
		})
	}
}
