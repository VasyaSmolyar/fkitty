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

func getUserData(user string) map[string]string {
	for _, val := range users {
		if val["login"] == user {
			return val
		}
	}
	return nil
}

func user(args []string, ftp *FtpConnect, ans *FtpAnswer) {
	if ftp.logged == true {
		ans.code = 530 
		ans.status = "You're already logged in"
		return
	}
	ftp.user = args[0]
	ans.code = 331
	ans.status = "User " + ftp.user + " OK. Password required"
}

func pass(args []string, ftp *FtpConnect, ans *FtpAnswer) {
	if ftp.logged == true {
		ans.code = 530 
		ans.status = "We can't do that in the current session"
		return
	}
	hash := sha512.New();
	hash.Write([]byte(args[0]))
	if hex.EncodeToString(hash.Sum(nil)) == passes[ftp.user] {
		data := getUserData(ftp.user)
		ftp.logged = true
		ftp.dir = data["dir"]
		ans.code = 230
		ans.status = "OK. Current directory is " + data["dir"]
	} else {
		ans.code = 530
		ans.status = "Login authentication failed"
	}

}