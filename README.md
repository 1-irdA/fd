# go-fd

A fast, concurrent file finder.   

## Features

```sh
go run go-fd -h
go run go-fd dir to-search
go run go-fd -p dir to-search
```

You can search a file by names or with regex.    

![alt text](assets/regex.PNG)

## Benchmark

Comparison was realised between my concurrent find function and the [Walk function](https://golang.org/pkg/path/filepath/#Walk) from Golang.    

- Walk function

![alt text](assets/walk.PNG)

- Concurrent find function

![alt text](assets/concurrency.PNG)