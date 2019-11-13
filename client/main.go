package main

import (
    "fmt"
    "net"
	"os"
	"flag"
	"bufio"
	"strings"
	"io/ioutil"
	"errors"
	"strconv"
)

var filehost string
var hosterror error
var fch chan net.Conn

func alloc() map[string]func (com string, conn net.Conn) bool {
	hosterror = nil
	return map[string]func (com string, conn net.Conn) bool {
		"QUIT" : quit,
		"PORT" : port,
		"STOR" : store,
	}
}

func getHost(host string) (string, error) {
	args := strings.Split(host, ",")
	if len(args) != 6 {
		return "", errors.New("Parsing error")
	}
	ips := make([]int, len(args))
	for i, v := range args {
		nv, err := strconv.Atoi(v)
		if err != nil {
			return "", err
		}
		ips[i] = nv
	}
	hex1 := strconv.FormatInt(int64(ips[4]), 16)
	if len(hex1) == 1 {
		hex1 = "0" + hex1 
	}
	hex2 := strconv.FormatInt(int64(ips[5]), 16)
	if len(hex2) == 1 {
		hex2 = "0" + hex2
	}
	hex := hex1 + hex2
	port, _ := strconv.ParseInt("0x" + hex, 0, 64)
	return args[0] + "." + args[1] + "." + args[2] + "." + args[3] + ":" + strconv.Itoa(int(port)), nil
}

func createActive(host string, ch chan net.Conn)  {
	server, err := net.Listen("tcp", host)
	if err != nil {
		fmt.Println("listen:", err)
		return
	}
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("listen:", err)
			break
		}
		ch <- conn
	}
}

func sendActive(ch chan net.Conn, data []byte) {
	conn := <- ch
	conn.Write(data)
	conn.Close()
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

func read(conn net.Conn) {
	b := bufio.NewReader(conn)
	res, _ := b.ReadBytes('\n')
	fmt.Print(string(res) + "fkitty: ")
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
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		com = scanner.Text()
		if com == "" {
			continue
		}
		res := send(conn, []byte(com))
		fmt.Print(res)
		for k,v := range methods {
			if strings.HasPrefix(strings.ToUpper(com), k) {
				work = v(com, conn)
			}
		}
	}
}

func quit(com string, conn net.Conn) bool {
	return false;
}

func port(com string, conn net.Conn) bool {
	names := strings.Split(com, " ")
	if len(names) < 2 {
		fmt.Println("fkitty PORT: Invalid arguments")
		return true
	}
	host, err := getHost(names[1])
	if err != nil {
		hosterror = err
	}
	filehost = host
	fch = make(chan net.Conn)
	go createActive(filehost, fch)
	return true;
}

func store(com string, conn net.Conn) bool {
	names := strings.Split(com, " ")
	if len(names) < 2 {
		fmt.Println("fkitty STORE: Invalid arguments")
		return true
	}
	if hosterror != nil {
		fmt.Println("fkitty STORE: " + hosterror.Error())
		hosterror = nil
		return true
	}
	var fn string
	fmt.Print("Enter a file name: ")
	fmt.Scanln(&fn)
	bs, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Print(err)
	} else {
		go sendActive(fch, bs)
	}
	go read(conn)
	return true;
}