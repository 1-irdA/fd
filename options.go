package fd

// just a find option to search
// with regex or not
// file, dir
// display absolute or relative path
// display search time
type Options struct {
	File     bool
	Dir      bool
	Regex    bool
	Absolute bool
	Bench    bool
}
