# Health Checks and DNS

## Default behavior

The bundled Compose file sets:

- `HEALTH_RESTART_VPN=on`
- `HEALTH_TARGET_ADDRESSES=1.1.1.1:443,8.8.8.8:443`
- `HEALTH_SMALL_CHECK_TYPE=icmp`
- `HEALTH_ICMP_TARGET_IPS=1.1.1.1,8.8.8.8`

This intentionally avoids DNS-based startup health checks. Using IP targets makes the first connectivity test independent from DNS resolution.

## Why not use hostnames by default

If the tunnel comes up but DNS-through-tunnel has not settled yet, hostname-based health checks can fail even though the VPN path is otherwise healthy. That commonly produces restart loops like:

- WireGuard setup completes
- DNS lookups time out
- health check fails
- `HEALTH_RESTART_VPN=on` restarts the tunnel

Using IP targets first removes that bootstrap dependency.

## DNS defaults

The Compose file keeps:

- `DNS_UPSTREAM_RESOLVER_TYPE=DoT`
- `DNS_UPSTREAM_RESOLVERS=` blank, so Gluetun falls back to its built-in defaults

If your network or ISP interferes with DNS over TLS, switch to plain upstream addresses:

```env
DNS_UPSTREAM_RESOLVER_TYPE=plain
DNS_UPSTREAM_PLAIN_ADDRESSES=1.1.1.1,1.0.0.1,8.8.8.8,8.8.4.4
```

## When to disable automatic restart temporarily

For debugging, turn this off:

```env
HEALTH_RESTART_VPN=off
```

That lets the container stay up long enough to inspect live state instead of restarting in a tight loop.

## ICMP and capabilities

The small health check defaults to ICMP, which requires:

- `NET_RAW`

If you remove `NET_RAW`, change the small health check to a non-ICMP mode before starting the container.
