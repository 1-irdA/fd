package main

import (
	"flag"

	"github.com/1-irdA/fd/finder"
	"github.com/fatih/color"
)

func printHelp() {
	color.Cyan(`
go-fd is a fast tool to find a directory, a file following some rules.
	
Usage:
	go run go-fd -h
	go run go-fd [-p=path] -s=to-search [-f] [-d] [-a] [-b] [-r]
The commands are:
	-h	Display this message
	-r	Search files that correspond to following regex
	-p	Path to search
	-s	Searched value
	-a	Display absolute path
	-d	Search directories
	-f	Search files
	-b	Display benchmark`)
}

var (
	path     = flag.String("p", "", "location to search")
	searched = flag.String("s", "", "thing to search")
	dir      = flag.Bool("d", false, "search folders")
	file     = flag.Bool("f", false, "search files")
	regex    = flag.Bool("r", false, "search with regex")
	absolute = flag.Bool("a", false, "display absolute path")
	bench    = flag.Bool("b", false, "display benchmark")
	help     = flag.Bool("h", false, "display help")
)

func main() {

	flag.Parse()

	if !*help {
		options := finder.Options{
			File:     *file,
			Dir:      *dir,
			Regex:    *regex,
			Absolute: *absolute,
			Bench:    *bench,
		}
		finder.New(*path, *searched, options).Find()
	} else {
		printHelp()
	}
}
