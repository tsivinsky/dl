package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/tsivinsky/dl/commander"
)

func findAppByName(config Conf, name string) (App, error) {
	for _, item := range config.DL {
		if item.Name == name {
			return item, nil
		}
	}

	return App{}, errors.New("app doesn't exist in config")
}

func isGitRepoExist(repoPath string) bool {
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return false
	}

	cmd := exec.Command("git", "status")
	cmd.Dir = repoPath
	if err := cmd.Run(); err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			if e.ProcessState.ExitCode() == 128 {
				return false
			}
			return false
		}
		return false
	}

	return true
}

func runInstructions(path string, instructions []string) error {
	for _, instruction := range instructions {
		parts := strings.Split(instruction, " ")
		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Dir = path
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf(`command "%s" failed: %v`, instruction, err)
		}
	}

	return nil
}

func getAppDestination(app App, configPath string) string {
	if app.Destination != "" {
		return app.Destination
	}

	return path.Join(configPath, app.Name)
}

func main() {
	cmder := commander.New()

	cmder.RegisterCommand("edit", "edit dl config", editCommand)
	cmder.RegisterCommand("list", "list all apps", listCommand)
	cmder.AddAliases("list", "ls")
	cmder.RegisterCommand("install", "install app defined in config", installCommand)
	cmder.RegisterCommand("update", "update app defined in config", updateCommand)
	cmder.RegisterCommand("check", "fetch new commits and print them", checkCommand)

	flag.Usage = cmder.Usage
	flag.Parse()

	if err := cmder.RunCommand(flag.Arg(0)); err != nil {
		fmt.Printf("%v\n", err)
	}
}
