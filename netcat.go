package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"regexp"
	"sync"
)

// NOTE; httpのプロトコル通りContent-Lengthを検知してwaitする
func httpNetcat(conn *net.TCPConn) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		writer := bufio.NewWriter(os.Stdout)

		contentLength := 0
		data := make([]byte, 1024)
		headerString := ""
		bodyFlag := false
		fullLength := 0

		for {
			n, err := conn.Read(data)
			if debugFlag {
				log.Println(data[:n])
			}
			writer.Write(data[:n])
			writer.Flush()
			if n > 0 {
				fullLength += n
				if !bodyFlag {
					if preBodyIndex := bytes.Index(data[:n], []byte("\r\n\r\n")); preBodyIndex != -1 {
						headerString += string(data[:preBodyIndex])
						bodyFlag = true
						for _, v := range regexp.MustCompile("\r\n|\n\r|\n|\r").Split(headerString, -1) {
							if index := bytes.Index([]byte(v), []byte("Content-Length: ")); index != -1 {
								fmt.Sscanf(v, "Content-Length: %d", &contentLength)
							}
						}
					} else {
						headerString += string(data[:n])
					}
				}
				if len(headerString)+4+contentLength == fullLength {
					break
				}
			} else {
				log.Fatalf("n:%d\n", n)
				break
			}
			if err != nil {
				log.Fatalf("err:%s\n", err)
				break
			}
		}
		if debugFlag {
			log.Println("stdout write done...")
		}
		wg.Done()
	}()
	go func() {
		io.Copy(conn, os.Stdin)
		if debugFlag {
			log.Println("stdin read done...")
		}
		wg.Done()
	}()
	wg.Wait()
	conn.Close()
}

func netcat(conn *net.TCPConn) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		//	receive
		io.Copy(os.Stdout, conn)
		conn.CloseRead()
		if debugFlag {
			log.Println("stdout write done...")
		}
		wg.Done()
	}()
	go func() {
		//	send
		io.Copy(conn, os.Stdin)
		conn.CloseWrite()
		if debugFlag {
			log.Println("stdin read done...")
		}
		wg.Done()
	}()
	wg.Wait()
}
