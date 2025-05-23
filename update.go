package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
)

func updateApp(app App) error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = app.Destination
	return cmd.Run()
}

func updateCommand() error {
	confDir, confFile, err := setupConfig()
	if err != nil {
		log.Fatal(err)
	}

	conf, err := parseConfig(confFile)
	if err != nil {
		log.Fatal(err)
	}

	app, err := findAppByName(conf, flag.Arg(1))
	if err != nil {
		return fmt.Errorf("coulnd't find app with this name: %v", err)
	}

	dest := getAppDestination(app, confDir)

	repoExists := isGitRepoExist(dest)
	if !repoExists {
		return fmt.Errorf("repo doesn't exist, install first")
	}

	if err := updateApp(app); err != nil {
		return fmt.Errorf("couldn't pull changes: %v", err)
	}

	return runInstructions(dest, app.Build)
}
