package core

import (
	"github.com/wrapped-owls/goremy-di/remy"
	"log"
)

var Injector remy.Injector

// create a new instance of the injector
func init() {
	log.Println("Initializing injector")
	Injector = remy.NewInjector()
}
