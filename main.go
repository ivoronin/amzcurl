package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func fatalf(format string, args ...interface{}) int {
	fmt.Fprintf(os.Stderr, "amzcurl: "+format+"\n", args...)

	return 1
}

func mustArg(name string, args *[]string, idx *int) string {
	*idx++
	if *idx >= len(*args) {
		fatalf("missing value for %s", name)
	}

	return (*args)[*idx]
}

func parseFlags(args []string) (profile, region, service string, curlArgs []string) { //nolint:nonamedreturns
	serviceOverride := ""

	for idx := 0; idx < len(args); idx++ { //nolint:intrange
		if strings.HasPrefix(args[idx], "https://") || strings.HasPrefix(args[idx], "http://") {
			service = guessService(args[idx])
		}

		switch args[idx] {
		case "--profile":
			profile = mustArg("--profile", &args, &idx)
		case "--region":
			region = mustArg("--region", &args, &idx)
		case "--service":
			serviceOverride = mustArg("--service", &args, &idx)
		default:
			curlArgs = append(curlArgs, args[idx])
		}
	}

	if serviceOverride != "" {
		service = serviceOverride
	}

	return
}

func buildCurlConfig(cfg aws.Config, regionOverride, service string) ([]string, error) {
	region := cfg.Region
	if regionOverride != "" {
		region = regionOverride
	}

	creds, err := cfg.Credentials.Retrieve(context.Background())
	if err != nil {
		return nil, err
	}

	curlConfigLines := []string{
		fmt.Sprintf(`--user "%s:%s"`, creds.AccessKeyID, creds.SecretAccessKey),
		fmt.Sprintf(`--aws-sigv4 "aws:amz:%s:%s"`, region, service),
	}

	if creds.SessionToken != "" {
		curlConfigLines = append(curlConfigLines, fmt.Sprintf(`-H "x-amz-security-token: %s"`, creds.SessionToken))
	}

	return curlConfigLines, nil
}

func amzcurl() int {
	profile, region, service, curlArgs := parseFlags(os.Args[1:])
	if service == "" {
		return fatalf("--service is required")
	}

	cfgOpts := []func(*config.LoadOptions) error{}
	if profile != "" {
		cfgOpts = append(cfgOpts, config.WithSharedConfigProfile(profile))
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), cfgOpts...)
	if err != nil {
		return fatalf("failed to load AWS config: %v", err)
	}

	curlConfigLines, err := buildCurlConfig(cfg, region, service)
	if err != nil {
		return fatalf("failed to build config: %v", err)
	}

	tmpFile, err := os.CreateTemp("", "amzcurl")
	if err != nil {
		return fatalf("failed to create temp file: %v", err)
	}

	defer os.Remove(tmpFile.Name())

	for _, line := range curlConfigLines {
		fmt.Fprintln(tmpFile, line)
	}

	tmpFile.Close()

	// #nosec G204
	cmd := exec.Command("curl", append([]string{"--config", tmpFile.Name()}, curlArgs...)...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}

		return fatalf("failed to execute curl: %v", err)
	}

	return 0
}

func main() {
	os.Exit(amzcurl())
}
