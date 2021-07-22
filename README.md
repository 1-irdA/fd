# go-fd

A fast, concurrent file finder.   

## Features

```sh
go run main.go -h
go run main.go -p=dir -s=to-search [-r] [-a]
```

You can search a file by names or with regex.    

![regex](assets/reg.PNG)

## Benchmark

Find in all D:// drive :    

![drive](assets/drive.PNG)

Comparison was realised between my concurrent recursive find function and the [Walk function](https://golang.org/pkg/path/filepath/#Walk) from Golang.    

- Walk function

![walk](assets/walk.PNG)

- Concurrent find function

![find](assets/compare.PNG)