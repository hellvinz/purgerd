package main

import (
	"flag"
	"fmt"
	zmq "github.com/alecthomas/gozmq"
	"log/syslog"
	"net"
	"os"
	"syscall"
	"time"
)

type clientPool struct {
	connections []net.Conn
}

var clients *clientPool
var context zmq.Context
var logger *syslog.Writer

func main() {
	// parse command line options
	var incomingAddress = flag.String("i", "0.0.0.0:8081", "incoming zmq purge address, eg: '0.0.0.0:8081'")
	var outgoingAddress = flag.String("o", "0.0.0.0:8080", "listening socket where purge message are sent to varnish reverse cli, eg: 0.0.0.0:8080")
	flag.Parse()

	// log to syslog
	logger, _ = syslog.New(syslog.LOG_INFO, "")

	// setup zmq
	context, _ = zmq.NewContext()
	defer context.Close()

	// the zmq REP socket where to send purge requests
	go setupPurgeReceiver(outgoingAddress)

	// we're ready to listen varnish cli connection
	setupPurgeSenderAndListen(incomingAddress)
}

//setupPurgeSenderAndListen create a clientPool and start listening to the socket where varnish cli connects
//when a client connects it is calling the handleConnection handler
func setupPurgeSenderAndListen(incomingAddress *string) {
	clients := &clientPool{connections: make([]net.Conn, 0)}
	go clients.ping()
	defer clients.close()
	ln, err := net.Listen("tcp", *incomingAddress)
	checkError(err)
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		logger.Info(fmt.Sprintln("New client: ", conn.RemoteAddr()))
		// flush the whole cache of the new client
		sendPurge(conn, ".*")
		// put it in the client pool
		clients.appendClient(conn)
		// connect client to the pubsub purge
		go connectClientToPusher(conn)
	}
	return
}

//setupPurgeReceiver set up the zmq REP socket where ban messages arrives
//when a purge pattern is received it dispatch it to a PUB socket
func setupPurgeReceiver(outgoingAddress *string) {
	receiver, _ := context.NewSocket(zmq.REP)
	defer receiver.Close()
	receiver.Bind("tcp://" + *outgoingAddress)

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

//connectClientToPusher is used to forward message received from the internal PUB socket to the client
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
			clients.removeClient(conn)
			break
		}
		logger.Debug(fmt.Sprintln("Client Purged", conn.RemoteAddr(), string(b)))
	}
}

//appendClient add a client connection to the clientPool
func (pool *clientPool) appendClient(conn net.Conn) {
	pool.connections = append(pool.connections, conn)
}

//removeClient remove a client connection from the clientPool
func (pool *clientPool) removeClient(conn net.Conn) {
	newConnections := make([]net.Conn, 0)
	for _, client := range pool.connections {
		if client != conn {
			newConnections = append(newConnections, client)
		}
	}
	pool.connections = newConnections
	return
}

//close close every connection with clients
func (pool *clientPool) close() {
	for _, client := range pool.connections {
		client.Close()
	}
}

//ping send ping message to every clients every 5 seconds
func (pool *clientPool) ping() {
	for {
		time.Sleep(5 * time.Second)
		fmt.Println(pool)
		for _, client := range pool.connections {
			n, err := client.Write([]byte("ping\n"))
			if n == 0 || err == syscall.EPIPE {
				logger.Debug(fmt.Sprintln("ping: client gone", client.RemoteAddr()))
				pool.removeClient(client)
				break
			}
		}
	}
}

//sendPurge send a purge message to a client
//it appends a ban.url to the pattern passed
func sendPurge(conn net.Conn, pattern string) (err error) {
	n, err := conn.Write([]byte("ban.url " + pattern + "\n"))
	if n == 0 {
		logger.Debug(fmt.Sprintln("failed to send message", conn.RemoteAddr()))
		err = syscall.EPIPE
	}
	return
}

//checkError basic error handling
func checkError(err error) {
	if err != nil {
		logger.Crit(fmt.Sprintln("Fatal error", err.Error()))
		os.Exit(1)
	}
}
