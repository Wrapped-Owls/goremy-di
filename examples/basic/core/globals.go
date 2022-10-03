package core

import (
	"log"

	"github.com/wrapped-owls/goremy-di/remy"
)

var Injector remy.Injector

// create a new instance of the injector
func init() {
	log.Println("Initializing injector")
	Injector = remy.NewInjector()
}
