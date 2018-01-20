package module

import (
	"github.com/march1993/gohive/api"
)

type Module interface {
	Create(name string) error
	Rename(oldName string, newName string) error
	Delete(name string) error
	Status(name string) api.Status
	Repair(name string) error
}

var Modules = map[string]Module{}

// var ModuleList []Module

func RegisterModule(name string, module Module) {

	Modules[name] = module
	// ModuleList = append(ModuleList, module)

}

func GetAppStatus(app string) map[string]api.Status {
	ret := map[string]api.Status{}

	for name, module := range Modules {
		ret[name] = module.Status(app)
	}

	return ret
}
