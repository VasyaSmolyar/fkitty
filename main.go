package main

import (
	"fmt"
	"net"
	"flag"
	"strconv"
	"bufio"
	"strings"
	"encoding/json"
	"io/ioutil"
)

type FtpConnect struct {
	conn net.Conn
	connected bool
	user string
	logged bool
	dir string
}

type FtpAnswer struct {
	code int
	status string
}

var params map[string]string
var users []map[string]string
var passes map[string]string

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
	ftp.connected = true
	ftp.user = "anonymous"
	for ftp.connected {
		line, err := b.ReadBytes('\n')
		if err != nil {
			break
		}
		ret := parser(string(line), ftp)
		conn.Write([]byte(ret))
	}
	conn.Close()
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

func loadParams() bool {
	var config map[string]interface{}
	bs, err := ioutil.ReadFile(params["config"])
	if err != nil {
		fmt.Println("config: ", err)
		return false
	}
	json.Unmarshal([]byte(bs), &config)
	pathes := config["pathes"].(map[string]interface{})
	for k,_ := range pathes {
		params[k] = pathes[k].(string)
	}
	us, err := ioutil.ReadFile(params["users"])
	if err != nil {
		fmt.Println("users: ", err)
		return false
	}
	json.Unmarshal([]byte(us), &users)
	ps, err := ioutil.ReadFile(params["shadow"])
	if err != nil {
		fmt.Println("passes: ", err)
		return false
	}
	passes = make(map[string]string)
	for _,v := range strings.Split(string(ps),"\n") {
		match := strings.Split(v,":")
		passes[match[0]] = match[1]
	}
	return true
}
	
func main() {
	params = make(map[string]string) 
	port := flag.Int("port", 8000, "Number of port")
	config := flag.String("users", "data/config.json", "Path to configs file")
	params["config"] = *config
	flag.Parse() 
	if !loadParams() {
		return
	}
	go accept(*port)
	var test string
	fmt.Scanln(&test)
}