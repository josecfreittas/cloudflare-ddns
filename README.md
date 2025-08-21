# Cloudflare DDNS

Keep a Cloudflare A/AAAA record in sync with the machine’s current public IP. Small binary, no daemon dependencies, and works in one-shot or looped mode.

Highlights
- IPv4 and IPv6 (choose either or both)
- Run once or on a fixed schedule (e.g., every 15m, 2h)
- Configuration via environment variables only

## Quick start

One-time update (runs and exits):

```bash
CF_TOKEN=YOUR_CLOUDFLARE_API_TOKEN \
CF_ZONE_ID=YOUR_ZONE_ID \
CF_HOST=example.com \
./cloudflare-ddns
```

Run as a simple loop (example: every 2 hours):

```bash
./cloudflare-ddns -duration 2h
```

Enable IPv6 updates (AAAA) alongside IPv4:

```bash
./cloudflare-ddns -ipv6=true -duration 1h
```

Notes
- The DNS record must already exist. This tool only updates the content value.
- For IPv6, create the AAAA record ahead of time and ensure the host has IPv6 connectivity.

## Configuration

Required environment variables
- CF_TOKEN: Cloudflare API Token with permission to edit DNS for the target zone
- CF_ZONE_ID: Cloudflare Zone ID for your domain
- CF_HOST: The fully-qualified DNS record to update (e.g., sub.example.com)

CLI flags
- -duration: Interval to check and update. Accepts Go duration strings like 30s, 5m, 2h30m. If unset or 0s, runs once and exits.
- -ipv4: Toggle IPv4 A record updates. Default: true
- -ipv6: Toggle IPv6 AAAA record updates. Default: false

## How it works

This tool reads your current public IP from well-known endpoints (IPv4 via checkip.amazonaws.com, IPv6 via v6.ident.me). It caches the last observed value and only talks to Cloudflare when the IP changes. When updating, it finds the first matching record for the host and preserves the existing proxied setting.

Limitations
- It won’t create records—only update existing ones.
- It targets a single host per run.

## Build

Using make (recommended):

```bash
# Build for the host platform (release)
make build

# Cross-compile
make build TARGET=linux-arm64
make build TARGET=windows-x64

# Debug build (disables optimizations and inlining)
make build DEBUG=1
# or
make build -- --debug
```

Artifacts are placed in:
```
dist/cloudflare-ddns-<os>-<arch>[.exe]
```

Build directly with Go (optional):

```bash
go build -trimpath -ldflags="-s -w" -o cloudflare-ddns .
```

## Troubleshooting

- Missing env vars: the program exits with a clear error if any of the three required variables are unset.
- “host not found”: ensure the A/AAAA record for CF_HOST already exists in the specified zone.
- No IPv6 updates: confirm your machine has IPv6 connectivity and the AAAA record exists.

## Acknowledgments

- This tool is heavily inspired by [hugomd/cloudflare-ddns](https://github.com/hugomd/cloudflare-ddns)!
