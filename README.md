Purger
======

This tool forward a purge request received to a pool of varnish connected via the reverse cli

Usage
=====

`
Usage of ./purger: 
  -i="0.0.0.0:8081": incoming zmq purge address, eg: '0.0.0.0:8081'                                 
  -o="0.0.0.0:8080": listening socket where purge message are sent to varnish reverse cli, eg: 0.0.0.0:8
`

Requirements
============

0MQ: http://www.zeromq.org/
Go: http://golang.org/

Install
=======

if you install 0MQ in a non-standard directory, for example /opt/local, export first:

`
export CGO_LDFLAGS=-L/opt/local/lib
export CGO_CFLAGS=-I/opt/local/include
`

then

`
go get github.com/hellvinz/purger
`
