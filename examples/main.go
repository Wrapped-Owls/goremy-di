package main

import (
	"database/sql"
	"github.com/wrapped-owls/talaria-di/examples/core"
	_ "github.com/wrapped-owls/talaria-di/examples/infra"
	"github.com/wrapped-owls/talaria-di/gotalaria"
	"log"
)

// register all dependencies into the injector
func init() {
	gotalaria.RegisterInstance(core.Injector, "That's a nice test")
}

func main() {
	// Executing create table query
	dbConn := gotalaria.Get[*sql.DB](core.Injector)
	if _, err := dbConn.Exec("CREATE TABLE toys(id INTEGER, name VARCHAR(60))"); err != nil {
		log.Fatalln(err)
	}

	names := []string{
		"Teddy",
		"Teddybear",
		"Spider-Man",
	}

	repository := gotalaria.Get[core.ToysRepository](core.Injector)
	for _, name := range names {
		err := repository.Save(name)
		if err != nil {
			log.Fatalln(err)
		}
	}

	log.Println(gotalaria.Get[string](core.Injector))
	log.Println(repository.ListAll())
}
