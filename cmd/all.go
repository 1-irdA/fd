package cmd

import (
	"github.com/garrou/fd/util"
	"github.com/spf13/cobra"
)

var (
	hidden  bool
	recurse bool
	allCmd = &cobra.Command{
		Use:   "all name|regex [path]",
		Short: "Find files or folders by name or regex",
		Long:  `Find files of folders by name or regex`,
		Run: func(cmd *cobra.Command, args []string) {
			find(args)
		},
	}
)

func init() {
	rootCmd.AddCommand(allCmd)
}

func find(args []string) {
	search, path := util.BindArgs(args)
	config := util.NewConfig(search, path, true, true, recurse, hidden)
	util.Search(config)
}