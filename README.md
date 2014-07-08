varniphite
==========

Sending Varnish 4.0 stats results to graphite never been so easy.

What
===

Every 10 seconds it runs `varnishstat`, parses the JSON result, does a minimal cleanup and sends to Graphite.

How
===

This app doesn't daemonize itself, so I suggest using it under `runit` or similar. All output goes to `STDOUT`

Parameters
==========

Parameter|Default|Function
---------|-------|----------
`-i`|`10`           |interval between runs
`-H`|`localhost`    |Graphite server host
`-m`|`varnish.stats`|Path to be appended to metrics
`-p`|`2003`         |Graphite server port
