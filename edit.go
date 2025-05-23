package main

import (
	"log"
	"os"
	"os/exec"
)

func editCommand() error {
	_, confFile, err := setupConfig()
	if err != nil {
		log.Fatal(err)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	cmd := exec.Command(editor, confFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	return cmd.Run()
}
