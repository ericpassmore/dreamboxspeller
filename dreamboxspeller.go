package main

import (
  "speller"
  "os"
  "log"
  "fmt"
)

const dictionary="wordsEn.txt"
const logFileName="dreamboxspeller-error.log"
const host="localhost"
const port=8080

func main() {
  // setup loging and close file at end
  var logFileHandle *os.File = setupLogging()
  defer logFileHandle.Close()

  // build index, a one time event
  speller.Build(getWorkingDirectory() + "/" + dictionary)
  var matched []speller.Word = speller.Search("a")
  for at:=0; at<len(matched); at++ {
    fmt.Println(matched[at].Raw)
  }
}
//############# UTIL FUNCTIONS ################
// get working directory so we can open dictionary file
func getWorkingDirectory() string {
  path, err := os.Getwd()
  if err != nil {
    log.Println("Error getting working directory", err)
    // attempt to recover
    path = "../../"
  }
  return path
}

// setup loggin file
func setupLogging() *os.File {
  logFileHandle, err := os.Create(logFileName)
  if err != nil {
    log.Fatal(err)
  }
  log.SetOutput(logFileHandle)
  return logFileHandle
}
