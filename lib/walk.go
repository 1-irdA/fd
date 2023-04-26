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
)

type config struct {
	search  string
	path    string
	routine uint
	pattern bool
	hidden  bool
	recurse bool
	count   bool
	occur   uint32
	state   state
}

type state struct {
	visit  chan string
	active sync.WaitGroup
}

type style struct {
	name, path           string
	nameColor, pathColor AnsiColor
}

func (s style) String() string {
	return fmt.Sprintf("%s%s%c%s%s", s.pathColor.String(), s.path, os.PathSeparator, s.nameColor.String(), s.name)
}

type Walker interface {
	Search()
}

// NewWalker init a walker
func NewWalker(args []string, routine uint, recurse, hidden, count, pattern bool) Walker {

	search, path := ".*", "."

	if len(args) == 1 {
		search = args[1]
	} else if len(args) == 2 {
		search, path = args[0], args[1]
		path, err := filepath.Abs(path)

		if err != nil {
			log.Fatalf("%s is not a valid path", path)
		}
	} else {
		pattern = true
	}

	if routine < 10 || routine > 100 {
		log.Fatal("invalide number of routine, -h for more informations")
	}
	if pattern {
		if _, err := regexp.Compile(search); err != nil {
			log.Fatalf("%s invalide pattern", path)
		}
	}
	return &config{
		search:  search,
		path:    path,
		routine: routine,
		pattern: pattern,
		recurse: recurse,
		hidden:  hidden,
		count:   count,
		state:   state{visit: make(chan string, 1024)},
	}
}

func (cg *config) check() {
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
		style := style{name: name, path: path, nameColor: Green, pathColor: SkyBlue}
		path := fmt.Sprintf("%s%c%s", path, os.PathSeparator, name)
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
			fmt.Println(style)

			if cg.count {
				atomic.AddUint32(&cg.occur, 1)
			}
		}
	}
}

// Search file from path following conditions
func (cg *config) Search() {
	_, err := os.Lstat(cg.path)

	if err != nil {
		fmt.Printf("%signore %s\n", Red.String(), cg.path)
		return
	}
	defer close(cg.state.visit)

	cg.state.active.Add(1)
	cg.state.visit <- cg.path

	for i := 0; i < int(cg.routine); i++ {
		go cg.check()
	}
	cg.state.active.Wait()

	if cg.count {
		fmt.Printf("%s%d files", Default.String(), cg.occur)
	}
}

func (cg *config) isRecursive(entry os.FileInfo) bool {
	return (cg.hidden && strings.HasPrefix(entry.Name(), ".")) || !strings.HasPrefix(entry.Name(), ".")
}

func (cg *config) isMatch(current os.FileInfo) bool {
	if cg.pattern {
		ok, _ := regexp.MatchString(cg.search, current.Name())
		return ok
	}
	return strings.Contains(current.Name(), cg.search)
}
