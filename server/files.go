package main

import (
	"net"
	"strconv"
	"bufio"
	"io/ioutil"
)

type FileConnect struct {
	conn net.Conn
	active bool
	write bool
	host string
	filename string
}

func createActive(ftp *FtpConnect, host string, filename string, write bool) error {
	file := new(FileConnect)
	file.host = host
	file.active = true
	file.filename = filename
	file.write = write
	conn, err := net.Dial("tcp", host)
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

func handleFile(ftp *FtpConnect) {
	ans := new(FtpAnswer)
	if ftp.file.write { 
		var bs []byte
		b := bufio.NewReader(ftp.file.conn)
		b.Read(bs)
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