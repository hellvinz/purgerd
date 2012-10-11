Purger
======

This tool forward a purge request received to a pool of varnish connected via the reverse cli

Usage
=====

```
Usage of ./purger: 
  -i="0.0.0.0:8081": incoming zmq purge address, eg: '0.0.0.0:8081'                                 
  -o="0.0.0.0:8080": listening socket where purge message are sent to varnish reverse cli, eg: 0.0.0.0:8
```

Run purger from $GOCODE/bin/purger. With no options it will listen to purge request on 0.0.0.0:8081 with a REP 0MQ socket.
Start varnish with the -M option to make it connect to the purger. (ex: -M localhost:8081 if you're running varnish on the same box)

Client example
==============

Ruby
```
require 'ffi-rzmq'

puts "new context"
context = ZMQ::Context.new
socket = context.socket ZMQ::REQ
socket.connect('tcp://127.0.0.1:8080')
socket.send_string('.*',0)
msg = ""
socket.recv_string(msg, 0)
puts msg
```

Go
```
package main

import (
	"fmt"
  "flag"
  zmq "github.com/alecthomas/gozmq"
)

func main(){
  var pattern = flag.String("p", ".*", "the url pattern you want to purge")
  flag.Parse()
    context, _ := zmq.NewContext()
    defer context.Close()
    requester, _ := context.NewSocket(zmq.REQ)
    defer requester.Close()
    requester.Connect("tcp://localhost:8080")
    fmt.Println("Sending pattern:", *pattern)
    requester.Send([]byte(*pattern), 0)
    for{
      b, _ := requester.Recv(0)
      if string(b) == "ok" {
        fmt.Println("Purge done of pattern: ", *pattern)
        break
      }
    }
}
```

Requirements
============

0MQ: http://www.zeromq.org/
Go: http://golang.org/

Install
=======

if you install 0MQ in a non-standard directory, for example /opt/local, export first:

```
export CGO_LDFLAGS=-L/opt/local/lib
export CGO_CFLAGS=-I/opt/local/include
```

then

`
go get github.com/hellvinz/purger
`
