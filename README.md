# certwatch

SSL Certificate Watcher - monitors SSL/TLS certificate expiration dates.

## Purpose

Track SSL certificate validity and expiration for domains. Provides renewal scheduling recommendations.

## Installation

```bash
go build -o certwatch ./cmd/certwatch
```

## Usage

```bash
certwatch <domain1> <domain2> ...
```

### Examples

```bash
# Check single domain
certwatch example.com

# Check multiple domains
certwatch example.com api.example.com mail.example.com

# Check with specific ports
certwatch example.com:443
```

## Output

```
=== CERTIFICATE WATCHER ===

Domain: example.com
Issuer: Let's Encrypt
Valid From: 2024-01-15
Valid Until: 2025-01-15
Days Remaining: 365
Status: OK

=== RENEWAL SCHEDULE ===
Recommended cron jobs for auto-renewal:

30 days before expiry:
  0 0 1 * * certbot renew --force-renewal

7 days before expiry:
  0 0 1 * * certbot renew
```

## Status Levels

- OK: Certificate valid with >30 days remaining
- EXPIRING_SOON: Certificate valid with 7-30 days remaining
- CRITICAL: Certificate valid with <7 days remaining
- EXPIRED: Certificate has expired

## Dependencies

- Go 1.21+
- github.com/fatih/color

## Build and Run

```bash
# Build
go build -o certwatch ./cmd/certwatch

# Run
go run ./cmd/certwatch example.com api.example.com
```

## License

MIT