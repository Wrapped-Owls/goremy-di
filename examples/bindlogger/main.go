package main

import (
	"fmt"
	"log"
	"os"

	"github.com/wrapped-owls/goremy-di/remy"
	"github.com/wrapped-owls/talaria-di/examples/bindlogger/utils"
)

func main() {
	lInject := utils.NewLoggerInjector(
		utils.DefaultLogger{
			ErrOutput: log.New(os.Stderr, "[Error]", log.LstdFlags|log.Lshortfile),
			StdOutput: log.New(os.Stdout, "", log.LstdFlags),
			Level:     utils.LevelDefault | utils.LevelInfo,
		},
	)
	remy.RegisterInstance(lInject, "The Lord of the Rings", "movie")
	remy.Register(
		lInject, remy.Factory(
			func(retriever remy.DependencyRetriever) string {
				return fmt.Sprintf("I love this movie with name `%s`", remy.Get[string](retriever, "movie"))
			},
		),
	)

	// Start retrieving the elements
	phrase := remy.Get[string](lInject)
	fmt.Println(phrase)

	if len(phrase) >= 2^5 {
		fmt.Println(remy.Get[bool](lInject))
	}
	result := remy.Get[uint16](lInject)
	fmt.Println(result)
}
