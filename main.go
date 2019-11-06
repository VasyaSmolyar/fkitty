package main

import (
	"fmt"
	"net"
	"flag"
	"strconv"
	"bufio"
	"strings"
)

type FtpConnect struct {
	conn net.Conn
	user string
	logged bool
	dir string
}

type FtpAnswer struct {
	code int
	status string
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
	ftp.logged = false
	ftp.user = "anonymous"
	for {
		line, err := b.ReadBytes('\n')
		if err != nil {
			break
		}
		ret := parser(string(line), ftp)
		conn.Write([]byte(ret))
	}
}

func parser(line string, ftp *FtpConnect) string {
	line = strings.Split(line, "\r")[0]
	line = strings.Split(line, "\n")[0]
	words := strings.Split(line, " ")
	tokens := make([]string, 0, len(words))
	for _, word := range words {
		if word != "" {
			tokens = append(tokens, word)
		}
	}
	res := new(FtpAnswer)
	execute(tokens, ftp, res)
	if res.code != 0 {
		return strconv.Itoa(res.code) + " " + res.status + "\n"
	}
	return ""
}
	
func main() { 
	port := flag.Int("port", 8000, "Number of port")
	flag.Parse() 
	go accept(*port)
	var test string
	fmt.Scanln(&test)
}