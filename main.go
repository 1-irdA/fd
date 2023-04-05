package main

import (
	"runtime"

	"github.com/garrou/fd/cmd"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() / 4)
	cmd.Execute()
}
