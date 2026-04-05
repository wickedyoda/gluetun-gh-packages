# Setup and Docker Compose

## Files

- [docker-compose.yml](../docker-compose.yml) is the runnable Compose definition.
- [.env.example](../.env.example) is the template for your local `.env`.
- The mounted Gluetun state directory defaults to `./gluetun` and is ignored by git.

## Quick start

1. Copy `.env.example` to `.env`.
2. Edit `.env` with your VPN provider and credentials.
3. Start the container:

```bash
docker compose up -d
```

1. Inspect the resolved configuration:

```bash
docker compose config
```

1. Follow logs:

```bash
docker compose logs -f gluetun
```

## Required settings

Every deployment needs:

- `VPN_SERVICE_PROVIDER`
- `VPN_TYPE`

Then choose one credential path:

- OpenVPN:
  - `OPENVPN_USER`
  - `OPENVPN_PASSWORD`
- WireGuard:
  - `WIREGUARD_PRIVATE_KEY`
  - `WIREGUARD_ADDRESSES`

## Exposed ports

The bundled Compose file publishes:

- `8000/tcp` for the control server
- `8888/tcp` for the built-in HTTP proxy
- `8388/tcp` and `8388/udp` for Shadowsocks

The HTTP proxy and Shadowsocks services stay disabled unless you explicitly turn them on in `.env`.

## Why this Compose file uses NET_RAW

The default health configuration uses ICMP for the smaller between-cycle checks:

- `HEALTH_SMALL_CHECK_TYPE=icmp`
- `HEALTH_ICMP_TARGET_IPS=1.1.1.1,8.8.8.8`

That requires `NET_RAW`, so the Compose file adds:

- `NET_ADMIN`
- `NET_RAW`

If you do not want ICMP health checks, you can remove `NET_RAW` and change `HEALTH_SMALL_CHECK_TYPE`.
