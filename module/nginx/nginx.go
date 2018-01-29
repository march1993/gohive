package nginx

import (
	"encoding/json"
	"github.com/march1993/gohive/config"
	"github.com/march1993/gohive/module"
	"github.com/march1993/gohive/module/linux"
	"github.com/march1993/gohive/util"
	"strconv"
)

func init() {
	// module.RegisterModule("profile", &profile{})
	module.RegisterRuncom(runcomHandler)
}

func runcomHandler(name string) []string {

	host := "127.0.0.1"
	port := strconv.FormatUint(uint64(getPort(name)), 10)

	return []string{
		"export HOST=" + host,
		"export PORT=" + port,
		"export HOST_PORT=" + host + ":" + port,
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
