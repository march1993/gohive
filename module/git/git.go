package git

import (
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/module"
	"github.com/march1993/gohive/module/linux"
	"os/exec"
)

type git struct{}

func init() {
	module.RegisterModule("git", &git{})
}

func (g *git) Create(name string) api.Status {
	unixname := linux.Prefix + name
	stdout, err := exec.Command("runuser",
		unixname,
		"-c",
		"cd . && git init .",
	).CombinedOutput()
	if err != nil {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: string(stdout),
		}
	} else {
		return api.Status{Status: api.STATUS_SUCCESS}
	}
}

func (g *git) Rename(oldName string, newName string) api.Status {
	// do nothing
	return api.Status{Status: api.STATUS_SUCCESS}
}

func (g *git) Remove(name string) api.Status {
	// do nothing
	return api.Status{Status: api.STATUS_SUCCESS}
}

func (g *git) Status(name string) api.Status {
	unixname := linux.Prefix + name
	stdout, err := exec.Command("runuser",
		unixname,
		"-c",
		"cd . && git status",
	).CombinedOutput()
	if err != nil {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: string(stdout),
		}
	} else {
		return api.Status{Status: api.STATUS_SUCCESS}
	}
}

func (g *git) Repair(name string) api.Status {
	unixname := linux.Prefix + name
	stdout, err := exec.Command("runuser",
		unixname,
		"-c",
		"cd . && git checkout .",
	).CombinedOutput()
	if err != nil {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: string(stdout),
		}
	} else {
		return api.Status{Status: api.STATUS_SUCCESS}
	}
	// TODO: Repair other problems
}

func (g *git) ListRemoved() []string {
	return []string{}
}
