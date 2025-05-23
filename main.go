package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/tsivinsky/dl/commander"
	"gopkg.in/yaml.v3"
)

type App struct {
	Name        string   `yaml:"name"`
	URL         string   `yaml:"url"`
	Build       []string `yaml:"build"`
	Destination string   `yaml:"dest"`
}

type Conf struct {
	DL []App `yaml:"dl"`
}

func setupConfig() (string, string, error) {
	userConfDir, err := os.UserConfigDir()
	if err != nil {
		return "", "", err
	}

	confDir := path.Join(userConfDir, "dl")
	if _, err := os.Stat(confDir); os.IsNotExist(err) {
		err = os.Mkdir(confDir, 0755)
		if err != nil {
			return "", "", err
		}
	}

	confFile := path.Join(confDir, "config.yml")
	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		_, err = os.OpenFile(confFile, os.O_CREATE, 0644)
		if err != nil {
			return "", "", err
		}
	}

	return confDir, confFile, nil
}

func parseConfig(configPath string) (Conf, error) {
	var conf Conf

	f, err := os.OpenFile(configPath, os.O_RDONLY, 0644)
	if err != nil {
		return conf, err
	}
	defer f.Close()

	if err := yaml.NewDecoder(f).Decode(&conf); err != nil {
		return conf, err
	}

	return conf, nil
}

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

func cloneGitRepo(repoUrl string, dest string) error {
	cmd := exec.Command("git", "clone", repoUrl, dest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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

func updateApp(app App) error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = app.Destination
	return cmd.Run()
}

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
