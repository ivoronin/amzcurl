# amzcurl
![GitHub release (with filter)](https://img.shields.io/github/v/release/ivoronin/amzcurl)
[![Go Report Card](https://goreportcard.com/badge/github.com/ivoronin/amzcurl)](https://goreportcard.com/report/github.com/ivoronin/amzcurl)
![GitHub last commit (branch)](https://img.shields.io/github/last-commit/ivoronin/amzcurl/main)
![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/ivoronin/amzcurl/goreleaser.yml)
![GitHub top language](https://img.shields.io/github/languages/top/ivoronin/amzcurl)

**`amzcurl`** is the thinnest possible wrapper around [`curl`](https://curl.se/), designed to transparently inject AWS SigV4 authentication using credentials discovered via the [AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2).

It allows you to use `curl` as-is to make signed requests to AWS APIs with minimal friction.

> You’ll need at least curl version 8.5 for everything to work properly - the newer, the better, as many issues with AWS SigV4 have been fixed in recent releases.

## 🧩 What It Does

- 🔍 **Discovers AWS credentials** from your environment (`~/.aws/config`, env vars, EC2/ECS/IAM, etc.)
- 🔎 **Auto-detects the AWS service name and region** based on the request URL (optional override)
- 🪪 **Injects SigV4 signing flags** into `curl` using `--aws-sigv4` and temporary credentials
- ✅ **Passes everything else to `curl` untouched**

That's it.

No extra abstractions. Just automatic AWS signing.

## 📦 Installation

### Download

https://github.com/ivoronin/amzcurl/releases

### Build from source

```bash
git clone https://github.com/ivoronin/amzcurl.git
cd amzcurl
go build -o amzcurl
```

### Or install via go install

```
go install github.com/ivoronin/amzcurl@latest
```

## 🚀 Usage

```
amzcurl [--profile PROFILE] [--region REGION] [--service SERVICE] <curl args...>
```

### Examples

```
# Auto-detect service from URL (S3)
amzcurl https://my-bucket.s3.amazonaws.com/file.txt

# Explicit service and region
amzcurl --service dynamodb --region us-west-2 \
  -X POST https://dynamodb.us-west-2.amazonaws.com/ \
  -H "Content-Type: application/x-amz-json-1.0" \
  -H "X-Amz-Target: DynamoDB_20120810.ListTables" \
  -d '{}'
```
