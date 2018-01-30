package nginx

import (
	"encoding/json"
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/config"
	"github.com/march1993/gohive/module"
	"github.com/march1993/gohive/module/linux"
	"github.com/march1993/gohive/util"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
)

type nginx struct{}

func init() {
	module.RegisterModule("nginx", &nginx{})
	module.RegisterRuncom(runcomHandler)
}

const (
	HOST = "127.0.0.1"
)

func runcomHandler(name string) []string {

	port := strconv.FormatUint(uint64(getPort(name)), 10)

	return []string{
		"export HOST=" + HOST,
		"export PORT=" + port,
		"export HOST_PORT=" + HOST + ":" + port,
	}
}

type PortAssign map[uint16]string

const CONFIG_KEY = "nginx-port-assign"

func getPort(name string) uint16 {
	list := linux.GetAppList()
	assign := PortAssign{}
	json.Unmarshal([]byte(config.Get(CONFIG_KEY, "{}")), assign)

	for i := config.APP_PORT_BEGIN; i < config.APP_PORT_BEGIN+config.APP_N_LIMIT; i++ {

		if !util.Includes(list, assign[i]) {
			assign[i] = name
			bytes, _ := json.Marshal(assign)
			config.Set(CONFIG_KEY, string(bytes))
			return i
		}

	}

	panic("APP_N_LIMIT: No more available port to assign")
}

const (
	NGINX_CONF_TEMPLATE = "./templates/app.nginx.conf"
)

func (n *nginx) Create(name string) api.Status {

	bytes, err := ioutil.ReadFile(NGINX_CONF_TEMPLATE)
	if err != nil {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: err.Error(),
		}
	}

	ServerName := name + "." + config.Get("server_name", "_")
	ProxyPass := "http://" + HOST + ":" + strconv.FormatUint(uint64(getPort(name)), 10)

	content := string(bytes)
	content = strings.Replace(content, "{{ServerName}}", ServerName, -1)
	content = strings.Replace(content, "{{ProxyPass}}", ProxyPass, -1)

	config.AppConfigSet(name, "nginx", "hash", util.Hash(content+config.Get("server_name", "_")))
	output := config.GetDataDir(name) + "/nginx.conf"
	err = ioutil.WriteFile(output, []byte(content), 0644)
	if err != nil {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: err.Error(),
		}
	}

	return nginxReload()
}

func (n *nginx) Remove(name string) api.Status {
	return nginxReload()
}

func (n *nginx) Status(name string) api.Status {

	hash := config.AppConfigGet(name, "nginx", "hash", "")

	output := config.GetDataDir(name) + "/nginx.conf"
	bytes, err := ioutil.ReadFile(output)
	if err != nil {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: err.Error(),
		}
	}

	if hash != util.Hash(string(bytes)+config.Get("server_name", "_")) || "" == hash {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: api.NGINX_CONF_EXPIRED,
		}
	}

	return api.Status{Status: api.STATUS_SUCCESS}
}

func (n *nginx) Repair(name string) api.Status {
	return n.Create(name)
}

func (n *nginx) ListRemoved() []string {
	return []string{}
}

func nginxReload() api.Status {
	output, err := exec.Command("systemctl", "reload", "nginx").CombinedOutput()
	if err != nil {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: string(output),
		}
	} else {

		return api.Status{Status: api.STATUS_SUCCESS}
	}
}
