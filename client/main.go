package main

import (
    "fmt"
    "net"
	"os"
	"flag"
	"bufio"
)

func send(conn net.Conn, data []byte) string {
	if data != nil {
		conn.Write(data)
		conn.Write([]byte("\n"))
	}
	b := bufio.NewReader(conn)
	res, _ := b.ReadBytes('\n')
	return string(res)
}

func login(conn net.Conn, name string) {
	if name == "" {
		fmt.Print("login as: ")
		fmt.Scanln(&name)
	}
	fmt.Print(send(conn, []byte("USER " + name)))
	var pass string
	fmt.Print("Password: ")
	fmt.Scanln(&pass)
	fmt.Print(send(conn, []byte("PASS " + pass)))
}

func main() {
	args := os.Args
	name := flag.String("u", "", "Username")
	if len(args) == 1 {
		fmt.Println("Fkitty Client: You must specify a host:port to connect to.")
		return
	}
	flag.Parse() 
	host := os.Args[1]
	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(send(conn, nil))
	login(conn, *name)
}