package main

type ToysRepository interface {
	Save(name string) error
	ListAll() ([]string, error)
}
