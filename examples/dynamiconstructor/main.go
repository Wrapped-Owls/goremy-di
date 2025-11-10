package main

import (
	"log"
	"time"

	"github.com/wrapped-owls/goremy-di/remy"
)

func registerInjections(ij remy.Injector) {
	remy.RegisterInstance(ij, "That's a nice test")
	remy.RegisterConstructor(ij, remy.Factory[time.Time], time.Now)
	remy.RegisterConstructorArgs2(ij, remy.Factory[Note], NewAnnotation)
	remy.RegisterConstructorArgs1(ij, remy.LazySingleton[FolderChecker], NewFolderChecker)
}

func main() {
	inj := remy.NewInjector(remy.Config{UseReflectionType: false, DuckTypeElements: true})
	// Registering injections
	registerInjections(inj)

	folderChecker := remy.MustGetWith[FolderChecker](inj, func(injector remy.Injector) error {
		remy.RegisterInstance(injector, "Trying to retrieve program current folder")
		return nil
	})
	absPath := folderChecker.RunningAbsolute()
	log.Println(absPath)
	log.Println(remy.MustGet[string](inj))
}
