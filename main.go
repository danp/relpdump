package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

// 0: new
// 1: open
func handle(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	state := 0

	for {
		txn, err := reader.ReadString(' ')
		if err != nil {
			log.Println(err)
			return
		}
		txn = strings.TrimSpace(txn)

		cmd, err := reader.ReadString(' ')
		if err != nil {
			log.Println(err)
			return
		}
		cmd = strings.TrimSpace(cmd)

		// TODO: handle 0 datalen -- loop on bytes until non-digit?
		dataLenString, err := reader.ReadString(' ')
		if err != nil {
			log.Println(err)
			return
		}
		dataLen, err := strconv.Atoi(strings.TrimSpace(dataLenString))
		if err != nil {
			log.Println(err)
			return
		}

		dataBytes := make([]byte, dataLen)
		_, err = io.ReadFull(reader, dataBytes)
		if err != nil {
			log.Println(err)
			return
		}

		switch cmd {
		case "open":
			_, err := conn.Write([]byte(fmt.Sprintf("%s rsp 92 200 OK\nrelp_version=0\nrelp_software=librelp,1.0.0,http://librelp.adiscon.com\ncommands=syslog\n", txn)))
			if err != nil {
				log.Println(err)
				return
			}
			state = 1
		default:
			fmt.Printf("%s\n", dataBytes)
			if state != 1 {
				_, err := conn.Write([]byte(fmt.Sprintf("%s rsp 7 500 ERR\n", txn)))
				if err != nil {
					log.Println(err)
				}
				return
			} else {
				_, err := conn.Write([]byte(fmt.Sprintf("%s rsp 6 200 OK\n", txn)))
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handle(conn)
	}
}
