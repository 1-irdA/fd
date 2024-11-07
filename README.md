# fd

A fast, concurrent files finder.   

## Commands

```sh
fd -h
fd [flags] [search]
```

## Examples

```sh
# search recursively
# exclude migrations folder 
# with extension 'sql'
fd -x migrations -e sql

## search recursively
# with extension 'go'
# in location parent directory
# entry who contains 'main'
fd -e go -l ../ main

# search recursively
# get search time
# count result
# exclude node_modules
# search with pattern
# in folder 'anothapp-services'
# search this pattern '.*Service.*$' 
fd -t -c -x node_modules -p -l .\anothapp-services .*Service.*$ 

# search recursively
# in folder anothapp-services
# exclude node_modules, models controllers services
# search with pattern "^*js"
# search file
fd -l ..\anothapp-services\ -x "node_modules models controllers services" -p -t f "^*js"
```