package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var installCommand = &cobra.Command{
	Use:   "install",
	Short: "installs app defined in config",
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

		if isGitRepoExist(app.Destination) {
			fmt.Println("Repository already exists, skipping cloning")
		} else {
			if err := cloneGitRepo(app.URL, app.Destination); err != nil {
				return fmt.Errorf("failed to clone git repository: %v", err)
			}
		}

		return runAppInstructions(app.Destination, app.Build)
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
