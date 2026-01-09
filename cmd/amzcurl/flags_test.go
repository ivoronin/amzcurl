package main

import (
	"reflect"
	"testing"
)

func TestParseFlags(t *testing.T) { //nolint:funlen
	t.Parallel()

	tests := []struct {
		name         string
		args         []string
		wantProfile  string
		wantRegion   string
		wantService  string
		wantCurlArgs []string
		expectErr    bool
	}{
		{
			name:         "no flags, just URL",
			args:         []string{"https://service.region.amazonaws.com", "-X", "GET"},
			wantProfile:  "",
			wantRegion:   "region",
			wantService:  "service",
			wantCurlArgs: []string{"https://service.region.amazonaws.com", "-X", "GET"},
		},
		{
			name:         "with --profile",
			args:         []string{"--profile", "dev", "https://service.region.amazonaws.com"},
			wantProfile:  "dev",
			wantRegion:   "region",
			wantService:  "service",
			wantCurlArgs: []string{"https://service.region.amazonaws.com"},
		},
		{
			name:         "with --region and --service override",
			args:         []string{"--region", "us-west-2", "--service", "s3", "https://service.region.amazonaws.com"},
			wantProfile:  "",
			wantRegion:   "us-west-2",
			wantService:  "s3",
			wantCurlArgs: []string{"https://service.region.amazonaws.com"},
		},
		{
			name:      "missing profile argument",
			args:      []string{"--profile"},
			expectErr: true,
		},
		{
			name:         "extra args preserved",
			args:         []string{"--profile", "dev", "https://service.region.amazonaws.com", "-H", "Accept: */*"},
			wantProfile:  "dev",
			wantRegion:   "region",
			wantService:  "service",
			wantCurlArgs: []string{"https://service.region.amazonaws.com", "-H", "Accept: */*"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			profile, region, service, curlArgs, err := parseFlags(test.args)

			if test.expectErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if profile != test.wantProfile {
				t.Errorf("profile = %q, want %q", profile, test.wantProfile)
			}

			if region != test.wantRegion {
				t.Errorf("region = %q, want %q", region, test.wantRegion)
			}

			if service != test.wantService {
				t.Errorf("service = %q, want %q", service, test.wantService)
			}

			if !reflect.DeepEqual(curlArgs, test.wantCurlArgs) {
				t.Errorf("curlArgs = %v, want %v", curlArgs, test.wantCurlArgs)
			}
		})
	}
}
