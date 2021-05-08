package main

import (
	"flag"
	"fmt"
	"os"

	fd "github.com/1-irdA/go-fd/src"
	"github.com/1-irdA/go-fd/src/color"
	"github.com/1-irdA/go-fd/src/help"
)

func main() {

	h := flag.Bool("h", false, "display help")
	r := flag.Bool("p", false, "search with pattern")
	flag.Parse()

	if *h {
		help.Help()
	} else if len(os.Args[1:]) >= 2 {
		finder := fd.New(os.Args[1:], *r)
		finder.Find()
	} else {
		fmt.Println(color.Red + "Needs 2 args, -h for more details" + color.White)
	}
}
