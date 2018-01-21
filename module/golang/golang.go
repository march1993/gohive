package golang

import (
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/module"
	_ "github.com/march1993/gohive/module/git"
	"github.com/march1993/gohive/module/linux"
	"os/exec"
	"strings"
)

type golang struct{}

func init() {
	module.RegisterModule("golang", &golang{})
}

func (g *golang) Create(name string) api.Status {
	// do nothing
	return api.Status{Status: api.STATUS_SUCCESS}
}

func (g *golang) Rename(oldName string, newName string) api.Status {
	// do nothing
	return api.Status{Status: api.STATUS_SUCCESS}
}

func (g *golang) Remove(name string) api.Status {
	// do nothing
	return api.Status{Status: api.STATUS_SUCCESS}
}

func (g *golang) Status(name string) api.Status {
	unixname := linux.Prefix + name
	stdout, err := exec.Command("runuser",
		unixname,
		"-s", "/bin/bash",
		"-c", "go version",
	).CombinedOutput()
	if err != nil {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: string(stdout),
		}
	} else {
		return api.Status{
			Status: api.STATUS_SUCCESS,
			Result: string(stdout),
		}
	}
}

func (g *golang) Repair(name string) api.Status {
	unixname := linux.Prefix + name
	errs := []string{}

	if stdout, err := exec.Command("runuser",
		unixname,
		"-s", "/bin/bash",
		"-c", "cd ~ && go get .",
	).CombinedOutput(); err != nil {
		errs = append(errs, string(stdout))
	}

	if len(errs) > 0 {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: strings.Join(errs, "\n"),
		}
	} else {
		return api.Status{Status: api.STATUS_SUCCESS}
	}
}

func (g *golang) ListRemoved() []string {
	return []string{}
}
