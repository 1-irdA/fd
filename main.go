package main

import (
	"runtime"

	"github.com/garrou/fd/cmd"
)

func main() {
	runtime.GOMAXPROCS(2)
	cmd.Execute()
}
