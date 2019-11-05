package main

import (
	"fmt"
	"net"
	"flag"
	"strconv"
	"bufio"
)

type FtpConnect struct {
	conn net.Conn
	user string
	dir string
}

func accept(port int) {
	server, err := net.Listen("tcp", ":" + strconv.Itoa(port))
	if err != nil {
		fmt.Println("listen:", err)
		return
	}
	ch := make(chan net.Conn)
	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				fmt.Println("listen:", err)
				break
			}
			ch <- conn
		}
	}()
	go func() {
		for {
			conn := <-ch
			go handle(conn)
		}
	}()
}

func handle(conn net.Conn) {
	b := bufio.NewReader(conn)
	ftp := new(FtpConnect)
	ftp.conn = conn
	ftp.user = "anonymous"
	for {
		line, err := b.ReadBytes('\n')
		if err != nil {
			fmt.Println("close conn")
			break
		}
		ret := parser(string(line), ftp)
		conn.Write([]byte(ret))
	}
}

func parser(line string, ftp *FtpConnect) string {
	if ftp.user == "anonymous" {
		ftp.user = "User: " + line
	}
	return ftp.user
}
	
func main() { 
	port := flag.Int("port", 8000, "Number of port")
	flag.Parse() 
	go accept(*port)
	var test string
	fmt.Scanln(&test)
}