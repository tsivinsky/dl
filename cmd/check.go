package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var checkCommand = &cobra.Command{
	Use:   "check",
	Short: "pulls new git commits",
	RunE: func(cmd *cobra.Command, args []string) error {
		var config Config
		if err := viper.Unmarshal(&config); err != nil {
			return err
		}

		if len(args) == 0 {
			return fmt.Errorf("app name is required")
		}

		app, err := findAppByName(config, args[0])
		if err != nil {
			return fmt.Errorf("app not found")
		}

		app.Destination = getAppDestination(app, config, configDir)

		if !isGitRepoExist(app.Destination) {
			return fmt.Errorf("repo doesn't exist, install first")
		}

		newCommits, err := fetchAppCommits(app)
		if err != nil {
			return fmt.Errorf("failed to fetch new commits: %v", err)
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
	},
}

func init() {
	rootCmd.AddCommand(checkCommand)
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
