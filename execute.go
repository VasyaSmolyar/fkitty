package main

import (
	"strings"
	"crypto/sha512"
	"encoding/hex"
)

func execute(tokens []string, ftp *FtpConnect, ans *FtpAnswer) {
	// create and add your methods in here
	methods := map[string]func(args []string, ftp *FtpConnect, ans *FtpAnswer){
		"USER" : user,
		"PASS" : pass,
	}

	if len(tokens) == 0 {
		ans.code = 0
		ans.status = ""
		return
	}
	if method, ok := methods[strings.ToUpper(tokens[0])]; ok {
		method(tokens[1:], ftp, ans)
	} else {
		ans.code = 500
		ans.status = "Unknown command"
	}
}

func user(args []string, ftp *FtpConnect, ans *FtpAnswer) {
	ftp.user = args[0]
	ftp.logged = false
	ans.code = 331
	ans.status = "User " + ftp.user + " OK. Password required"
}

func pass(args []string, ftp *FtpConnect, ans *FtpAnswer) {
	hash := sha512.New();
	hash.Write([]byte(args[0]))
	if hex.EncodeToString(hash.Sum(nil)) == passes[ftp.user] {
		ftp.logged = true
		ans.code = 230
		// TODO: add a getting directory name
		ans.status = "OK. Current directory is"
	} else {
		ans.code = 530
		ans.status = "Login authentication failed"
	}

}