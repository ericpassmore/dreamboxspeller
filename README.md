# dreamboxspeller
Project for Dreambox first go-land project ever.
## Description
create a web interface that takes a word and has a submit button
returns one of three responses
* Correct - found the word in the dictionary
* Not Found - did not find the work in the dictionary
* Incorrect - not found but close enough to offer suggested spellings

## High Level Design
Running `dreamboxspeller` creates a bitmapped index and a small data-structure for each word in the dictionary file, and starts a webservice. The webservice support an HTML interface and a JSON-over-HTTP api. The spelling service does three things
* Searches the dictionary (indexes+search.go) using the bitmapped indexes, returning an array of possible matches
* Filters (http+services.go) the search results down to a set of suggestions.
* Returns the suggestions encoded as a JSON object, or returns an HTTP error
Navigate to "http://localhost:8080/about.html" and "http://localhost:8080/userguide.html" for additional usage instruction and HTTP API details.

Note: both indexes+search.go and http+services.go utilize the functions in vowels+consonants.go


## Scripted Build
This is the recommended way, it is simpler
* cd into dreamboxspeller director
* `./build.sh`

## Manual Building
If you have problems with the scripted build try these steps
Need to have go installed. If you don't have it download from https://golang.org/dl/
* make sure go is in your PATH
* cd into dreamboxspeller director
* set import path `$ GOPATH=$(pwd); export GOPATH`
* build `$ go build dreamboxspeller.go`

## Running
Runs the program, and goto URL
* run `$ ./dreamboxspeller &`
* goto 'http://localhost:8080/'
Check error log as needed
* `$ cat dreamboxspeller-error.log`

## Testing
Runs test for package. No tests for main.
* cd into dreamboxspeller director
* Tests run are both unit tests and integration tests. The integration tests require the service to be up answering on port 8080. This is my first time using go-lang, going forward it would be best to separate out the unit tests from the integration test.
* set import path `$ GOPATH=$(pwd); export GOPATH`
* run service `$ ./dreamboxspeller &`
* `$ go test speller`
* shutdown service `$ kill pid`

## Files and Directors
Under the main directory you will find the following
* src/speller - go files, the program executables
  * http+service.go - the web service supporting /speller and /hello paths
  * indexes+search.go - builds the search index and performs the search
  * vowels+consonants.go - functions that perform operations on strings and characters, to determine if they are consonants or vowels, or to extract the consonants or vowels found in the strings
  * package_test.go - performs unit tests and integration tests via http GET
* web - html and javascript supporting the user interface
  * about.html - html page documenting the HTTP service calls, along with JSON response
  * colortheme.js - logic to change from dark theme to light theme
  * favicon-32.png - tiny picture in your tab of the web browser
  * index.html - the user interface, submit a word and get spelling suggestions
  * jquery-3.5.1.min.js - javascript library to enable DOM filtering and manipulation
  * main.css - style sheet for web pages
  * sliders.css - styling for dark mode toggle switch
  * userguide.html - explains how to run the user interface and the options
* build.sh - convenience function to compile the go program, also inserts the build date into /hello
* dreamboxspeller.go - main loop of program, it only does two things, builds the index and starts the webservice
* LICENSE - anyone can use it
* README.md - this file
* wordsEn.txt - our dictionary of words, this file is indexed to support our spelling program
* dreamboxspeller-error.log - This file is created after the program is running. Contains error logs and debug statments
