package main

import (
	"github.com/march1993/gohive/module/linux"
	"github.com/march1993/gohive/util"
	"os"
	"os/exec"
)

func tryRestartService() {

	unixname := os.Getenv("SUDO_USER")
	list := linux.GetAppList()

	if util.Includes(list, unixname) {
		if output, err := exec.Command("systemctl", "restart", unixname).CombinedOutput(); err != nil {
			panic(string(output))
		}
	} else {
		panic("Invalid service name")
	}
}
