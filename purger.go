package main

import (
	"fmt"
	"net"
	"os"
	"time"
  zmq "github.com/alecthomas/gozmq"
)

type clientPool []net.Conn

var clients clientPool

func main() {
	go setupPurgeReceiver()
	clients = make(clientPool, 0)
	defer clients.close()
	ln, err := net.Listen("tcp", ":8081")
	checkError(err)
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		go handleConnection(conn)
	}

}

func setupPurgeReceiver() (err error){
    context, _ := zmq.NewContext()
    defer context.Close()
    receiver, _ := context.NewSocket(zmq.REP)
    defer receiver.Close()
    receiver.Bind("tcp://*:8080")
    for {
      b, _ := receiver.Recv(0)
      clients.dispatchPurge(string(b))
      receiver.Send([]byte("ok"),0)
    }
    return
}

func handleConnection(conn net.Conn) {
	fmt.Println("New client", conn.RemoteAddr())
	// flush the whole cache
	sendPurge(conn, ".*")
	// put it in the client pool
	clients = append(clients, conn)
}

func (clients clientPool) dispatchPurge(pattern string) {
  fmt.Println("Going to purge:", pattern)
	for _, client := range clients {
		err := sendPurge(client, pattern)
		checkError(err)
	}
}

func (clients clientPool) close() {
	for _, client := range clients {
		client.Close()
	}
}

func sendPurge(conn net.Conn, pattern string) (err error) {
	n, err := conn.Write([]byte("ban.url " + pattern + "\n"))
	if n == 0 {
		fmt.Println("failed to send message")
	}
	return
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
