package linux

import (
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/module"
)

type linux struct{}

func init() {
	module.RegisterModule("linux", &linux{})
}

func (l *linux) Create(name string) error {
	return nil
}

func (l *linux) Rename(oldName string, newName string) error {
	return nil
}

func (l *linux) Delete(name string) error {
	return nil
}

func (l *linux) Status(name string) api.Status {
	return api.Status{}
}
