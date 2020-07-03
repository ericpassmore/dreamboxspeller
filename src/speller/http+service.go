package speller

import (
  "fmt"
  "strconv"
  "log"
  "net/http"
  "unicode"
)

// ############### CONSTANTS ##########################
const doctype="<!DOCTYPE html>"
const defaultContentType="Content-Type: application/json"
const defaultCharset="Charset=utf-8"

// ##################### PUBLIC INTERFACE #######################
func StartHTTP(port int) {
  // say hello
  http.HandleFunc("/hello",hello)
  // static file handler powers User Interface
  webDirectory := http.FileServer(http.Dir("./web"))
  http.Handle("/", webDirectory)
  // spelling service
  http.HandleFunc("/spelling",spelling)

  // start service at endpoint
  endpoint := ":" + strconv.Itoa(port)
  log.Println("connecting to endpoint ", endpoint)
  log.Fatal(http.ListenAndServe(endpoint,nil))
}
// ###################### END PUBLIC INTERFACE ####################
type response struct {
  responseCode int
  userInput string
  suggestions string
  repeating bool
  missingVowels bool
  mixedCase bool
  notFound bool
}

// spelling service
func spelling(w http.ResponseWriter, req *http.Request) {
  // get params
  query, ok := req.URL.Query()["q"]
  //optionalOrder, isOrderImportant := req.URL.Query()["orderImportant"]

  if !ok || len(query[0]) < 1 {
    log.Println("query param is missing from search spelling request")
    http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
    return
  }

  fmt.Fprintf(w,"%v\n",query[0])
  if isMixedCase(query[0]) {
    fmt.Fprintf(w,"failed mixed case\n")
  } else {
    fmt.Fprintf(w,"ok\n")
  }

  // not found error
  if (false) {
    http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
  }
}

// hello handler
func hello(w http.ResponseWriter, req *http.Request) {
  fmt.Fprintf(w,"%v\n",doctype)
  fmt.Fprintf(w,"<html>\n")
  fmt.Fprintf(w,"<body>\n")
  fmt.Fprintf(w,"<h1>My First Heading</h1>\n")
  fmt.Fprintf(w,"<p>Greetings</p>\n")
  fmt.Fprintf(w,"</body>\n")
  fmt.Fprintf(w,"</html>\n")
}

// return ture when all UPPER OR all lower OR Capitalized
// otherwise false
func isMixedCase(userInput string) bool {
  // assume good
  isUpper := true
  isLower := true
  isCapitalized := true

  for i := 0; i < len(userInput); i++ {
    char := []rune(userInput)[i]
    // not upper (aka lower), is letter
    if !unicode.IsUpper(char) && unicode.IsLetter(char) {
      isUpper = false
      // lower case to start the word
      if ( i == 0 ) { isCapitalized = false }
    }
    // not lower (aka upper), is letter
    if !unicode.IsLower(char) && unicode.IsLetter(char) {
      isLower = false
      // upper case in middle of the word
      if ( i > 0 ) { isCapitalized = false }
    }
    // stop if we have failed
    if !(isUpper || isLower || isCapitalized) { break }
    // otherwise keep going
  }

  return !(isUpper || isLower || isCapitalized)
}
