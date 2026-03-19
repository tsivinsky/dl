package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateCommand = &cobra.Command{
	Use:   "update",
	Short: "updates app",
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

		if !isGitRepoExist(dest) {
			return fmt.Errorf("repo doesn't exist, install first")
		}

		if err := updateApp(app); err != nil {
			return fmt.Errorf("failed to pull changed: %v", err)
		}

		return runAppInstructions(dest, app.Build)
	},
}

func init() {
	rootCmd.AddCommand(updateCommand)
}

func updateApp(app App) error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = app.Destination
	return cmd.Run()
}
