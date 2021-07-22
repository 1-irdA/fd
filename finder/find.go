package finder

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

	"github.com/fatih/color"
)

type find struct {
	path     string
	searched string
	nbFiles  uint32
	reg      bool
	absolute bool
	wg       sync.WaitGroup
	elapsed  time.Duration
}

func New(path string, searched string, pattern bool, absolute bool) *find {
	_, err := os.Lstat(path)
	if err != nil {
		printErr("Invalid path")
	}
	return &find{path: path, searched: searched, reg: pattern, absolute: absolute}
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

	folder, _ := os.Open(dirPath)
	dirs, _ := folder.ReadDir(-1)

	for _, entry := range dirs {
		atomic.AddUint32(&f.nbFiles, 1)
		if entry.IsDir() {
			f.wg.Add(1)
			go f.worker(path.Join(dirPath, entry.Name()))
		} else if f.correspond(entry) {
			f.printPath(dirPath, entry)
		}
	}

	folder.Close()
}

func (f *find) details() {
	color.Yellow("Files browsed %d, search duration : %v\n", f.nbFiles, f.elapsed)
}

func (f *find) checkRegIfNeed() {
	if f.reg {
		if _, err := regexp.Compile(f.searched); err != nil {
			printErr("Invalid regex : " + f.searched)
		}
	}
}

func (f *find) correspond(entry fs.DirEntry) bool {
	var result bool = false
	if f.reg {
		match, err := regexp.MatchString(f.searched, entry.Name())
		if err != nil {
			printErr("Regex error")
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
	if f.absolute {
		color.Green(path.Join(dirPath, entry.Name()))
	} else if path, err := filepath.Rel(f.path, path.Join(dirPath, entry.Name())); err == nil {
		color.Green(path)
	}
}

func printErr(msg string) {
	color.Red(msg)
	os.Exit(1)
}
