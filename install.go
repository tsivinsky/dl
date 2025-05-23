package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func cloneGitRepo(repoUrl string, dest string) error {
	cmd := exec.Command("git", "clone", repoUrl, dest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installCommand() error {
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
	if repoExists {
		fmt.Printf("%s already exists, skipping cloning\n", dest)
	}

	if !repoExists {
		err = cloneGitRepo(app.URL, dest)
		if err != nil {
			return fmt.Errorf("couldn't clone git repo: %v", err)
		}
	}

	return runInstructions(dest, app.Build)
}
