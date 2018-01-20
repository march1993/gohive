package module

import (
	"github.com/march1993/gohive/api"
)

type Module interface {
	Create(name string) error // create application
	Rename(oldName string, newName string) error
	Remove(name string) error      // remove application
	Status(name string) api.Status // show status

	Repair(name string) error // Repair an application when it's broken or copied to the hive directory
	ListRemoved() []string    // List application names whose have been removed
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
