package infra

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wrapped-owls/talaria-di/examples/core"
	"github.com/wrapped-owls/talaria-di/examples/infra/repositories"
	"github.com/wrapped-owls/talaria-di/gotalaria"
	"log"
)

func init() {
	log.Println("Registering repositories")
	gotalaria.Register(
		core.Injector,
		gotalaria.Factory(func(retriever gotalaria.DependencyRetriever) core.ToysRepository {
			return repositories.NewToysDbRepository(gotalaria.Get[*sql.DB](retriever))
		}),
	)
}

// Create an instance of the database connection
func init() {
	gotalaria.Register(
		core.Injector,
		gotalaria.Singleton(func(retriever gotalaria.DependencyRetriever) *sql.DB {
			db, err := sql.Open("sqlite3", "file:locked.sqlite?cache=shared&mode=memory")
			if err != nil {
				panic(err)
			}
			return db
		}),
	)
}
