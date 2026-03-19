package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "dl",
	Short: "download git repositories easier",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	if err := initConfig(); err != nil {
		log.Fatalf("%v", err)
	}
}

type App struct {
	Name        string   `yaml:"name" mapstructure:"name"`
	URL         string   `yaml:"url" mapstructure:"url"`
	Build       []string `yaml:"build" mapstructure:"build"`
	Destination string   `yaml:"dest" mapstructure:"dest"`
}

type Config struct {
	RootDir string `yaml:"root_dir" mapstructure:"root_dir"`
	DL      []App  `yaml:"dl" mapstructure:"dl"`
}

var (
	configDir  string
	configFile string
)

func initConfig() error {
	viper.SetEnvPrefix("DL")

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get user config directory: %v", err)
	}

	configDir = path.Join(userConfigDir, "dl")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.Mkdir(configDir, 0777); err != nil {
			return fmt.Errorf("failed to create config directory: %v", err)
		}
	}

	configFile = path.Join(configDir, "config.yml")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if _, err := os.Create(configFile); err != nil {
			return fmt.Errorf("failed to create config file: %v", err)
		}
	}

	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	return nil
}
