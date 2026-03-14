# URL Availability Checker

This repository contains a Go program (`main.go`) that checks the availability of URLs passed as command-line arguments. It performs DNS lookups and HTTP requests concurrently using goroutines.

## Purpose

The `main.go` program demonstrates how to use concurrency in Go to efficiently check the status of multiple URLs. It resolves hostnames to IP addresses and sends HTTP GET requests to determine if the URLs are up or down.

## Prerequisites

- Go 1.25.1 or newer installed. Verify with:

```bash
go version
```

## Usage

1. Run the program directly with `go run`:

```bash
go run main.go https://example.com https://google.com
```

2. Alternatively, build the program into a binary:

```bash
go build -o url-checker main.go
./url-checker https://example.com https://google.com
```

## Example Output

```plaintext
[UP]   status:200      url:https://example.com          ip:[93.184.216.34]
[UP]   status:200      url:https://google.com           ip:[142.250.190.14]
```

- `[UP]` indicates the URL is reachable with a status code ≤ 299.
- `[DOWN]` indicates the URL is unreachable or returned a status code > 299.

## Notes

- The program uses a timeout of 5 seconds for HTTP requests to avoid hanging.
- DNS lookups are performed to resolve hostnames to IP addresses.
- Ensure the URLs provided include the scheme (e.g., `https://`).

