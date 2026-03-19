package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var updateCommand = &cobra.Command{
	Use:   "update",
	Short: "updates app",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var apps []App
		if err := viper.UnmarshalKey("dl", &apps); err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		options := []string{}
		for _, app := range apps {
			if strings.HasPrefix(app.Name, toComplete) {
				options = append(options, app.Name)
			}
		}

		return options, cobra.ShellCompDirectiveDefault
	},
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

		if err := updateApp(app); err != nil {
			return fmt.Errorf("failed to pull changed: %v", err)
		}

		return runAppInstructions(app.Destination, app.Build)
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
