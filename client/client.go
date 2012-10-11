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
