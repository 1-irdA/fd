package cmd

import (
	"fmt"
	"os"

	"github.com/garrou/fd/util"
	"github.com/spf13/cobra"
)

var (
	hidden  bool
	recurse bool
)

var rootCmd = &cobra.Command{
	Use:   "fd",
	Short: "fd is a file and folder finder",
	Long: `A fast file and folder finder built in Go.
			Complete documentation is available at https://github.com/garrou/fd`,
	Run: func(cmd *cobra.Command, args []string) {
		find(cmd, args)
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&hidden, "hidden", "c", false, "Print hidden file or folder")
	rootCmd.PersistentFlags().BoolVarP(&recurse, "recurse", "r", true, "Search recursively")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func find(cmd *cobra.Command, args []string) {
	search, path := util.BindArgs(args)
	config := util.NewConfig(search, path, true, true, recurse, hidden)
	util.Search(config)
}
