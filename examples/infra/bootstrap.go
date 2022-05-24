package infra

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wrapped-owls/goremy-di/remy"
	"github.com/wrapped-owls/talaria-di/examples/core"
	"github.com/wrapped-owls/talaria-di/examples/infra/repositories"
	"log"
)

func init() {
	log.Println("Registering repositories")
	remy.Register(
		core.Injector,
		remy.Factory(func(retriever remy.DependencyRetriever) core.ToysRepository {
			return repositories.NewToysDbRepository(remy.Get[*sql.DB](retriever))
		}),
	)
}

// Create an instance of the database connection
func init() {
	remy.Register(
		core.Injector,
		remy.Singleton(func(retriever remy.DependencyRetriever) *sql.DB {
			db, err := sql.Open("sqlite3", "file:locked.sqlite?cache=shared&mode=memory")
			if err != nil {
				panic(err)
			}
			return db
		}),
	)
}
