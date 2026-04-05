# Operations and Troubleshooting

## Day-to-day commands

Render the final Compose configuration:

```bash
docker compose config
```

Start or restart:

```bash
docker compose up -d
```

Stop and remove:

```bash
docker compose down
```

Tail logs:

```bash
docker compose logs -f gluetun
```

## Useful checks

Check WireGuard handshake state:

```bash
docker exec -it gluetun wg show
```

Check outbound connectivity from inside the container:

```bash
docker exec -it gluetun wget -qO- --timeout=10 https://ipinfo.io/ip
```

Check local DNS inside the container:

```bash
docker exec -it gluetun nslookup github.com 127.0.0.1
```

## Common startup issues

### Blank variable warnings from Compose

Cause:

- `.env` does not exist, or required values are empty

Fix:

1. Copy [.env.example](../.env.example) to `.env`
2. Fill in the required provider and credential values
3. Re-run `docker compose config`

### WireGuard connects but traffic still times out

Likely causes:

- bad `WIREGUARD_PRIVATE_KEY`
- bad `WIREGUARD_ADDRESSES`
- upstream UDP blocking
- overly narrow server filters

Recommended checks:

1. clear `SERVER_COUNTRIES`, `SERVER_CITIES`, and `SERVER_HOSTNAMES`
2. confirm `wg show` reports a recent handshake
3. temporarily set `HEALTH_RESTART_VPN=off`

### DNS over TLS is unstable

Symptoms:

- TLS handshake resets to `1.1.1.1:853`
- intermittent lookup failures

Fix:

- switch to plain upstream DNS as documented in [Health Checks and DNS](./healthchecks-and-dns.md)

## Logging expectations

The control server listens on port `8000` by default. The HTTP proxy and Shadowsocks ports are published too, but their services stay off until enabled.

If you need persistent host-side log shipping for a larger stack, add a separate collector such as Vector or Promtail at the deployment level. That is outside the scope of this repo-level compose example.
