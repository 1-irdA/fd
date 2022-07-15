package cmd

import (
	"github.com/garrou/fd/util"
	"github.com/spf13/cobra"
)

var (
	ext bool
)

var fileCmd = &cobra.Command{
	Use:   "file name|regex [path]",
	Short: "Find a file by name or regex",
	Long:  `Find a file by name or regex`,
	Run: func(cmd *cobra.Command, args []string) {
		findFile(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(fileCmd)
	rootCmd.Flags().BoolVarP(&ext, "extension", "e", false, "Print file by extension")
}

func findFile(cmd *cobra.Command, args []string) {
	search, path := util.BindArgs(args)
	config := util.NewConfig(search, path, true, false, false, hidden)
	util.Search(config)
}
