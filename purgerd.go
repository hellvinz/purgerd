package main

import (
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
//when a client connects it is calling the handleConnection handler
func setupPurgeSenderAndListen(outgoingAddress *string, purgeOnStartup bool, publisher *Publisher, secret *string) {
	ln, err := net.Listen("tcp", *outgoingAddress)
	checkError(err)
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		logger.Info(fmt.Sprintln("New client: ", conn.RemoteAddr()))

		// connect client to the pubsub purge
		go connectClientToPublisher(conn, publisher, purgeOnStartup, secret)
	}
	return
}

//setupPurgeReceiver set up the tcp socket where ban messages come
//when a purge pattern is received it dispatches it to a Pub object
func setupPurgeReceiver(incomingAddress *string, publisher *Publisher) {
	receiver, err := net.Listen("tcp", *incomingAddress)
	checkError(err)

	go func() {
		for {
			time.Sleep(5 * time.Second)
			publisher.Pub([]byte("ping"))
		}
	}()
	for {
		conn, err := receiver.Accept()
		checkError(err)
		go func(c net.Conn) {
			b, err := ioutil.ReadAll(conn)
			if err != nil {
				conn.Close()
			}
			logger.Info(fmt.Sprintln("i've received to purge from client:", string(b)))
			publisher.Pub(b)
		}(conn)
	}
	return
}

//connectClientToPusher is used to forward message received from the internal PUB socket to the client
func connectClientToPublisher(conn net.Conn, publisher *Publisher, purgeOnStartup bool, secret *string) {
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
		sendPurge(conn, ".*")
	}

	subscriber := new(Subscriber)
	subscriber.Channel = make(chan []byte, 3)
	publisher.Sub(subscriber.Channel)
	defer publisher.Unsub(subscriber.Channel)
	for {
		b := <-subscriber.Channel
		var err error
		if string(b) == "ping" {
			err = sendString(conn, string(b))
		} else {
			err = sendPurge(conn, string(b))
		}
		if err == syscall.EPIPE {
			logger.Info(fmt.Sprintln("client gone", conn.RemoteAddr()))
			break
		} else {
			logger.Debug(fmt.Sprintln("Client got", conn.RemoteAddr(), string(b)))
		}
	}
}

//sendPurge send a purge message to a client
//it appends a ban.url to the pattern passed
func sendPurge(conn net.Conn, pattern string) (err error) {
	err = sendString(conn, "ban.url "+pattern)
	return
}

//sendString is sending a raw string to a client
func sendString(conn net.Conn, message string) (err error) {
	n, err := conn.Write([]byte(message + "\n"))
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

//monitorSignals trap SIGUSR1 to print stats
func monitorSignals(p *Publisher) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)
	for {
		<-c
		logger.Info(fmt.Sprintln("Purges sent: ", p.Publishes))
	}
}

//version
func printVersion() {
	fmt.Println("0.0.2")
}
