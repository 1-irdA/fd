package utils

import (
	"os"

	"github.com/fatih/color"
)

func Err(msg string) {
	color.Red(msg)
	os.Exit(1)
}
