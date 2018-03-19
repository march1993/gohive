package config

import (
	"encoding/json"
	. "github.com/march1993/gohive/db"
	"io/ioutil"
)

// Configurations for compilation time
const (
	LISTEN_HOST_PORT   = "127.0.0.1:1033"
	APP_DIR            = "/gohive"
	APP_DIR_PERM       = 0755
	APP_DIR_O_USER     = 0 // APP_DIR Owner
	APP_DIR_O_GROUP    = 0 // APP_DIR Group
	GOLANG_DIR         = "/gohive.go"
	GOLANG_DIR_PERM    = 0755
	GOLANG_DIR_O_USER  = 0           // GOLANG_DIR Owner
	GOLANG_DIR_O_GROUP = 0           // GOLANG_DIR Group
	SSH_SHELL          = "/bin/bash" // "/usr/bin/git-shell"

	APP_N_LIMIT uint16 = 1024 // support maximum of 1Ki applications
)

// Linux setting
const (
	APP_PREFIX                  = "gohive_app_"
	APP_GROUP                   = "gohive_app"
	APP_DATA_SUFFIX             = ".data"
	APP_PORT_BEGIN       uint16 = 2000
	APP_RESTART_INTERVAL        = "60s" // restart interval when the service fails
	APP_WORKSPACE               = "workspace"
)

// Golang setting
const (
	GOLANG_GOPATH     = "$HOME/workspace:$HOME/workspace/vendor"
	GOLANG_EXECUTABLE = "output.exe"
)

/**
 *	Directories
 */

func GetHomeDir(name string) string {
	return APP_DIR + "/" + APP_PREFIX + name
}

func GetDataDir(name string) string {
	return APP_DIR + "/" + APP_PREFIX + name + APP_DATA_SUFFIX
}

/**
 *	Configurations stored in database
 */

type Config struct {
	Key   string
	Value string
}

func init() {
	DB.AutoMigrate(&Config{})
}

func Get(key string, voreinstellung string) string {

	config := Config{}
	DB.Where(&Config{Key: key}).First(&config)
	if config.Value != "" {
		return config.Value
	} else {
		return voreinstellung
	}

}

func Set(key string, value string) {
	DB.Where(Config{Key: key}).Assign(Config{Value: value}).FirstOrCreate(&Config{})
}

func getConfigFile(name, module string) string {
	return GetDataDir(name) + "/" + module + ".conf.json"
}

type ConfigList map[string]string

func AppConfigGet(name, module, key, voreinstellung string) string {
	config := getConfigFile(name, module)

	bytes, err := ioutil.ReadFile(config)
	if err != nil {
		return voreinstellung
	}

	configList := ConfigList{}

	if err := json.Unmarshal(bytes, &configList); err != nil {
		return voreinstellung
	}

	if configList[key] == "" {
		return voreinstellung
	} else {
		return configList[key]
	}
}

func AppConfigSet(name, module, key, value string) {
	config := getConfigFile(name, module)

	configList := ConfigList{}

	bytes, err := ioutil.ReadFile(config)
	if err == nil {
		json.Unmarshal(bytes, &configList)
	}

	configList[key] = value

	bytes, err = json.Marshal(configList)
	if err != nil {
		panic(err.Error())
	}
	ioutil.WriteFile(config, bytes, 0644)

}
