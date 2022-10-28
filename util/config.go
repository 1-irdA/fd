package util

import (
	"regexp"
)

type config struct {
	Search  string
	Path    string
	Regex   bool
	Match   func(string, string) bool
	File    bool
	Dir     bool
	Hidden  bool
	Recurse bool
	Count   bool
}

func NewConfig(search, path string, file, dir, recurse, hidden, count bool) *config {
	_, err := regexp.Compile(search)
	return &config{
		Search:  search,
		Path:    path,
		Regex:   err == nil,
		File:    file,
		Dir:     dir,
		Recurse: recurse,
		Hidden:  hidden,
		Count:   count,
	}
}

func BindArgs(args []string) (search, path string) {
	argsLen := len(args)

	search = ".*."
	path = "."

	if argsLen == 1 {
		search = args[0]
	} else if argsLen > 1 {
		search = args[0]
		path = args[1]
	}
	return
}
