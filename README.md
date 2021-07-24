# go-fd

A fast, concurrent file finder.   

# Usage

```sh
go get github.com/1-irdA/go-fd
```

```go
package main

import "github.com/1-irdA/go-fd/finder"

func main() {
	options := finder.Options{
		File:     true,
		Dir:      false,
		Regex:    false,
		Absolute: true,
		Bench:    false,
	}
	fd := finder.New("", "Memoire_GARROUSTE", options)
	fd.Find()
}
```

## Features

```sh
go run main.go -h
go run main.go -p=path -s=to-search [-f] [-d] [-a] [-b] [-r]
```

You can search a file or directory by names or with regex.    