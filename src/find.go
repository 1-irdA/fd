package src

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/1-irdA/go-fd/src/utils"
	"github.com/fatih/color"
)

type Finder interface {
	Find()
	details()
}

type find struct {
	path     string
	searched string
	nbFiles  uint32
	reg      bool
	wg       sync.WaitGroup
	elapsed  time.Duration
}

func New(args []string, r bool) *find {
	path, searched := getPathAndSearch(args, r)
	_, err := os.Lstat(path)
	if err != nil {
		utils.Err("Invalid path")
	}
	return &find{path: path, searched: searched, reg: r}
}

func (f *find) Find() {
	defer f.details()
	f.checkRegIfNeed()
	start := time.Now()
	f.wg.Add(1)
	f.worker(f.path)
	f.wg.Wait()
	f.elapsed = time.Since(start)
}

func (f *find) worker(dirPath string) {
	defer f.wg.Done()
	dirs, err := os.ReadDir(dirPath)
	if err != nil {
		utils.Err("Error during files browsing")
	}
	for _, entry := range dirs {
		atomic.AddUint32(&f.nbFiles, 1)
		if entry.IsDir() {
			f.wg.Add(1)
			go f.worker(path.Join(dirPath, entry.Name()))
		} else if f.correspond(entry) {
			f.printPath(dirPath, entry)
		}
	}
}

func (f *find) details() {
	color.Yellow("Files browsed %d, search duration : %v\n", f.nbFiles, f.elapsed)
}

func (f *find) checkRegIfNeed() {
	if f.reg {
		if _, err := regexp.Compile(f.searched); err != nil {
			utils.Err("Invalid regex : " + f.searched)
		}
	}
}

func (f *find) correspond(entry fs.DirEntry) bool {
	var result bool = false
	if f.reg {
		match, err := regexp.MatchString(f.searched, entry.Name())
		if err != nil {
			utils.Err("Regex error")
		}
		if match {
			result = true
		}
	} else {
		if strings.Contains(entry.Name(), f.searched) {
			result = true
		}
	}
	return result
}

func (f *find) printPath(dirPath string, entry fs.DirEntry) {
	if path, err := filepath.Rel(f.path, path.Join(dirPath, entry.Name())); err == nil {
		color.Green(path)
	}
}

func getPathAndSearch(args []string, reg bool) (string, string) {
	if reg {
		return args[1], args[2]
	}
	return args[0], args[1]
}
