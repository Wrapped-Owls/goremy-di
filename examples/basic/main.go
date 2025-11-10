package main

import (
	"database/sql"
	"log"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/wrapped-owls/goremy-di/examples/basic/core"
	_ "github.com/wrapped-owls/goremy-di/examples/basic/infra"
)

// register all dependencies into the injector
func init() {
	remy.RegisterInstance(core.Injector, "That's a nice test")
}

func main() {
	// Executing create table query
	dbConn := remy.MustGet[*sql.DB](core.Injector)
	if _, err := dbConn.Exec("CREATE TABLE toys(id INTEGER, name VARCHAR(60))"); err != nil {
		log.Fatalln(err)
	}

	names := []string{
		"Teddy",
		"Teddy-bear",
		"Spider-Man",
	}

	repository := remy.MustGet[core.ToysRepository](core.Injector)
	for _, name := range names {
		err := repository.Save(name)
		if err != nil {
			log.Fatalln(err)
		}
	}

	log.Println(remy.MustGet[string](core.Injector))
	log.Println(repository.ListAll())
}
