# dreamboxspeller
Project for Dreambox first go-land project ever.
## Description
create a web interface that takes a word and has a submit button
returns one of three responses
* Correct - found the word in the dictionary
* Not Found - did not find the work in the dictionary
* Incorrect - not found but close enough to offer suggested spellings

## Scripted Build
This is the recommended way, it is simpler
* cd into dreamboxspeller director
* `./build.sh`

## Manual Building
If you have problems with the scripted build & run
Need to have go installed. If you don't have it download from https://golang.org/dl/
* make sure go is in your PATH
* cd into dreamboxspeller director
* set import path `$ GOPATH=$(pwd); export GOPATH`
* build `$ go build dreamboxspeller.go`

## Running
Runs the program, and goto URL
* run `$ ./dreamboxspeller`
* goto 'http://localhost:8080/'
Check error log as needed
* `$ cat dreamboxspeller-error.log`
