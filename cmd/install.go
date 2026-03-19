package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var installCommand = &cobra.Command{
	Use:   "install",
	Short: "installs app defined in config",
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

		dest := getAppDestination(app, config, configDir)

		if isGitRepoExist(dest) {
			fmt.Println("Repository already exists, skipping cloning")
		} else {
			if err := cloneGitRepo(app.URL, dest); err != nil {
				return fmt.Errorf("failed to clone git repository: %v", err)
			}
		}

		return runAppInstructions(dest, app.Build)
	},
}

func init() {
	rootCmd.AddCommand(installCommand)
}

func cloneGitRepo(repoUrl string, dest string) error {
	cmd := exec.Command("git", "clone", repoUrl, dest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
