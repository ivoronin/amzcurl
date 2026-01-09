package main

import (
	"fmt"
	"strings"
)

// Returns the argument at args[*idx+1] and advances the index.
// If no such argument exists, returns "", ErrMissingArgumentValue.
func shiftArg(name string, args *[]string, idx *int) (string, error) {
	*idx++
	if *idx >= len(*args) {
		return "", fmt.Errorf("%w: %s", ErrMissingArgumentValue, name)
	}

	return (*args)[*idx], nil
}

// Returns the first string that isn't empty or an empty string.
func coalesce(items ...string) string {
	for _, i := range items {
		if i != "" {
			return i
		}
	}

	return ""
}

func parseFlags(args []string) (string, string, string, []string, error) {
	var (
		profile                  string
		region, regionOverride   string
		service, serviceOverride string
		curlArgs                 []string
	)

	for idx := 0; idx < len(args); idx++ { //nolint:intrange
		if strings.HasPrefix(args[idx], "https://") || strings.HasPrefix(args[idx], "http://") {
			service, region = guessServiceAndRegion(args[idx])
		}

		var err error

		switch args[idx] {
		case "--profile":
			profile, err = shiftArg("--profile", &args, &idx)
		case "--region":
			regionOverride, err = shiftArg("--region", &args, &idx)
		case "--service":
			serviceOverride, err = shiftArg("--service", &args, &idx)
		default:
			curlArgs = append(curlArgs, args[idx])
		}

		if err != nil {
			return "", "", "", nil, err
		}
	}

	service = coalesce(serviceOverride, service)
	region = coalesce(regionOverride, region)

	return profile, region, service, curlArgs, nil
}
