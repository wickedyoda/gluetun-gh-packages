# Environment Variables

This page documents the variables used by the bundled Compose setup. For the exhaustive image defaults, refer to [Dockerfile](../Dockerfile).

## Container identity

- `TZ`: container timezone used for timestamps
- `PUID`: runtime user id inside the container
- `PGID`: runtime group id inside the container

## Compose host settings

- `GLUETUN_STATE_DIR`: host directory mounted to `/gluetun`
- `GLUETUN_CONTROL_PORT`: host port mapped to the control server
- `GLUETUN_HTTP_PROXY_PORT`: host port mapped to the HTTP proxy
- `GLUETUN_SHADOWSOCKS_PORT`: host port mapped to Shadowsocks

## VPN selection

- `VPN_SERVICE_PROVIDER`: provider slug such as `mullvad`, `pia`, `ivpn`, `protonvpn`
- `VPN_TYPE`: `openvpn` or `wireguard`

## OpenVPN

- `OPENVPN_USER`
- `OPENVPN_PASSWORD`

Use these only when `VPN_TYPE=openvpn`.

## WireGuard

- `WIREGUARD_PRIVATE_KEY`
- `WIREGUARD_ADDRESSES`
- `WIREGUARD_MTU`

Use these when `VPN_TYPE=wireguard`.

## Server filtering

- `SERVER_COUNTRIES`
- `SERVER_CITIES`
- `SERVER_HOSTNAMES`

Leave them blank to let Gluetun choose from all eligible servers for the provider.

## Firewall

- `FIREWALL_OUTBOUND_SUBNETS`

Use this to allow LAN access or other explicitly permitted outbound destinations outside the tunnel.

## DNS

- `DNS_UPSTREAM_RESOLVER_TYPE`
- `DNS_UPSTREAM_RESOLVERS`
- `DNS_UPSTREAM_PLAIN_ADDRESSES`
- `BLOCK_MALICIOUS`
- `BLOCK_ADS`
- `BLOCK_SURVEILLANCE`

Recommended defaults in this repo:

- `DNS_UPSTREAM_RESOLVER_TYPE=DoT`
- block malicious domains on
- ads and surveillance blocking off unless you explicitly want it

## Health checks

- `HEALTH_RESTART_VPN`
- `HEALTH_TARGET_ADDRESSES`
- `HEALTH_SMALL_CHECK_TYPE`
- `HEALTH_ICMP_TARGET_IPS`

The bundled defaults use raw IP addresses instead of hostnames for startup health checks. That avoids bootstrap failures where the tunnel is up but DNS has not stabilized yet.

## Built-in services

- `HTTP_CONTROL_SERVER_LOG`
- `HTTPPROXY`
- `HTTPPROXY_LOG`
- `SHADOWSOCKS`
- `SHADOWSOCKS_LOG`
- `SHADOWSOCKS_PASSWORD`

The Compose file publishes the ports, but the services stay disabled until their corresponding enable flags are set to `on`.

## Updater

- `UPDATER_PERIOD`

This controls how often server metadata is refreshed. The local default is `24h`.
