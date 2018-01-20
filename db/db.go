package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

const DB_NAME = "gohive"

func init() {

	var err error

	if DB, err = gorm.Open("mysql", "root@unix(/var/run/mysqld/mysqld.sock)/?charset=utf8"); err != nil {
		panic(err.Error())
	}

	DB.Exec("CREATE DATABASE IF NOT EXISTS " + DB_NAME)

	DB.Close()

	if DB, err = gorm.Open("mysql", "root@unix(/var/run/mysqld/mysqld.sock)/"+DB_NAME+"?charset=utf8"); err != nil {
		panic(err.Error())
	}

}
