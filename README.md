# URL Health Checker

A simple Go program that periodically checks the health of a given URL by sending GET requests at a specified interval.

## Usage

To use this program, you need to provide two arguments: the URL to check and the interval in seconds between checks.

```bash
go run main.go <url> <interval_in_seconds>
```

### Example

```bash
go run main.go https://example.com 10
```

This command will check the health of `https://example.com` every 10 seconds.
