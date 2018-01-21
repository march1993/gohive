package main

import (
	"fmt"
	"github.com/march1993/gohive/admin"
	"github.com/march1993/gohive/module/nginx"
	"github.com/march1993/gohive/module/systemd"
	"github.com/march1993/gohive/util"
	"os/user"
)

func Install() {
	testRoot()
	systemd.RegisterService()
	nginx.RegisterNginx()
	updateAdminToken()
}

func testRoot() {
	fmt.Print("Testing if you are root ...")
	if user, err := user.Current(); err != nil {
		panic(err.Error())
	} else if (*user).Uid != "0" {
		panic("You should install gohive as root")
	}
	fmt.Println(" [ok]")
}

func updateAdminToken() {
	fmt.Println("Checking administration token... [ok]")
	if admin.GetToken() == "" {
		token := util.RandomString(32)
		admin.SetToken(token)
		fmt.Println("Generating new token: >>> " + token + " <<<")
		fmt.Println("Please use this token to login web panel.")
	} else {
		fmt.Println("Already set. Skipping.")
	}
}
