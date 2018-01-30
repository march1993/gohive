package mariadb

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/config"
	. "github.com/march1993/gohive/db"
	"github.com/march1993/gohive/module"
	_ "github.com/march1993/gohive/module/linux"
	"github.com/march1993/gohive/util"
)

type mariadb struct{}

func init() {
	module.RegisterModule("mariadb", &mariadb{})
	module.RegisterRuncom(runcomHandler)
}

const (
	DB_HOST = "127.0.0.1"
	DB_PORT = "3306"
)

func runcomHandler(name string) []string {
	unixname := config.APP_PREFIX + name
	password := config.AppConfigGet(name, "mariadb", "db_password", "")

	return []string{
		"export DB_HOST=" + DB_HOST,
		"export DB_PORT=" + DB_PORT,
		"export DB_NAME=" + unixname,
		"export DB_USERNAME=" + unixname,
		"export DB_PASSWORD=" + password,
	}
}

func (m *mariadb) Create(name string) api.Status {
	unixname := config.APP_PREFIX + name
	password := util.RandomString(32)
	config.AppConfigSet(name, "mariadb", "db_password", password)

	DB.Exec("REVOKE ALL PRIVILEGES FROM '" + unixname + "'@'localhost';")
	DB.Exec("DROP USER '" + unixname + "'@'localhost';")
	DB.Exec("CREATE USER '" + unixname + "'@'localhost' IDENTIFIED BY '" + password + "';")
	DB.Exec("CREATE DATABASE " + unixname)
	DB.Exec("GRANT ALL PRIVILEGES ON " + unixname + " . * TO '" + unixname + "'@'localhost';")

	return api.Status{Status: api.STATUS_SUCCESS}
}

func (m *mariadb) Remove(name string) api.Status {
	unixname := config.APP_PREFIX + name

	DB.Exec("REVOKE ALL PRIVILEGES FROM '" + unixname + "'@'localhost';")
	DB.Exec("DROP DATABASE " + unixname)
	DB.Exec("DROP USER '" + unixname + "'@'localhost';")

	return api.Status{Status: api.STATUS_SUCCESS}
}

func (m *mariadb) Status(name string) api.Status {
	unixname := config.APP_PREFIX + name
	password := config.AppConfigGet(name, "mariadb", "db_password", "")

	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		unixname, password, DB_HOST, DB_PORT, unixname))

	if err != nil {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: err.Error(),
		}
	}
	db.Close()

	return api.Status{Status: api.STATUS_SUCCESS}
}

func (m *mariadb) Repair(name string) api.Status {
	return m.Create(name)
}

func (m *mariadb) ListRemoved() []string {
	return []string{}
}
