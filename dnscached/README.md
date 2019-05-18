# dnscached

A simple example of using CoreDNS as a library. This embeds only the `bind`,
`cache`, `errors`, `forward` and `log` plugins, and only allows configuration
from flags.

## Usage

```
  dnscached [ OPTIONS ] [ DESTINATIONS ]
```

### OPTIONS

```
  -bind IP(s)
    	IP(s) to which to bind (default '127.0.0.1 ::1')
  -denial entries
    	Number of denial cache entries (default 9984)
  -dry-run
    	Prints out the internally generated Corefile and exits
  -log
    	Enable query logging
  -pidfile path
    	File path to write pid file
  -port number
    	Local port number to use (default 5300)
  -success entries
    	Number of success cache entries (default 9984)
  -ttl seconds
    	Maximum seconds to cache records, zero disables caching (default 60)
  -version
    	Show version
```

### DESTINATIONS

One or more forwarding destinations. Each can be a file in /etc/resolv.conf
format or a destination IP or IP:PORT, with or without a with or without a
protocol (leading `PROTO://`), with `dns` and `tls` as the supported `PROTO`
values. If omitted, "dns" is assumed as the protocol. The default destination is
/etc/resolv.conf.

## Examples

Run with default settings: port 5300 on localhost, forwarding to nameservers from /etc/resolv.conf.
```
  dnscached
```

Keep the cache contents for at most 30 seconds.

```
  dnscached -ttl 30
```

Use TLS for upstream traffic.

```
  dnscached tls://8.8.8.8 tls://8.8.4.4
```

Bind to two IPs:

```
  dnscached -bind "10.0.0.1 10.10.0.1"
```
