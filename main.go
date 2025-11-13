package main

import (
	"cover-letter-templates/pkg/gui"
	"log"
)

// Project represents a single project entry.
type Project struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	URL         string `yaml:"url"`
}

func main() {
	// err := cli.Run()
	//if err != nil {
	//	log.Fatal(err)
	//}
	err := gui.Run()
	if err != nil {
		log.Fatal(err)
	}
}
