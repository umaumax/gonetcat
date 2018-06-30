package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"strings"
)

var (
	tlsFlag   bool
	httpFlag  bool
	debugFlag bool
)

func init() {
	flag.BoolVar(&tlsFlag, "tls", false, "tls flag")
	flag.BoolVar(&httpFlag, "http", false, "http or https flag")
	flag.BoolVar(&debugFlag, "debug", false, "debug log flag")
}

func main() {
	flag.Parse()

	n := flag.NArg()
	if n == 0 {
		log.Fatalln("Usage: netcat host:port")
	}
	addr := strings.Join(flag.Args(), ":")

	var conn net.Conn
	var err error
	if tlsFlag {
		conn, err = tls.Dial("tcp", addr, nil)
	} else {
		conn, err = net.Dial("tcp", addr)
	}
	if err != nil {
		log.Fatalln("dial error:", err)
	}

	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		log.Fatalln("Cannot convert %+v", conn)
	}

	if httpFlag {
		httpNetcat(tcpConn)
	} else {
		netcat(tcpConn)
	}
}
