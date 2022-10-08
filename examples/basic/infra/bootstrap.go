package infra

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wrapped-owls/goremy-di/remy"
	"github.com/wrapped-owls/talaria-di/examples/basic/core"
	"github.com/wrapped-owls/talaria-di/examples/basic/infra/repositories"
)

func init() {
	log.Println("Registering repositories")
	remy.Register(
		core.Injector,
		remy.Factory(
			func(retriever remy.DependencyRetriever) (core.ToysRepository, error) {
				db, err := remy.DoGet[*sql.DB](retriever)
				return repositories.NewToysDbRepository(db), err
			},
		),
	)
}

// Create an instance of the database connection
func init() {
	remy.Register(
		core.Injector,
		remy.Singleton(
			func(retriever remy.DependencyRetriever) (*sql.DB, error) {
				return sql.Open("sqlite3", "file:locked.sqlite?cache=shared&mode=memory")
			},
		),
	)
}
