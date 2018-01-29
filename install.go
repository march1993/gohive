package main

import (
	"fmt"
	"github.com/march1993/gohive/admin"
	"github.com/march1993/gohive/config"
	"github.com/march1993/gohive/module/nginx"
	"github.com/march1993/gohive/module/systemd"
	"github.com/march1993/gohive/util"
	"io/ioutil"
	"os/user"
	"strings"
)

func Install() {
	testRoot()
	systemd.RegisterService()
	nginx.RegisterNginx()
	updateAdminToken()
	updateSudoer()
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

const (
	SUDOER_TEMPLATE  = "./templates/sudoer"
	SUDOER_GENERATED = "/etc/sudoers.d/gohive"
)

func updateSudoer() {
	fmt.Print("Updating sudoers policy ...")
	bytes, err := ioutil.ReadFile(SUDOER_TEMPLATE)
	if err != nil {
		panic(err.Error())
	}

	content := string(bytes)
	content = strings.Replace(content, "{{APP_GROUP}}", config.APP_GROUP, -1)
	content = strings.Replace(content, "{{GOHIVE}}", systemd.ExecStart, -1)

	err = ioutil.WriteFile(SUDOER_GENERATED, []byte(content), 0440)

	if err != nil {
		panic(err.Error())
	}
	fmt.Println(" [ok]")
}
