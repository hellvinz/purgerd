package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

type clientPool []net.Conn

var clients clientPool

func main() {
	go sendPurges()
	clients = make(clientPool, 0)
	defer clients.close()
	ln, err := net.Listen("tcp", ":8080")
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

func handleConnection(conn net.Conn) {
	fmt.Println("New client", conn.RemoteAddr())
    credentials, _ := ioutil.ReadAll(conn)
	if strings.Contains(string(credentials), "200 194") {
        fmt.Println("credentials:",credentials)
		// flush the whole cache
		sendPurge(conn, ".*")
		// put it in the client pool
		clients = append(clients, conn)
	} else {
        conn.Write([]byte("no you can't\n"))
		conn.Close()
	}
}

func (clients clientPool) dispatchPurge(pattern string) {
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

func sendPurges() {
	fmt.Println("going to purge in 20 secondes")
	time.Sleep(20 * time.Second)
	fmt.Println("ok time to purge")
	clients.dispatchPurge(".*")
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
