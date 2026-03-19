package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCommand = &cobra.Command{
	Use:     "list",
	Short:   "lists all available apps",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		var config Config
		if err := viper.Unmarshal(&config); err != nil {
			return err
		}

		tw := tabwriter.NewWriter(os.Stdout, 0, 2, 4, ' ', 0)
		defer tw.Flush()
		for _, item := range config.DL {
			info, err := getRepoInfo(item.Destination)
			if err != nil {
				fmt.Fprintf(tw, "%s\t\tfailed to retrieve latest commit\t[%v]\n", item.Name, err)
				continue
			}

			fmt.Fprintf(tw, "%s\t\t%s\t%s\n", item.Name, info.hash, info.relativeDate)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCommand)
}

type commitInfo struct {
	hash         string
	relativeDate string
}

func getRepoInfo(repoPath string) (*commitInfo, error) {
	cmd := exec.Command("git", "log", "-n", "1", `--pretty=format:%H;%ar`)
	cmd.Dir = repoPath

	b, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	s := strings.Split(string(b), ";")
	if len(s) != 2 {
		return nil, fmt.Errorf("invalid info format")
	}

	info := &commitInfo{
		hash:         s[0],
		relativeDate: s[1],
	}
	return info, nil
}
