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
	Paths    []string
	opts     Options
	wg       sync.WaitGroup
	elapsed  time.Duration
}

func New(path string, searched string, options Options) *find {
	checkParamsAndOpts(searched, options)
	return &find{path: path, searched: searched, opts: options}
}

func (f *find) Find() {
	defer f.details()
	f.checkRegIfNeed()
	f.launch()
}

func checkParamsAndOpts(searched string, opts Options) {
	if searched == "" {
		printErr("Needs value to search")
	} else if !opts.Dir && !opts.File {
		printErr("Needs to search files, folders or both")
	}
}

func (f *find) launch() {
	start := time.Now()
	f.wg.Add(1)

	if f.path == "" {
		f.checkEverywhere()
	} else {
		_, err := os.Lstat(f.path)
		if err != nil {
			printErr("Invalid path")
		}
		f.worker(f.path)
	}

	f.wg.Wait()
	f.elapsed = time.Since(start)
}

func (f *find) checkEverywhere() {
	for _, drive := range getDrives() {
		f.worker(drive)
	}
}

func (f *find) worker(dirPath string) {
	defer f.wg.Done()

	folder, _ := os.Open(dirPath)
	dirs, _ := folder.ReadDir(-1)

	for _, entry := range dirs {
		atomic.AddUint32(&f.nbFiles, 1)
		if entry.IsDir() {
			if f.correspond(entry) {
				f.printPath(dirPath, entry, color.New(color.FgBlue))
				f.Paths = append(f.Paths, dirPath)
			}
			f.wg.Add(1)
			go f.worker(path.Join(dirPath, entry.Name()))
		} else if f.correspond(entry) {
			f.printPath(dirPath, entry, color.New(color.FgGreen))
			f.Paths = append(f.Paths, filepath.Join(dirPath, entry.Name()))
		}
	}

	folder.Close()
}

func (f *find) details() {
	if f.opts.Bench {
		color.Yellow("Files browsed %d, search duration : %v\n", f.nbFiles, f.elapsed)
	}
}

func (f *find) checkRegIfNeed() {
	if f.opts.Regex {
		if _, err := regexp.Compile(f.searched); err != nil {
			printErr("Invalid regex : " + f.searched)
		}
	}
}

func (f *find) correspond(entry fs.DirEntry) (result bool) {
	if f.opts.File && f.opts.Dir {
		result = f.checkName(entry)
	} else if f.opts.File && !entry.IsDir() {
		result = f.checkName(entry)
	} else if f.opts.Dir && entry.IsDir() {
		result = f.checkName(entry)
	}
	return
}

func (f *find) checkName(entry fs.DirEntry) (result bool) {
	if f.opts.Regex {
		match, err := regexp.MatchString(f.searched, entry.Name())
		if err != nil {
			printErr("Regex error")
		}
		result = match
	} else if strings.Contains(entry.Name(), f.searched) {
		result = true
	}
	return
}

func getDrives() (drives []string) {
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		f, err := os.Open(string(drive) + ":/")
		if err == nil {
			drives = append(drives, string(drive)+"://")
			f.Close()
		}
	}
	return
}

func (f *find) printPath(dirPath string, entry fs.DirEntry, c *color.Color) {
	if f.opts.Absolute {
		c.Println(path.Join(dirPath, entry.Name()))
	} else if path, err := filepath.Rel(f.path, path.Join(dirPath, entry.Name())); err == nil {
		c.Println(path)
	}
}

func printErr(msg string) {
	color.Red(msg)
	os.Exit(1)
}
