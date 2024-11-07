package lib

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type config struct {
	search    string
	path      string
	pattern   bool
	hidden    bool
	recurse   bool
	count     bool
	extension string
	bench     bool
	occur     uint32
	exclude   []string
	entryType string
	location  string
	state     state
}

type state struct {
	visit  chan string
	active sync.WaitGroup
}

type print struct {
	name, path string
}

func (p print) print(info os.FileInfo) {
	if info.IsDir() {
		fmt.Println(fmt.Sprintf("%s%s%c%s%s", SkyBlue, p.path, os.PathSeparator, SkyBlue, p.name))
	} else {
		fmt.Println(fmt.Sprintf("%s%s%c%s%s", SkyBlue, p.path, os.PathSeparator, Green, p.name))
	}
}

func handleStop() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Print(Default.String())
			os.Exit(0)
		}
	}()

}

func getSearch(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	return ""
}

type Walker interface {
	Search()
}

// NewWalker init a walker
func NewWalker(args []string, recurse, hidden, count, pattern, bench bool, extension, exclude, entryType, location string) Walker {
	config := &config{
		search:    getSearch(args),
		path:      location,
		pattern:   pattern,
		recurse:   recurse,
		hidden:    hidden,
		count:     count,
		extension: extension,
		bench:     bench,
		exclude:   strings.Split(exclude, " "),
		entryType: entryType,
		state:     state{visit: make(chan string, 1024)},
	}
	return config.checkConfig()
}

func (cg *config) checkConfig() *config {

	if cg.pattern {
		if _, err := regexp.Compile(cg.search); err != nil {
			log.Fatalf("%s invalide pattern", cg.search)
		}
	}
	if cg.isSearchExt() && !strings.HasPrefix(cg.extension, ".") {
		cg.extension = "." + cg.extension
	}
	if !cg.isSearchFile() && !cg.isSearchDir() {
		cg.entryType = ""
	}
	if path, err := filepath.Abs(cg.location); err != nil {
		log.Fatalf("%s is not a valid path", path)
	} else {
		cg.location = path
	}
	return cg
}

func (cg *config) isSearchExt() bool {
	return cg.extension != ""
}

func (cg *config) isSearchAll() bool {
	return cg.entryType == ""
}

func (cg *config) isSearchFile() bool {
	return cg.entryType == "f"
}

func (cg *config) isSearchDir() bool {
	return cg.entryType == "d"
}

func (cg *config) walk() {
	for path := range cg.state.visit {
		cg.checkFile(path)
		cg.state.active.Done()
	}
}

func (cg *config) checkFile(path string) {
	f, err := os.Open(path)

	if err != nil {
		// fmt.Printf("%signore %s\n", Red, path)
		return
	}
	defer f.Close()
	names, err := f.Readdirnames(-1)

	if err != nil {
		return
	}
	for _, name := range names {
		print := print{name: name, path: path}
		path := filepath.Join(path, name)
		info, err := os.Lstat(path)

		if err != nil {
			fmt.Printf("%signore %s\n", Red, path)
			return
		}

		if cg.isMatch(info) {
			print.print(info)

			if cg.count {
				atomic.AddUint32(&cg.occur, 1)
			}
		}
		if info.IsDir() && cg.isRecursive(info) && cg.isNotExclude(info) {
			cg.state.active.Add(1)
			select {
			case cg.state.visit <- path:
			default:
				cg.state.active.Done()
				cg.checkFile(path)
			}
		}
	}
}

func (cg *config) launch() {

	if _, err := os.Lstat(cg.path); err != nil {
		fmt.Printf("%signore %s\n", Red, cg.path)
		return
	}
	defer close(cg.state.visit)

	handleStop()
	cg.state.active.Add(1)
	cg.state.visit <- cg.path

	for i := 0; i < 50; i++ {
		go cg.walk()
	}
	cg.state.active.Wait()
}

// Search file from path following conditions
func (cg *config) Search() {
	var start time.Time

	if cg.bench {
		start = time.Now()
	}
	cg.launch()
	fmt.Print(Default.String())

	if cg.count {
		fmt.Printf("%d match(es)\n", cg.occur)
	}
	if cg.bench {
		fmt.Printf("%s", time.Since(start))
	}
}

func (cg *config) isNotExclude(entry os.FileInfo) bool {
	for _, ex := range cg.exclude {
		if ex == entry.Name() {
			return false
		}
	}
	return true
}

func (cg *config) isRecursive(entry os.FileInfo) bool {
	return cg.recurse && ((cg.hidden && strings.HasPrefix(entry.Name(), ".")) || !strings.HasPrefix(entry.Name(), "."))
}

func (cg *config) isMatch(entry os.FileInfo) bool {
	if !cg.isNotExclude(entry) ||
		(cg.isSearchFile() && entry.IsDir()) ||
		(cg.isSearchDir() && !entry.IsDir()) ||
		(cg.isSearchExt() && entry.IsDir()) {
		return false
	}
	match := true
	if cg.isSearchExt() {
		match = match && strings.EqualFold(filepath.Ext(entry.Name()), cg.extension)
	}
	if cg.pattern {
		ok, _ := regexp.MatchString(cg.search, entry.Name())
		return match && ok
	}
	return match && strings.Contains(entry.Name(), cg.search) &&
		(cg.isSearchAll() || ((cg.isSearchFile() && !entry.IsDir()) || (cg.isSearchDir() && entry.IsDir())))
}
