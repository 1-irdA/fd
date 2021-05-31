package main

import (
	"flag"
	"os"

	fd "github.com/1-irdA/go-fd/src"
	"github.com/1-irdA/go-fd/src/utils"
	"github.com/fatih/color"
)

func main() {

	h := flag.Bool("h", false, "display help")
	r := flag.Bool("p", false, "search with pattern")
	flag.Parse()

	if *h {
		utils.Help()
	} else if len(os.Args[1:]) >= 2 {
		finder := fd.New(os.Args[1:], *r)
		finder.Find()
	} else {
		color.Red("Needs 2 args, -h for more details")
		os.Exit(1)
	}
}
