package main

import (
	"fmt"
	"log/syslog"
	"net"
	"os"
)

var hostCache = map[net.Conn]string{}

//checkError basic error handling
func checkError(err error, logger *syslog.Writer) {
	if err != nil {
		logger.Crit(fmt.Sprintln("Fatal error", err.Error()))
		os.Exit(1)
	}
}

//get reverse name
func reverseName(conn net.Conn) (name string) {
	name = hostCache[conn]
	if name == "" {
		ip, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
		names, err := net.LookupAddr(ip)
		if err == nil {
			name = names[0]
		} else {
			name = ip
		}
		hostCache[conn] = name
	}
	return name
}
