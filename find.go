package fd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

type config struct {
	path   string
	search string
	regex  bool
}

func New(path, search string, regex bool) *config {
	return &config{path: path, search: search, regex: regex}
}

func (c *config) isValid() bool {
	_, err := os.Lstat(c.path)
	return err == nil && c.search != ""
}

type find struct {
	conf *config
}

func Fd(config *config) *find {
	if !config.isValid() {
		log.Fatal("Invalid config")
	}
	return &find{conf: config}
}

func (f *find) Find() {
	wg := sync.WaitGroup{}
	start := time.Now()

	wg.Add(1)
	go f.worker(&wg, f.conf.path)
	wg.Wait()
	fmt.Println(time.Since(start))
}

func (f *find) worker(wg *sync.WaitGroup, path string) {
	defer wg.Done()
	dir, _ := os.Open(path)

	defer dir.Close()

	dirs, _ := dir.Readdir(-1)

	for _, entry := range dirs {
		if f.isMatch(entry.Name()) {
			fmt.Println(filepath.Join(path, entry.Name()))
		}
		if entry.IsDir() && entry.Name()[0] != '.' {
			wg.Add(1)
			go f.worker(wg, filepath.Join(path, entry.Name()))
		}
	}
}

func (f *find) isMatch(name string) (result bool) {
	if f.conf.regex {
		result, _ = regexp.MatchString(f.conf.search, name)
	} else {
		result = name == f.conf.search
	}
	return
}
