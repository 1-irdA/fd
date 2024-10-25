package main

import (
	"flag"

	"github.com/garrou/fd/lib"
)

func main() {
	hidden := flag.Bool("i", false, "search hidden entries")
	recurse := flag.Bool("r", true, "find recursively")
	count := flag.Bool("c", false, "count the number of results")
	pattern := flag.Bool("p", false, "search pattern")
	extension := flag.String("e", "", "search files with extension")
	bench := flag.Bool("b", false, "get the execution time")
	exclude := flag.String("x", "", "exclude matching folders (space separated)")
	entryType := flag.String("t", "", "type of searched entry (f)ile / (d)irectory")
	location := flag.String("l", ".", "place to search")
	flag.Parse()

	walker := lib.NewWalker(
		flag.Args(),
		*recurse,
		*hidden,
		*count,
		*pattern,
		*bench,
		*extension,
		*exclude,
		*entryType,
		*location,
	)
	walker.Search()
}
