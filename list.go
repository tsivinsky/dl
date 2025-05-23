package main

import (
	"fmt"
	"log"
)

func listCommand() error {
	_, confFile, err := setupConfig()
	if err != nil {
		log.Fatal(err)
	}

	conf, err := parseConfig(confFile)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range conf.DL {
		fmt.Printf("%s\n", item.Name)
	}

	return nil
}
