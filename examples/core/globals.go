package core

import (
	"github.com/wrapped-owls/talaria-di/gotalaria"
	"log"
)

var Injector gotalaria.Injector

// create a new instance of the injector
func init() {
	log.Println("Initializing injector")
	Injector = gotalaria.NewInjector()
}
