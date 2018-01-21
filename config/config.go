package config

import (
	. "github.com/march1993/gohive/db"
)

// Configurations for compilation time
const (
	LISTEN_HOST_PORT   = "127.0.0.1:1033"
	APP_DIR            = "/gohive"
	APP_DIR_PERM       = 0755
	APP_DIR_O_USER     = 0 // APP_DIR Owner
	APP_DIR_O_GROUP    = 0 // APP_DIR Group
	GOLANG_DIR         = "/gohove.go"
	GOLANG_DIR_PERM    = 0755
	GOLANG_DIR_O_USER  = 0 // GOLANG_DIR Owner
	GOLANG_DIR_O_GROUP = 0 // GOLANG_DIR Group
	SSH_SHELL          = "/usr/bin/git-shell"
)

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
