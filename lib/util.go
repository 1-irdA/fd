package lib

import (
	"log"
	"path/filepath"
)

func checkArgs(args []string) (search, path string) {
	search, path = "", "."

	if len(args) == 1 {
		search = args[0]
	} else if len(args) == 2 {
		search, path = args[0], args[1]
		path, err := filepath.Abs(path)

		if err != nil {
			log.Fatalf("%s is not a valid path", path)
		}
	}
	return
}
