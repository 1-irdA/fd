package cmd

import (
	"github.com/garrou/fd/lib"
	"github.com/spf13/cobra"
)

var (
	ext     bool
	fileCmd = &cobra.Command{
		Use:   "file name|regex [path]",
		Short: "Find files by name or regex",
		Long:  `Find files by name or regex`,
		Run: func(cmd *cobra.Command, args []string) {
			findFile(args)
		},
	}
)

func init() {
	rootCmd.AddCommand(fileCmd)
}

func findFile(args []string) {
	search, path := lib.BindArgs(args)
	config := lib.NewConfig(search, path, true, false, recurse, hidden, count)
	lib.Search(config)
}
