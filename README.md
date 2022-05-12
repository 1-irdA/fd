# fd

A fast, concurrent file finder.   

# Usage

```sh
go get github.com/1-irdA/fd
```

```go
package main

import (
	"github.com/1-irdA/fd"
)

func main() {
	config := fd.New("E:/Developpement", "^.*\\.sql$", true)
	fd.Fd(config).Find()
}
```

## Features

You can search a file or directory by names or with regex.    
