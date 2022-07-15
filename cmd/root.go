package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fd",
	Short: "fd is a file and folder finder",
	Long: `A fast file and folder finder built in Go.
			Complete documentation is available at https://github.com/garrou/fd`,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&hidden, "hidden", "c", false, "Print hidden file or folder")
	rootCmd.PersistentFlags().BoolVarP(&recurse, "recurse", "r", false, "Search recursively")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
