Purgerd
=======

This tool forward a purge request received to a pool of varnish connected via the Varnish reverse cli.

It uses [ban](https://www.varnish-cache.org/docs/4.0/reference/varnish-cli.html#ban-expressions) and basically sends to the connected varnish the string it receives prefixed by ban req.url ~

It is implement as a PubSub, varnishes subscribes to purges on the purgerd daemon, and clients push purges to the purgerd daemon. When a varnish vanish (server reboot...) it is removed from subscribers.

Usage
=====

```
Usage of purgerd:
  -i="0.0.0.0:8111": socket where purge messages are sent, '0.0.0.0:8111'
  -o="0.0.0.0:1118": listening socket where purge message are sent to varnish reverse cli, 0.0.0.0:1118
  -p=false: purge all the varnish cache on connection
  -s="": path of the file containing the varnish secret
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

Netcat

to receive purges
```
echo "200 0" | nc localhost 1118
```

to send purges
```
echo "mypurge" | nc localhost 8111
```

Requirements
============

Go: http://golang.org/
Gb: http://getgb.io/

Install
=======

This project uses [gb](http://getgb.io/) to build the project

```
git clone https://github.com/hellvinz/purgerd.git
cd purgerd
gb generate (optional: if you want to regenerate varnish cli parser. Need [ragel](http://www.colm.net/open-source/ragel/))
gb build all
```

Logging
=======

the purgerd logs to syslog

to have some stats: killall -USR1 purgerd
