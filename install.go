package main

import (
	"fmt"
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

func registerService() {

	var err error

	fmt.Print("Getting gohive executable directory ...")
	execStart, err := os.Executable()
	if err != nil {
		panic(err.Error())
	}

	workingDirectory := filepath.Dir(execStart)
	fmt.Println(" [ok]")

	fmt.Print("Generating systemd service file...")
	bytes, err := ioutil.ReadFile(serviceTemplate)
	if err != nil {
		panic(err.Error())
	}

	content := string(bytes)
	content = Replace(content, "{{ExecStart}}", execStart, -1)
	content = Replace(content, "{{WorkingDirectory}}", workingDirectory, -1)
	ioutil.WriteFile(serviceGenerated, []byte(content), 0644)
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

	fmt.Print("Enabling and starting service...")
	cmd := exec.Command("systemctl", "enable", "gohive.service")
	if stdout, err := cmd.CombinedOutput(); err != nil {
		panic(string(stdout) + err.Error())
	}

	cmd = exec.Command("systemctl", "start", "gohive.service")
	if stdout, err := cmd.CombinedOutput(); err != nil {
		panic(string(stdout) + err.Error())
	}
	fmt.Println(" [ok]")

}
