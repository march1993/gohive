package profile

import (
	"github.com/march1993/gohive/api"
	"github.com/march1993/gohive/config"
	"github.com/march1993/gohive/module"
	_ "github.com/march1993/gohive/module/linux"
	"github.com/march1993/gohive/util"
	"io/ioutil"
)

type profile struct{}

func init() {
	module.RegisterModule("profile", &profile{})
}

const (
	bashrcFilename = ".bashrc"
	bashrcTemplate = "./templates/bashrc"
)

func (p *profile) Create(name string) api.Status {
	bytes, err := ioutil.ReadFile(bashrcTemplate)
	if err != nil {
		panic(err.Error())
	}
	content := string(bytes)

	for _, handler := range module.Environ {
		kvs := handler(name)
		for _, line := range kvs {
			content += line + "\n"
		}
	}

	hash := util.Hash(content)
	config.AppConfigSet(name, "profile", "hash", hash)
	bashrcGenerated := config.GetHomeDir(name) + "/" + bashrcFilename

	if err = ioutil.WriteFile(bashrcGenerated, []byte(content), 0644); err != nil {
		panic(err.Error())
	}

	return api.Status{Status: api.STATUS_SUCCESS}
}

func (p *profile) Rename(oldName, newName string) api.Status {
	// do nothing
	return api.Status{Status: api.STATUS_SUCCESS}
}

func (p *profile) Remove(name string) api.Status {
	// do nothing
	return api.Status{Status: api.STATUS_SUCCESS}
}

func (p *profile) Status(name string) api.Status {
	hash := config.AppConfigGet(name, "profile", "hash", "")
	bashrcGenerated := config.GetHomeDir(name) + "/" + bashrcFilename

	bytes, err := ioutil.ReadFile(bashrcGenerated)

	if err != nil {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: err.Error(),
		}
	}

	if hash != util.Hash(string(bytes)) {
		return api.Status{
			Status: api.STATUS_FAILURE,
			Reason: api.PROFILE_BASHRC_EXPIRED,
		}
	} else {
		return api.Status{Status: api.STATUS_SUCCESS}
	}

}

func (p *profile) Repair(name string) api.Status {
	return p.Create(name)
}

func (p *profile) ListRemoved() []string {
	return []string{}
}
