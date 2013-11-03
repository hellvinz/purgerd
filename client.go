package main

import (
	"net"
	"syscall"
)

type Client struct {
	messagesChannel chan []byte
	conn            *net.Conn
	workInProgess   chan bool
	//exitFunc func(*Client)
}

func NewClient(conn *net.Conn, workInProgress chan bool) *Client {
	client := new(Client)
	client.messagesChannel = make(chan []byte, 10)
	client.conn = conn
	client.workInProgess = workInProgress
	go client.monitorMessages()
	return client
}

func (c *Client) Receive(message []byte) {
	c.messagesChannel <- message
}

func (c *Client) monitorMessages() {
	defer close(c.messagesChannel)
	for {
        message := <-c.messagesChannel
        err := c.sendMessage(message)
        if err == syscall.EPIPE {
            break
        }
    }
    c.exit()
}

func (c *Client) sendMessage(message []byte) (err error) {
	if string(message) == "ping" {
		err = c.sendString(message)
	} else {
		err = c.sendPurge(message)
	}
	return
}

func (c *Client) exit() {
	c.workInProgess <- false
}

//sendPurge send a purge message to a client
//it appends a ban.url to the pattern passed
func (c *Client) sendPurge(pattern []byte) (err error) {
	err = c.sendString(append([]byte("ban.url "), pattern...))
	return
}

//sendString is sending a raw string to a client
func (c *Client) sendString(message []byte) (err error) {
	n, err := (*c.conn).Write(append(message, []byte("\n")...))
	if n == 0 {
		err = syscall.EPIPE
	}
	return
}

func (c *Client) String() string {
    return reverseName(*c.conn)
}
