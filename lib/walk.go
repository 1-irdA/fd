package lib

import (
	"fmt"
	"log"
	"os"
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
	routine   uint
	pattern   bool
	hidden    bool
	recurse   bool
	count     bool
	extension bool
	time      bool
	occur     uint32
	state     state
}

type state struct {
	visit  chan string
	active sync.WaitGroup
}

type print struct {
	name, path string
}

func (p print) Print() {
	fmt.Println(fmt.Sprintf("%s%s%c%s%s", SkyBlue, p.path, os.PathSeparator, Green, p.name))
}

type Walker interface {
	Search()
}

// NewWalker init a walker
func NewWalker(args []string, routine uint, recurse, hidden, count, pattern, extension, time bool) Walker {

	search, path := checkArgs(args)

	config := &config{
		search:    search,
		path:      path,
		routine:   routine,
		pattern:   pattern,
		recurse:   recurse,
		hidden:    hidden,
		count:     count,
		extension: extension,
		time:      time,
		state:     state{visit: make(chan string, 1024)},
	}
	return config.checkConfig()
}

func (cg *config) checkConfig() *config {

	if cg.routine < 10 || cg.routine > 100 {
		log.Fatal("invalide number of routine, -h for more informations")
	}
	if cg.pattern {
		if _, err := regexp.Compile(cg.search); err != nil {
			log.Fatalf("%s invalide pattern", cg.search)
		}
	}
	if cg.extension && !strings.HasPrefix(cg.search, ".") {
		cg.search = "." + cg.search
	}
	return cg
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
		fmt.Printf("%signore %s\n", Red.String(), path)
		return
	}
	defer f.Close()

	if info, _ := f.Stat(); !info.IsDir() {
		return
	}
	names, err := f.Readdirnames(-1)

	if err != nil {
		return
	}
	for _, name := range names {
		print := print{name: name, path: path}
		path := filepath.Join(path, name)
		info, err := os.Lstat(path)

		if err != nil {
			fmt.Printf("%signore %s\n", Red.String(), path)
			return
		}
		if info.IsDir() && cg.recurse && cg.isRecursive(info) {
			cg.state.active.Add(1)
			select {
			case cg.state.visit <- path:
			default:
				cg.state.active.Done()
				cg.checkFile(path)
			}
		} else if cg.isMatch(info) {
			print.Print()

			if cg.count {
				atomic.AddUint32(&cg.occur, 1)
			}
		}
	}
}

func (cg *config) launch() {
	_, err := os.Lstat(cg.path)

	if err != nil {
		fmt.Printf("%signore %s\n", Red.String(), cg.path)
		return
	}
	defer close(cg.state.visit)

	cg.state.active.Add(1)
	cg.state.visit <- cg.path

	for i := 0; i < int(cg.routine); i++ {
		go cg.walk()
	}
	cg.state.active.Wait()
}

// Search file from path following conditions
func (cg *config) Search() {
	var start time.Time

	if cg.time {
		start = time.Now()
	}
	cg.launch()
	fmt.Print(Default.String())

	if cg.count {
		fmt.Printf("%d files\n", cg.occur)
	}
	if cg.time {
		fmt.Printf("%s", time.Since(start))
	}
}

func (cg *config) isRecursive(entry os.FileInfo) bool {
	return (cg.hidden && strings.HasPrefix(entry.Name(), ".")) || !strings.HasPrefix(entry.Name(), ".")
}

func (cg *config) isMatch(current os.FileInfo) bool {
	if cg.pattern {
		ok, _ := regexp.MatchString(cg.search, current.Name())
		return ok
	} else if cg.extension {
		return filepath.Ext(current.Name()) == cg.search
	}
	return strings.Contains(current.Name(), cg.search)
}
