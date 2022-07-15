package util

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

func Search(config *config) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go walk(&wg, config.Path, config)
	wg.Wait()
}

func walk(wg *sync.WaitGroup, path string, config *config) {
	defer wg.Done()
	dir, errOp := os.Open(path)

	if errOp != nil {
		return
	}
	defer dir.Close()

	dirs, errRead := dir.Readdir(-1)

	if errRead != nil {
		return
	}
	for _, entry := range dirs {
		if match(config, entry) {
			fmt.Println(filepath.Join(path, entry.Name()))
		}
		if entry.IsDir() && config.Recurse && recurse(config, entry) {
			wg.Add(1)
			go walk(wg, filepath.Join(path, entry.Name()), config)
		}
	}
}

func recurse(config *config, entry fs.FileInfo) bool {
	if (config.Hidden && strings.HasPrefix(entry.Name(), ".")) || !strings.HasPrefix(entry.Name(), ".") {
		return true
	}
	return false
}

func match(config *config, current fs.FileInfo) bool {
	var ok bool

	if !config.Hidden && strings.HasPrefix(current.Name(), ".") {
		return false
	}

	if config.Regex {
		ok, _ = regexp.MatchString(config.Search, current.Name())
	} else {
		ok = strings.Contains(current.Name(), config.Search)
	}

	if config.File && config.Dir && ok {
		return true
	} else if config.Dir && current.IsDir() && ok {
		return true
	} else if config.File && !current.IsDir() && ok {
		return true
	}
	return false
}
