# fd

A fast, concurrent files finder.   

## Commands

```sh
fd -h
fd [flags] [search] [path]
```

## Examples

```sh
# search recursively
# exclude migrations folder 
# with extension sql
fd -r -x=migrations -e sql

# search recursively
# get search time
# count result
# exclude node_modules
# search with pattern '.*Service.*$' 
# in folder 'anothapp-services'
fd -r -t -c -x=node_modules -p .*Service.*$ .\anothapp-services
```