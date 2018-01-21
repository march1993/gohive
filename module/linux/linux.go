package linux

import (
	"errors"
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/config"
	"github.com/march1993/gohive/module"
	"os/exec"
)

type linux struct{}

const (
	Prefix = "gohive_app_"
	Group  = "gohive_app"
)

func init() {
	module.RegisterModule("linux", &linux{})

	cmd := exec.Command("groupadd", "-f", Group)
	if stdout, err := cmd.CombinedOutput(); err != nil {
		panic(string(stdout) + err.Error())
	}
}
func (l *linux) Create(name string) error {

	if l.Status(name).Status == api.APP_NON_EXIST {
		unixname := Prefix + name

		cmd := exec.Command("useradd", "-b", config.APP_DIR, "-m", "-s", config.SSH_SHELL, "-G", Group, unixname)
		stdout, err := cmd.CombinedOutput()

		if err != nil {
			return errors.New(string(stdout) + err.Error())
		}

		return nil

	} else {
		return errors.New(api.APP_ALREADY_EXISTING)
	}
}

func (l *linux) Rename(oldName string, newName string) error {
	return nil
}

func (l *linux) Remove(name string) error {
	return nil
}

func (l *linux) Status(name string) api.Status {
	return api.Status{
		Status: api.APP_NON_EXIST,
	}
}

func (l *linux) Repair(name string) error {
	return nil
}

func (l *linux) ListRemoved() []string {
	return []string{}
}
