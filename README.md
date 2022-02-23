# fd

A fast, concurrent file finder.   

# Usage

```sh
go get github.com/1-irdA/fd
```

```go
package main

import "github.com/1-irdA/fd"

func main() {
	options := fd.Options{
		File:     true,
		Dir:      false,
		Regex:    false,
		Absolute: true,
		Bench:    true,
	}
	fd.New("D:/", "main.cpp", options).Find()
}
```

## Features

```sh
go run main.go -h
go run main.go -p=path -s=to-search [-f] [-d] [-a] [-b] [-r]
```

You can search a file or directory by names or with regex.    
