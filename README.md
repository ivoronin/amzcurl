# amzcurl

Curl wrapper that injects AWS SigV4 authentication using credentials from the AWS SDK credential chain

[![CI](https://github.com/ivoronin/amzcurl/actions/workflows/test.yml/badge.svg)](https://github.com/ivoronin/amzcurl/actions/workflows/test.yml)
[![Release](https://img.shields.io/github/v/release/ivoronin/amzcurl)](https://github.com/ivoronin/amzcurl/releases)

[Overview](#overview) · [Features](#features) · [Installation](#installation) · [Usage](#usage) · [Configuration](#configuration) · [Requirements](#requirements) · [License](#license)

```bash
# Before: manually constructing SigV4 with curl
curl --aws-sigv4 "aws:amz:us-west-2:s3" \
  --user "$AWS_ACCESS_KEY_ID:$AWS_SECRET_ACCESS_KEY" \
  -H "x-amz-security-token: $AWS_SESSION_TOKEN" \
  https://my-bucket.s3.us-west-2.amazonaws.com/file.txt

# After: credentials and signing handled automatically
amzcurl https://my-bucket.s3.us-west-2.amazonaws.com/file.txt
```

## Overview

amzcurl discovers AWS credentials using the AWS SDK for Go v2 credential chain (environment variables, shared credentials file, EC2/ECS instance metadata). It parses the request URL to detect the AWS service and region, then passes signing flags (`--aws-sigv4`, `--user`, session token header) to curl via a temporary config file. All other arguments pass through to curl unchanged.

## Features

- Discovers credentials from full AWS SDK chain (env vars, `~/.aws/credentials`, `~/.aws/config`, EC2/ECS IMDS, SSO)
- Auto-detects service name and region from standard AWS endpoint URLs
- Supports named profiles via `--profile`
- Supports manual region and service override via `--region` and `--service`
- Handles session tokens for temporary credentials (IAM roles, SSO)
- Supports standard, dual-stack, FIPS, and Chinese region endpoints
- All curl arguments pass through unchanged

## Installation

### GitHub Releases

Download from [Releases](https://github.com/ivoronin/amzcurl/releases).

### Homebrew

```bash
brew install ivoronin/tap/amzcurl
```

### Go Install

```bash
go install github.com/ivoronin/amzcurl/cmd/amzcurl@latest
```

### Build from Source

```bash
git clone https://github.com/ivoronin/amzcurl.git
cd amzcurl
go build -o amzcurl ./cmd/amzcurl
```

## Usage

```
amzcurl [--profile PROFILE] [--region REGION] [--service SERVICE] [curl args...]
```

### S3

```bash
# List bucket contents
amzcurl https://my-bucket.s3.us-west-2.amazonaws.com/

# Download file
amzcurl -o file.txt https://my-bucket.s3.us-west-2.amazonaws.com/file.txt

# Upload file
amzcurl -X PUT -T file.txt https://my-bucket.s3.us-west-2.amazonaws.com/file.txt
```

### DynamoDB

```bash
amzcurl -X POST https://dynamodb.us-west-2.amazonaws.com/ \
  -H "Content-Type: application/x-amz-json-1.0" \
  -H "X-Amz-Target: DynamoDB_20120810.ListTables" \
  -d '{}'
```

### Named Profile

```bash
amzcurl --profile production https://my-bucket.s3.us-west-2.amazonaws.com/
```

### Explicit Region and Service

For non-standard endpoints or when auto-detection fails:

```bash
amzcurl --service execute-api --region us-east-1 \
  https://abc123.execute-api.us-east-1.amazonaws.com/prod/resource
```

### Version

```bash
amzcurl --version
```

## Configuration

amzcurl uses the AWS SDK for Go v2 default credential chain. No tool-specific configuration is needed.

### Credential Chain Order

1. Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_SESSION_TOKEN`)
2. Shared credentials file (`~/.aws/credentials`)
3. Shared config file (`~/.aws/config`) with profile support
4. EC2/ECS instance metadata service (IMDS)
5. SSO credentials when configured

### AWS Environment Variables

| Variable | Description |
|----------|-------------|
| `AWS_ACCESS_KEY_ID` | Access key ID |
| `AWS_SECRET_ACCESS_KEY` | Secret access key |
| `AWS_SESSION_TOKEN` | Session token for temporary credentials |
| `AWS_PROFILE` | Named profile to use |
| `AWS_REGION` | Default region |
| `AWS_CONFIG_FILE` | Path to config file (default: `~/.aws/config`) |
| `AWS_SHARED_CREDENTIALS_FILE` | Path to credentials file (default: `~/.aws/credentials`) |

## Requirements

- curl 8.5 or newer (required for `--aws-sigv4` flag; newer versions recommended as SigV4 bugs have been fixed in recent releases)
- Valid AWS credentials accessible via the SDK credential chain

## License

[GPL-3.0](LICENSE)
