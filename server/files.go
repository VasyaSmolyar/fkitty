package main

import (
	"net"
	"strconv"
	"bufio"
	"io/ioutil"
	"path/filepath"
)

type FileConnect struct {
	conn net.Conn
	active bool
	write bool
	host string
	filename string
}

func createActive(ftp *FtpConnect, filename string, write bool) error {
	file := new(FileConnect)
	file.host = ftp.filehost
	file.active = true
	file.filename = ftp.dir + string(filepath.Separator) + filename
	file.write = write
	conn, err := net.Dial("tcp", ftp.filehost)
	if err != nil {
		return err
	}
	file.conn = conn
	ftp.file = *file
	go handleFile(ftp)
	return nil
}

func createPassive(ftp *FtpConnect, host string, filename string, write bool) error {
	file := new(FileConnect)
	file.host = host
	port := 1024
	file.active = false
	file.filename = filename
	file.write = write
	server, err := net.Listen("tcp", host + ":" + strconv.Itoa(port))
	for err != nil {
		port += 1
		if port > 9999 {
			return err
		}
		server, err = net.Listen("tcp", host + ":" + strconv.Itoa(port))
	}
	conn, err := server.Accept()
	if err != nil {
		return err
	}
	file.conn = conn
	ftp.file = *file
	go handleFile(ftp)
	return nil
}

func readAll(conn net.Conn) []byte {
	bs := make([]byte, 1024)
	res := make([]byte, 0)
	b := bufio.NewReader(conn)
	for {
		n, err := b.Read(bs)
		if err != nil {
			break
		}
		if n == 0 {
			break
		}
		res = append(res, bs[0:n-1]...)
	}
	return res 
}

func handleFile(ftp *FtpConnect) {
	ans := new(FtpAnswer)
	if ftp.file.write { 
		bs := readAll(ftp.file.conn)
		err := ioutil.WriteFile(ftp.file.filename, bs, 0644)
		if err != nil {
			ans.code = 526
			ans.status = err.Error()
		} else {
			ans.code = 226
			ans.status = "The file was successfully saved"
		}
	} else {
		bs, err := ioutil.ReadFile(ftp.file.filename)
		if err != nil {
			ans.code = 526
			ans.status = err.Error()
		} else {
			ftp.file.conn.Write(bs)
		}
	}
	ftp.write(ans)
}