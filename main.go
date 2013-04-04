package main

import "net"
import "fmt"
import "log"
import "bufio"
import "os"
import "strconv"
import "strings"

// 0: new
// 1: open
func handle(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	dataRead := 0
	
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
		dataRead = 0
		for dataRead < dataLen {
			n, err := reader.Read(dataBytes[dataRead:])
			if err != nil {
				log.Println(err)
				return
			}
			dataRead += n
		}
		
/*
18:01:09.603945 IP 10.31.156.23.52077 > 23.23.116.103.9996: Flags [P.], seq 1:97, ack 1, win 12, options [nop,nop,TS val 176021 ecr 138772348], length 96
	0x0030:  0845 7f7c 3120 6f70 656e 2038 3520 7265  .E.|1.open.85.re
	0x0040:  6c70 5f76 6572 7369 6f6e 3d30 0a72 656c  lp_version=0.rel
	0x0050:  705f 736f 6674 7761 7265 3d6c 6962 7265  p_software=libre
	0x0060:  6c70 2c31 2e30 2e30 2c68 7474 703a 2f2f  lp,1.0.0,http://
	0x0070:  6c69 6272 656c 702e 6164 6973 636f 6e2e  librelp.adiscon.
	0x0080:  636f 6d0a 636f 6d6d 616e 6473 3d73 7973  com.commands=sys
	0x0090:  6c6f 670a                                log.

18:01:09.605465 IP 23.23.116.103.9996 > 10.31.156.23.52077: Flags [P.], seq 1:103, ack 97, win 12, options [nop,nop,TS val 138772348 ecr 176021], length 102
	0x0030:  0002 af95 3120 7273 7020 3932 2032 3030  ....1.rsp.92.200
	0x0040:  204f 4b0a 7265 6c70 5f76 6572 7369 6f6e  .OK.relp_version
	0x0050:  3d30 0a72 656c 705f 736f 6674 7761 7265  =0.relp_software
	0x0060:  3d6c 6962 7265 6c70 2c31 2e30 2e30 2c68  =librelp,1.0.0,h
	0x0070:  7474 703a 2f2f 6c69 6272 656c 702e 6164  ttp://librelp.ad
	0x0080:  6973 636f 6e2e 636f 6d0a 636f 6d6d 616e  iscon.com.comman
	0x0090:  6473 3d73 7973 6c6f 670a                 ds=syslog.

18:06:45.418428 IP 127.0.0.1.10000 > 127.0.0.1.36232: Flags [P.], seq 1:107, ack 97, win 64, options [nop,nop,TS val 209602 ecr 209592], length 106
	0x0030:  0003 32b8 3120 7273 7020 3938 2032 3030  ..2.1.rsp.98.200
	0x0040:  204f 4b0a 7265 6c70 5f76 6572 7369 6f6e  .OK.relp_version
	0x0050:  3d30 0a72 656c 705f 736f 6674 7761 7265  =0.relp_software
	0x0060:  3d67 6f72 656c 702c 302e 302e 312c 6874  =gorelp,0.0.1,ht
	0x0070:  7470 3a2f 2f67 6974 6875 622e 636f 6d2f  tp://github.com/
	0x0080:  6470 6964 6479 2f67 6f72 656c 700a 636f  dpiddy/gorelp.co
	0x0090:  6d6d 616e 6473 3d73 7973 6c6f 670a       mmands=syslog.
*/

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
