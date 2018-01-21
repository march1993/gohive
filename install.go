package main

import (
	"fmt"
	"github.com/march1993/gohive/admin"
	"github.com/march1993/gohive/config"
	"github.com/march1993/gohive/util"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	. "strings"
)

func Install() {
	testRoot()
	registerService()
	registerNginx()
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

const (
	serviceFilename  = "gohive.service"
	serviceTemplate  = "./templates/" + serviceFilename
	serviceGenerated = "./generated/" + serviceFilename
	systemdOutput    = "/lib/systemd/system/gohive.service"
)

var workingDirectory string

func registerService() {

	var err error

	fmt.Print("Getting gohive executable directory ...")
	execStart, err := os.Executable()
	if err != nil {
		panic(err.Error())
	}

	workingDirectory = filepath.Dir(execStart)
	fmt.Println(" [ok]")

	fmt.Print("Generating systemd service file...")
	bytes, err := ioutil.ReadFile(serviceTemplate)
	if err != nil {
		panic(err.Error())
	}

	content := string(bytes)
	content = Replace(content, "{{ExecStart}}", execStart, -1)
	content = Replace(content, "{{WorkingDirectory}}", workingDirectory, -1)
	if err = ioutil.WriteFile(serviceGenerated, []byte(content), 0644); err != nil {
		panic(err.Error())
	}
	fmt.Println(" [ok]")

	fmt.Print("Creating symbol link for systemd...")
	if abs, err := filepath.Abs(serviceGenerated); err != nil {
		panic(err.Error())
	} else if err := os.Remove(systemdOutput); err != nil {
		panic(err.Error())
	} else if err := os.Symlink(abs, systemdOutput); err != nil {
		panic(err.Error())
	}
	fmt.Println(" [ok]")

	fmt.Print("Enabling and starting gohive service...")
	cmd := exec.Command("systemctl", "enable", "gohive")
	if stdout, err := cmd.CombinedOutput(); err != nil {
		panic(string(stdout) + err.Error())
	}

	cmd = exec.Command("systemctl", "restart", "gohive")
	if stdout, err := cmd.CombinedOutput(); err != nil {
		panic(string(stdout) + err.Error())
	}
	fmt.Println(" [ok]")

}

const (
	nginxFilename  = "nginx.conf"
	nginxOutput    = "/etc/nginx/sites-enabled/default"
	nginxTemplate  = "./templates/" + nginxFilename
	nginxGenerated = "./generated/" + nginxFilename
)

func registerNginx() {

	fmt.Print("Generating nginx site file...")
	bytes, err := ioutil.ReadFile(nginxTemplate)
	if err != nil {
		panic(err.Error())
	}

	serverName := config.Get("server_name", "_")

	content := string(bytes)
	content = Replace(content, "{{Root}}", workingDirectory, -1)
	content = Replace(content, "{{ServerName}}", serverName, -1)
	content = Replace(content, "{{APP_DIR}}", config.APP_DIR, -1)
	if err = ioutil.WriteFile(nginxGenerated, []byte(content), 0644); err != nil {
		panic(err.Error())
	}
	fmt.Println(" [ok]")

	fmt.Print("Creating symbol link for nginx...")
	if abs, err := filepath.Abs(nginxGenerated); err != nil {
		panic(err.Error())
	} else if err := os.Remove(nginxOutput); err != nil {
		panic(err.Error())
	} else if err := os.Symlink(abs, nginxOutput); err != nil {
		panic(err.Error())
	}
	fmt.Println(" [ok]")

	fmt.Print("Reloading nginx service...")
	cmd := exec.Command("systemctl", "reload", "nginx")
	if stdout, err := cmd.CombinedOutput(); err != nil {
		panic(string(stdout) + err.Error())
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
