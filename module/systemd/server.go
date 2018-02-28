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
	var err error
	ExecStart, err = os.Executable()
	if err != nil {
		panic(err.Error())
	}
	WorkingDirectory = filepath.Dir(ExecStart)
}

const (
	SERVICE_FILENAME       = "gohive.service"
	SERVICE_TEMPLATE       = "./templates/" + SERVICE_FILENAME
	SERVICE_GENERATED      = "./generated/" + SERVICE_FILENAME
	SERVICE_SYSTEMD_OUTPUT = "/lib/systemd/system/gohive.service"
)

func RegisterService() {

	fmt.Print("Generating systemd service file...")
	bytes, err := ioutil.ReadFile(SERVICE_TEMPLATE)
	if err != nil {
		panic(err.Error())
	}

	content := string(bytes)
	content = Replace(content, "{{ExecStart}}", ExecStart, -1)
	content = Replace(content, "{{WorkingDirectory}}", WorkingDirectory, -1)
	if err = ioutil.WriteFile(SERVICE_GENERATED, []byte(content), 0644); err != nil {
		panic(err.Error())
	}
	fmt.Println(" [ok]")

	fmt.Print("Creating symbol link for systemd...")
	if abs, err := filepath.Abs(SERVICE_GENERATED); err != nil {
		panic(err.Error())
	} else if err := os.Remove(SERVICE_SYSTEMD_OUTPUT); err != nil {
		// do nothing
		fmt.Println("Unable to remove outdated systemd services. [Ignored]")
	} else if err := os.Symlink(abs, SERVICE_SYSTEMD_OUTPUT); err != nil {
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
