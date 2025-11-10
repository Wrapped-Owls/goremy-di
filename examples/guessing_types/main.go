package main

import (
	"fmt"
	"log"
	"reflect"
	"slices"
	"strings"

	"github.com/wrapped-owls/goremy-di/remy"
	"github.com/wrapped-owls/goremy-di/remy/test/fixtures"

	"github.com/wrapped-owls/goremy-di/examples/guessing_types/models"
)

func guessBindType(i remy.Injector) {
	expectedElement := fixtures.CountryLanguage{Language: "pt_BR"}
	remy.Register(
		i,
		remy.Factory(func(retriever remy.DependencyRetriever) (fixtures.CountryLanguage, error) {
			return expectedElement, nil
		}),
	)

	if _, err := remy.DoGet[fixtures.Language](i); err != nil {
		log.Printf("As said before, it shouldn't work. Received error `%v`", err)
	}
}

func main() {
	i := remy.NewInjector(remy.Config{DuckTypeElements: true})

	type Calculator[T interface {
		int32 | int64 | uint32 | uint64 | float32 | float64
	}] interface {
		Calculate(T) T
	}

	// This feature only works if the element was bind as Instance,
	// otherwise because of Go limitations it could not be guessed
	remy.RegisterInstance(i, models.FibonacciCalculator{})
	remy.RegisterInstance(i, models.RootCalculator{Radical: 2, Precision: 6})

	// Perform calculator that should be fibonacci
	{
		calculator := remy.Get[Calculator[uint64]](i)
		fmt.Println("Fibonacci: ", calculator.Calculate(11))
	}

	// Perform calculator that should be a square
	{
		calculator := remy.Get[Calculator[float64]](i)
		fmt.Println("Square: ", calculator.Calculate(1764))
	}

	// Calling in a way that doesn't work
	guessBindType(i)

	// Starting from here it will raise an error, as it will find many possible elements
	remy.RegisterInstance(i, models.SineCalculator{})

	{
		calculator, err := remy.DoGet[Calculator[float64]](i)
		if err == nil {
			fmt.Println("Sine: ", calculator.Calculate(180))
		} else {
			log.Printf("Error raisen when getting %T: %v", &calculator, err)
		}
	}

	// If you want to get and find more than one element, use remy.GetAll
	{
		listOfCalculators, err := remy.DoGetAll[Calculator[float64]](i)
		if err != nil {
			log.Printf("An error was raised, it is not expected")
		} else {
			var (
				value             float64 = 19980707
				totalCalculations uint8
			)
			slices.SortFunc(listOfCalculators, func(a, b Calculator[float64]) int {
				aType := reflect.TypeOf(a)
				bType := reflect.TypeOf(b)
				return strings.Compare(aType.Name(), bType.Name())
			})
			for _, calculator := range listOfCalculators {
				value = calculator.Calculate(value)
				totalCalculations += 1
			}
			fmt.Printf("Result after perform %d calculations: `%v`\n", totalCalculations, value)
		}
	}

	// Once we register anything as the actual interface, it will start to get only it
	remy.RegisterInstance[Calculator[uint64]](i, models.FactorialCalculator{})
	{
		calculator, err := remy.DoGet[Calculator[uint64]](i)
		if err == nil {
			fmt.Printf("Perform operation in type %T: %v\n", calculator, calculator.Calculate(5))
		} else {
			log.Fatalf("This error should not be raisen: %v", err)
		}
	}
}
