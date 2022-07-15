package cmd

import (
	"github.com/garrou/fd/util"
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
	fileCmd.Flags().BoolVarP(&ext, "extension", "e", false, "Print file by extension")
	rootCmd.AddCommand(fileCmd)
}

func findFile(args []string) {
	search, path := util.BindArgs(args)
	config := util.NewConfig(search, path, true, false, recurse, hidden)
	util.Search(config)
}
