package cmd

import (
	"github.com/garrou/fd/util"
	"github.com/spf13/cobra"
)

var folderCmd = &cobra.Command{
	Use:   "folder  name|regex [path]",
	Short: "Find a folder by name or regex",
	Long:  `Find a folder by name or regex`,
	Run: func(cmd *cobra.Command, args []string) {
		findFolder(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(folderCmd)
}

func findFolder(cmd *cobra.Command, args []string) {
	search, path := util.BindArgs(args)
	config := util.NewConfig(search, path, false, true, recurse, hidden)
	util.Search(config)
}
