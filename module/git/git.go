package git

import (
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/config"
	"github.com/march1993/gohive/module"
	"github.com/march1993/gohive/module/linux"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type git struct{}

func init() {
	module.RegisterModule("git", &git{})
}

func (g *git) Create(name string) api.Status {
	unixname := linux.Prefix + name
	stdout, err := exec.Command("runuser",
		unixname,
		"-s", "/bin/bash",
		"-c", "cd ~ && git init .",
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
		"-s", "/bin/bash",
		"-c", "cd ~ && git status",
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
	errs := []string{}

	if stdout, err := exec.Command("runuser",
		unixname,
		"-s", "/bin/bash",
		"-c", "cd ~ && git init .",
	).CombinedOutput(); err != nil {
		errs = append(errs, string(stdout))
	}

	if stdout, err := exec.Command("runuser",
		unixname,
		"-s", "/bin/bash",
		"-c", "cd ~ && git checkout .",
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
	// TODO: Repair other problems
}

func (g *git) ListRemoved() []string {
	return []string{}
}

func GetGitUrl(name string) string {
	unixname := linux.Prefix + name
	return unixname + "@" + config.Get("server_name", "${server_name}") + ":~/.git"
}

const (
	SSH_DIR      = "/.ssh"
	SSH_KEY_FILE = SSH_DIR + "/authorized_keys"
	SSH_KEY_PERM = 0644
)

func SetGitKeys(name string, keys []string) api.Status {
	unixname := linux.Prefix + name
	home := linux.GetHomeDir(name)

	errs := []string{}

	if err := os.MkdirAll(home+SSH_DIR, 0700); err != nil {
		errs = append(errs, err.Error())
	}

	if err := ioutil.WriteFile(home+SSH_KEY_FILE, []byte(strings.Join(keys, "\n")), SSH_KEY_PERM); err != nil {
		errs = append(errs, err.Error())
	}

	if stdout, err := exec.Command("chown",
		unixname+":"+linux.Group,
		"-R",
		heom+SSH_DIR).CombinedOutput(); err != nil {
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

func GetGitKeys(name string) api.Status {
	home := linux.GetHomeDir(name)

	if stdout, err := exec.Command("ssh-keygen",
		"-lf",
		home+SSH_KEY_FILE).CombinedOutput(); err != nil {
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
