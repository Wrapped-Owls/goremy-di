package main

import (
	"fmt"
	"log"
	"os"

	"github.com/wrapped-owls/goremy-di/remy"

	"github.com/wrapped-owls/goremy-di/examples/bindlogger/utils"
)

func main() {
	lInject := utils.NewLoggerInjector(
		utils.DefaultLogger{
			ErrOutput: log.New(os.Stderr, "[Error] ", log.LstdFlags|log.Lshortfile),
			StdOutput: log.New(os.Stdout, "", log.LstdFlags),
			Level:     utils.LevelDefault | utils.LevelInfo,
		},
		remy.Config{},
	)
	remy.RegisterInstance(lInject, "The Lord of the Rings", "movie")
	remy.Register(
		lInject, remy.Factory(
			func(retriever remy.DependencyRetriever) (result string, err error) {
				var movieName string
				movieName, err = remy.Get[string](retriever, "movie")
				result = fmt.Sprintf("I love this movie with name `%s`", movieName)
				return
			},
		),
	)

	// Start retrieving the elements
	phrase := remy.MustGet[string](lInject)
	fmt.Println(phrase)

	if len(phrase) >= 2^5 {
		fmt.Println(remy.MaybeGet[bool](lInject))
	}

	result := remy.MaybeGet[error](lInject)
	fmt.Println(result)
}
