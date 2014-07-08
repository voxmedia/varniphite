varniphite
==========

Sending Varnish 4.0 stats results to graphite never been so easy.

What
===

Every 10 seconds it runs `varnishstat`, parses the JSON result, does a minimal cleanup and sends to Graphite.

How
===

This app doesn't daemonize itself, so I suggest using it under `runit` or similar. All output goes to `STDOUT`

Installation
=============

```
$ go get github.com/voxmedia/varniphite
```

In case of the server not having `go` installed, you can compile/cross-compile (if you are on MacOS X) for linux:

```
$ GOOS=linux go build github.com/voxmedia/varniphite
```

And copy the generated static binary to your servers.

Parameters
==========

Parameter|Default|Function
---------|-------|----------
`-i`|`10`           |interval between runs
`-H`|`localhost`    |Graphite server host
`-m`|`varnish.stats`|Path to be appended to metrics
`-p`|`2003`         |Graphite server port

Example `runit` script
======================

```bash
#!/bin/sh
HOSTNAME=$(hostname)
exec 2>&1
exec /opt/go/bin/varniphite -H stats.example.com -p 2003 -m "varnish_stats.$HOSTNAME" -i 10
```
