package systemd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	. "strings"
)

var WorkingDirectory, ExecStart string

func init() {
	ExecStart, err := os.Executable()
	if err != nil {
		panic(err.Error())
	}
	WorkingDirectory = filepath.Dir(ExecStart)
}

const (
	serviceFilename  = "gohive.service"
	serviceTemplate  = "./templates/" + serviceFilename
	serviceGenerated = "./generated/" + serviceFilename
	systemdOutput    = "/lib/systemd/system/gohive.service"
)

func RegisterService() {

	fmt.Print("Generating systemd service file...")
	bytes, err := ioutil.ReadFile(serviceTemplate)
	if err != nil {
		panic(err.Error())
	}

	content := string(bytes)
	content = Replace(content, "{{ExecStart}}", ExecStart, -1)
	content = Replace(content, "{{WorkingDirectory}}", WorkingDirectory, -1)
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
