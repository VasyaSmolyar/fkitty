package main

import (
	"strings"
	"crypto/sha512"
	"encoding/hex"
	"os"
	"path/filepath"
)

func execute(tokens []string, ftp *FtpConnect, ans *FtpAnswer) {
	// create and add your methods in here
	methods := map[string]func(args []string, ftp *FtpConnect, ans *FtpAnswer){
		"USER" : user,
		"PASS" : pass,
		"QUIT" : quit,
		"MKD" : isLogged(mkd),
		"RMD" : isLogged(rmd),
		"CWD" : isLogged(cwd),
		"PWD" : isLogged(pwd),
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

func isLogged(ex func(args []string, ftp *FtpConnect, ans *FtpAnswer)) func(args []string, ftp *FtpConnect, ans *FtpAnswer) {
	return func(args []string, ftp *FtpConnect, ans *FtpAnswer) {
		if !ftp.logged {
			ans.code = 530
			ans.status = "You aren't logged in"
			return
		}
		ex(args, ftp, ans)
	}
}

func user(args []string, ftp *FtpConnect, ans *FtpAnswer) {
	if ftp.logged == true {
		ans.code = 530 
		ans.status = "You're already logged in"
		return
	}
	if len(args) == 0 {
		ans.code = 530
		ans.status = "This is a private system - No anonymous login"
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
	if len(args) == 0 {
		ans.code = 530
		ans.status = "Login authentication failed"
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

func quit(args []string, ftp *FtpConnect, ans *FtpAnswer) {
	ftp.connected = false
	ans.code = 221
	ans.status = "Logout."
}

func mkd(args []string, ftp *FtpConnect, ans *FtpAnswer) {
	dir := args[0]
	err := os.Mkdir(ftp.dir + string(filepath.Separator) + args[0], os.ModePerm)
	if err != nil {
		ans.code = 527
		ans.status = err.Error()
		return
	}
	ans.code = 257
	ans.status = string('"') + dir + string('"') + " : The directory was successfully created"
}

func rmd(args []string, ftp *FtpConnect, ans *FtpAnswer) {
	dir := args[0]
	err := os.Remove(ftp.dir + string(filepath.Separator) + dir)
	if err != nil {
		ans.code = 550
		ans.status = err.Error()
		return
	}
	ans.code = 250
	ans.status = string('"') + dir + string('"') + " : The directory was successfully removed"
}

func cwd(args []string, ftp *FtpConnect, ans *FtpAnswer) {
	dir := args[0]
	if _, err := os.Stat(ftp.dir + string(filepath.Separator) + args[0]); os.IsNotExist(err) {
		ans.code = 550
		ans.status = err.Error()
		return
	}
	ftp.dir = ftp.dir + string(filepath.Separator) + dir
	ans.code = 250
	ans.status = "OK. Current directory is " + ftp.dir
}

func pwd(args []string, ftp *FtpConnect, ans *FtpAnswer) {
	ans.code = 257
	ans.status = string('"') + ftp.dir + string('"') + " : is your current location"
}