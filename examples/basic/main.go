package main

import (
	"database/sql"
	"github.com/wrapped-owls/goremy-di/remy"
	"github.com/wrapped-owls/talaria-di/examples/basic/core"
	_ "github.com/wrapped-owls/talaria-di/examples/basic/infra"
	"log"
)

// register all dependencies into the injector
func init() {
	remy.RegisterInstance(core.Injector, "That's a nice test")
}

func main() {
	// Executing create table query
	dbConn := remy.Get[*sql.DB](core.Injector)
	if _, err := dbConn.Exec("CREATE TABLE toys(id INTEGER, name VARCHAR(60))"); err != nil {
		log.Fatalln(err)
	}

	names := []string{
		"Teddy",
		"Teddy-bear",
		"Spider-Man",
	}

	repository := remy.Get[core.ToysRepository](core.Injector)
	for _, name := range names {
		err := repository.Save(name)
		if err != nil {
			log.Fatalln(err)
		}
	}

	log.Println(remy.Get[string](core.Injector))
	log.Println(repository.ListAll())
}
