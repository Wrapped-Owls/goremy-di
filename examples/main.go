package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wrapped-owls/talaria-di/examples/repositories"
	"github.com/wrapped-owls/talaria-di/gotalaria"
	"log"
)

var injector gotalaria.Injector

// create a new instance of the injector
func init() {
	injector = gotalaria.NewInjector()
}

// Create an instance of the database connection
func init() {
	db, err := sql.Open("sqlite3", "file:locked.sqlite?cache=shared&mode=memory")
	if err != nil {
		panic(err)
	}

	gotalaria.RegisterSingleton(injector, db)
}

// register all dependencies into the injector
func init() {
	gotalaria.RegisterInstance(injector, "That's a nice test")
	gotalaria.Register(injector, gotalaria.Factory(func(retriever gotalaria.DependencyRetriever) ToysRepository {
		return repositories.NewToysDbRepository(gotalaria.Get[*sql.DB](injector))
	}))
}

func main() {
	// Executing create table query
	dbConn := gotalaria.Get[*sql.DB](injector)
	if _, err := dbConn.Exec("CREATE TABLE toys(id INTEGER, name VARCHAR(60))"); err != nil {
		log.Fatalln(err)
	}

	names := []string{
		"Teddy",
		"Teddybear",
		"Spider-Man",
	}

	repository := gotalaria.Get[ToysRepository](injector)
	for _, name := range names {
		err := repository.Save(name)
		if err != nil {
			log.Fatalln(err)
		}
	}

	log.Println(gotalaria.Get[string](injector))
	log.Println(repository.ListAll())
}
