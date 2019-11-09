package main

import (
    "fmt"
    "net"
	"os"
	"flag"
	"bufio"
	"strings"
)

func alloc() map[string]func (com string) bool {
	return map[string]func (com string) bool {
		"QUIT" : quit,
	}
}

func send(conn net.Conn, data []byte) string {
	if data != nil {
		conn.Write(data)
		conn.Write([]byte("\n"))
	}
	b := bufio.NewReader(conn)
	res, _ := b.ReadBytes('\n')
	return string(res)
}

func login(conn net.Conn, name string) bool {
	if name == "" {
		fmt.Print("login as: ")
		fmt.Scanln(&name)
	}
	fmt.Print(send(conn, []byte("USER " + name)))
	var pass string
	fmt.Print("Password: ")
	fmt.Scanln(&pass)
	res := send(conn, []byte("PASS " + pass))
	fmt.Print(res)
	return strings.HasPrefix(res, "230")
}

func main() {
	args := os.Args
	var name string 
	flag.StringVar(&name ,"u", "", "Username")
	if len(args) == 1 {
		fmt.Println("Fkitty Client: You must specify a host:port to connect to.")
		return
	}
	flag.Parse() 
	host := args[len(args)-1]
	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(send(conn, nil))
	log := login(conn, name)
	for !log {
		log = login(conn, "")
	}
	var com string
	work := true
	methods := alloc()
	for work {
		fmt.Print("fkitty: ")
		l, _ := fmt.Scanln(&com)
		if l == 0 {
			continue
		}
		res := send(conn, []byte(com))
		fmt.Print(res)
		for k,v := range methods {
			if strings.HasPrefix(strings.ToUpper(com), k) {
				work = v(com)
			}
		}
	}
}

func quit(com string) bool {
	return false;
}