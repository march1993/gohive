package nginx

import (
	"fmt"
	"github.com/march1993/gohive/config"
	"github.com/march1993/gohive/module/systemd"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	. "strings"
)

const (
	nginxFilename  = "nginx.conf"
	nginxOutput    = "/etc/nginx/sites-enabled/default"
	nginxTemplate  = "./templates/" + nginxFilename
	nginxGenerated = "./generated/" + nginxFilename
)

func RegisterNginx() {

	fmt.Print("Generating nginx site file...")
	bytes, err := ioutil.ReadFile(nginxTemplate)
	if err != nil {
		panic(err.Error())
	}

	serverName := config.Get("server_name", "_")

	content := string(bytes)
	content = Replace(content, "{{Root}}", systemd.WorkingDirectory, -1)
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
