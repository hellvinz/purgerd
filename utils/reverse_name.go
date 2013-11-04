package utils

import "net"

var hostCache map[net.Conn]string

func init(){
    hostCache = map[net.Conn]string{}
}

//get reverse name
func ReverseName(conn net.Conn) (name string) {
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
