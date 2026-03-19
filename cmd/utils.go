package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

func findAppByName(config Config, name string) (App, error) {
	for _, item := range config.DL {
		if item.Name == name {
			return item, nil
		}
	}

	return App{}, errors.New("app doesn't exist in config")
}

func runAppInstructions(path string, instructions []string) error {
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

func getAppDestination(app App, conf Config, configPath string) string {
	if app.Destination != "" {
		return app.Destination
	}

	if conf.RootDir != "" {
		return path.Join(conf.RootDir, app.Name)
	}

	return path.Join(configPath, app.Name)
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
