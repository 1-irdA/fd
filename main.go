package main

import (
	"flag"

	"github.com/1-irdA/go-fd/finder"
	"github.com/fatih/color"
)

func printHelp() {
	color.Cyan(`
go-fd is a fast tool to find a directory, a file following some rules.
	
Usage:
	go run go-fd -h
	go run go-fd <dir> <to-search>
	go run go-fd -p <dir> <to-search>
The commands are:
	-h		Display this message
	-r		Search files that correspond to following regex
	-p      Path to search
	-s      Searched value
	-a      Display absolute path`)
}

var (
	path     = flag.String("p", "", "location to search")
	searched = flag.String("s", "", "thing to search")
	regex    = flag.Bool("r", false, "search with regex")
	absolute = flag.Bool("a", false, "display absolute path")
	help     = flag.Bool("h", false, "display help")
)

func main() {

	flag.Parse()

	if !*help {
		fd := finder.New(*path, *searched, *regex, *absolute)
		fd.Find()
	} else {
		printHelp()
	}
}
