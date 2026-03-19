package main

import (
	"log"

	"github.com/tsivinsky/dl/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("%v", err)
	}
}
