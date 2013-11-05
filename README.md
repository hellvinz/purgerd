Purgerd
=======

This tool forward a purge request received to a pool of varnish connected via the reverse cli

Usage
=====

```
Usage of purgerd:
  -i="0.0.0.0:8111": socket where purge messages are sent, '0.0.0.0:8111'
  -o="0.0.0.0:1118": listening socket where purge message are sent to varnish reverse cli, 0.0.0.0:1118
  -p=false: purge all the varnish cache on connection
  -s="": varnish secret
  -v=false: display version
```

Run purgerd from $GOCODE/bin/purgerd. With no options it will listen to purge requests on 0.0.0.0:8111.
Start varnish with the -M option to make it connect to the purger. (ex: -M localhost:1118 if you're running varnish on the same box)
If your varnish cli needs [authentication](https://www.varnish-cache.org/trac/wiki/CLI#Authentication:Thegorydetails) pass the password with -s

Client example
==============

Ruby
```
require "socket"
socket = TCPSocket.new('127.0.0.1',8111)
socket.write('.*')
socket.close()
```

Requirements
============

Go: http://golang.org/

Install
=======

go get github.com/hellvinz/purgerd

Logging
=======

the purgerd logs to syslog
