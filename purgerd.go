package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io/ioutil"
	"log/syslog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var logger *syslog.Writer

func main() {

	// log to syslog
	logger, _ = syslog.New(syslog.LOG_INFO, "")

	// parse command line options
	incomingAddress := flag.String("i", "0.0.0.0:8111", "socket where purge messages are sent, '0.0.0.0:8111'")
	outgoingAddress := flag.String("o", "0.0.0.0:1118", "listening socket where purge message are sent to varnish reverse cli, 0.0.0.0:1118")
	version := flag.Bool("v", false, "display version")
	purgeOnStartUp := flag.Bool("p", false, "purge all the varnish cache on connection")
	secret := flag.String("s", "", "varnish secret")
	flag.Parse()
	if *version {
		printVersion()
		os.Exit(0)
	}

	publisher := NewPublisher()

	go monitorSignals(publisher)

	go setupPurgeReceiver(incomingAddress, publisher)

	// we're ready to listen varnish cli connection
	setupPurgeSenderAndListen(outgoingAddress, *purgeOnStartUp, publisher, secret)
}

//setupPurgeSenderAndListen start listening to the socket where varnish cli connects
//when a client connects it calls handleClient
func setupPurgeSenderAndListen(outgoingAddress *string, purgeOnStartup bool, publisher *Publisher, secret *string) {
	ln, err := net.Listen("tcp", *outgoingAddress)
	checkError(err,logger)
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		logger.Info(fmt.Sprintln("New client:", reverseName(conn)))

		// connect client to the pubsub purge
		go handleClient(conn, publisher, purgeOnStartup, secret)
	}
	return
}

//setupPurgeReceiver set up the tcp socket where ban messages come
//when a purge pattern is received it dispatches it to a Pub object
func setupPurgeReceiver(incomingAddress *string, publisher *Publisher) {
	receiver, err := net.Listen("tcp", *incomingAddress)
	checkError(err,logger)

	go func() {
		for {
			time.Sleep(5 * time.Second)
			publisher.Pub([]byte("ping"))
		}
	}()
	for {
		conn, err := receiver.Accept()
		checkError(err,logger)
		go func(c net.Conn) {
			defer conn.Close()
			b, err := ioutil.ReadAll(conn)
			if err != nil {
				logger.Info(fmt.Sprintln("Client connection error:", err))
			} else {
				logger.Info(fmt.Sprintln("<-", reverseName(conn), string(b)))
				publisher.Pub(bytes.TrimSpace(b))
				conn.Write([]byte("OK\n"))
			}
		}(conn)
	}
	return
}

//handleClient is used to forward message received to the client
func handleClient(conn net.Conn, publisher *Publisher, purgeOnStartup bool, secret *string) {
    defer conn.Close()

	// check if client need auth
	message := make([]byte, 512)
	conn.Read(message)
	cli := Cliparser(message)
	if cli.status == 107 {
		if *secret == "" {
			logger.Crit("Client varnish asked for a secret, provide one with -s")
			return
		}
		challenge := cli.body[:32]
		response := fmt.Sprintf("%s\n%s\n%s\n", challenge, *secret, challenge)
		hasher := sha256.New()
		hasher.Write([]byte(response))
		conn.Write([]byte(fmt.Sprintf("auth %s\n", hex.EncodeToString(hasher.Sum(nil)))))
	}

	if purgeOnStartup {
		// flush the whole cache of the new client
		conn.Write([]byte("ban.url .*\n"))
	}

	// wait for purges
	wait := make(chan bool, 1)
	client := NewClient(&conn, wait)
	publisher.Sub(client)
	<-wait
    publisher.Unsub(client)
    logger.Info(fmt.Sprintln(reverseName(conn),"gone"))
}

//monitorSignals trap SIGUSR1 to print stats
func monitorSignals(p *Publisher) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)
	for {
		<-c
        clients := make([]string,0)
        callback := func(client Subscriber){
            clients = append(clients,client.String())
        }
        p.dowithsubscribers(callback)
		logger.Info(fmt.Sprintln("Purges sent:", p.Publishes,". Connected Clients",clients))
	}
}

//version
func printVersion() {
	fmt.Println("0.0.2")
}
