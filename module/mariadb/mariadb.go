package mariadb

import (
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/module"
	_ "github.com/march1993/gohive/module/linux"
)

type mariadb struct{}

func init() {
	module.RegisterModule("mariadb", &mariadb{})
}

func (m *mariadb) Create(name string) api.Status {
	return api.Status{Status: api.STATUS_SUCCESS}
}

func (m *mariadb) Rename(oldName, newName string) api.Status {
	return api.Status{Status: api.STATUS_SUCCESS}
}

func (m *mariadb) Remove(name string) api.Status {
	return api.Status{Status: api.STATUS_SUCCESS}
}

func (m *mariadb) Status(name string) api.Status {
	return api.Status{Status: api.STATUS_SUCCESS}
}

func (m *mariadb) Repair(name string) api.Status {
	return api.Status{Status: api.STATUS_SUCCESS}
}

func (m *mariadb) ListRemoved() []string {
	return []string{}
}
