package main

import (
	"fmt"
	zmq "github.com/alecthomas/gozmq"
	"log/syslog"
	"net"
	"os"
	"syscall"
	"time"
)

type clientPool []net.Conn

var clients clientPool
var context zmq.Context
var logger *syslog.Writer

func main() {
	logger, _ = syslog.New(syslog.LOG_INFO, "")
	context, _ = zmq.NewContext()
	defer context.Close()
	go setupPurgeReceiver()
	clients = make(clientPool, 0)
	go ping()
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

func setupPurgeReceiver() (err error) {
	receiver, _ := context.NewSocket(zmq.REP)
	defer receiver.Close()
	receiver.Bind("tcp://*:8080")

	pusher, _ := context.NewSocket(zmq.PUB)
	defer pusher.Close()
	pusher.Bind("inproc://pusher")
	for {
		b, _ := receiver.Recv(0)
		pusher.Send(b, 0)
		receiver.Send([]byte("ok"), 0)
	}
	return
}

func handleConnection(conn net.Conn) {
	logger.Info(fmt.Sprintln("New client: ", conn.RemoteAddr()))
	// flush the whole cache
	sendPurge(conn, ".*")
	// put it in the client pool
	clients = append(clients, conn)
	// connect client to the pubsub purge
	go connectClientToPusher(conn)
}

func connectClientToPusher(conn net.Conn) {
	puller, _ := context.NewSocket(zmq.SUB)
	puller.SetSockOptString(zmq.SUBSCRIBE, "")
	defer puller.Close()
	puller.Connect("inproc://pusher")
	for {
		b, _ := puller.Recv(0)
		err := sendPurge(conn, string(b))
		if err == syscall.EPIPE {
			logger.Info(fmt.Sprintln("client gone", conn.RemoteAddr()))
			remove(conn)
			break
		}
		logger.Debug(fmt.Sprintln("Client Purged", conn.RemoteAddr(), string(b)))
	}
}

func remove(conn net.Conn) {
	newClients := make(clientPool, len(clients)-1)
	for _, client := range clients {
		if client != conn {
			newClients = append(newClients, client)
		}
	}
	clients = newClients
	return
}

func (clients clientPool) close() {
	for _, client := range clients {
		client.Close()
	}
}

func ping() {
	for {
		time.Sleep(5 * time.Second)
		for _, client := range clients {
			n, err := client.Write([]byte("ping\n"))
			if n == 0 || err == syscall.EPIPE {
				logger.Debug(fmt.Sprintln("ping: client gone", client))
				remove(client)
				break
			}
		}
	}
}

func sendPurge(conn net.Conn, pattern string) (err error) {
	n, err := conn.Write([]byte("ban.url " + pattern + "\n"))
	if n == 0 {
		logger.Debug(fmt.Sprintln("failed to send message", conn))
		err = syscall.EPIPE
	}
	return
}

func checkError(err error) {
	if err != nil {
		logger.Crit(fmt.Sprintln("Fatal error", err.Error()))
		os.Exit(1)
	}
}
