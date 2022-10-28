package util

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
)

var occur uint32

func Search(config *config) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go walk(&wg, config.Path, config)
	wg.Wait()

	if config.Count {
		fmt.Println(fmt.Sprintf("%d occurrences", occur))
	}
}

func walk(wg *sync.WaitGroup, path string, config *config) {
	defer wg.Done()
	dir, errOp := os.Open(path)

	if errOp != nil {
		log.Fatal(errOp)
	}
	defer func(dir *os.File) {
		err := dir.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(dir)

	dirs, errRead := dir.Readdir(-1)

	if errRead != nil {
		log.Fatal(errRead)
	}
	for _, entry := range dirs {
		if isMatch(config, entry) {
			fmt.Println(filepath.Join(path, entry.Name()))

			if config.Count {
				atomic.AddUint32(&occur, 1)
			}
		}
		if entry.IsDir() && config.Recurse && canRecurse(config, entry) {
			wg.Add(1)
			go walk(wg, filepath.Join(path, entry.Name()), config)
		}
	}
}

func canRecurse(config *config, entry os.FileInfo) bool {
	if (config.Hidden && strings.HasPrefix(entry.Name(), ".")) || !strings.HasPrefix(entry.Name(), ".") {
		return true
	}
	return false
}

func isMatch(config *config, current os.FileInfo) bool {
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
