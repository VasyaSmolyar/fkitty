package main

func execute(tokens []string, ftp *FtpConnect, ans *FtpAnswer) {
	ans.code = 530
	ans.status = "You aren't logged in"
}