package config

import (
	. "github.com/march1993/gohive/db"
)

type Config struct {
	Key   string `gorm:"primary_key"`
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
