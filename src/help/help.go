package help

import "fmt"

func Help() {
	fmt.Println(`
go-fd is a fast tool to find a directory, a file following some rules.
	
Usage:

	go run go-fd -h
	go run go-fd <dir> <to-search>
	go run go-fd -p <dir> <to-search>

The commands are:

	-h		Display this message
	-p		Search files that correspond to following regex`)
}
