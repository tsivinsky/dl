package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func fetchAppCommits(app App) ([]string, error) {
	fetchCmd := exec.Command("git", "fetch", "--all")
	fetchCmd.Dir = app.Destination
	if err := fetchCmd.Run(); err != nil {
		return nil, fmt.Errorf("command to fetch updates failed: %v", err)
	}

	logCmd := exec.Command("git", "log", "HEAD..origin", "--oneline")
	logCmd.Dir = app.Destination
	out, err := logCmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git log command failed: %v", err)
	}

	lines := strings.Split(string(out), "\n")
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines, nil
}

func checkCommand() error {
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

	if !isGitRepoExist(dest) {
		return fmt.Errorf("repo doesn't exist, install first")
	}

	newCommits, err := fetchAppCommits(app)
	if err != nil {
		return fmt.Errorf("couldn't fetch new commits: %v", err)
	}

	if len(newCommits) > 0 {
		fmt.Printf("%s has new commits\n", app.Name)
		for _, commit := range newCommits {
			fmt.Println(commit)
		}
	} else {
		fmt.Printf("%s has no new commits\n", app.Name)
	}

	return nil
}
