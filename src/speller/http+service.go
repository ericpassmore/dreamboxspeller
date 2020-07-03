package speller

import (
  "fmt"
  "strconv"
  "log"
  "net/http"
  "unicode"
  "strings"
  "encoding/json"
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
type ResponseBody struct {
  HTTPCode int
  UserInput string
  Suggestions []string
  Repeating bool
  MissingVowels bool
  MixedCase bool
  NotFound bool
}

// spelling service
func spelling(w http.ResponseWriter, req *http.Request) {
  emptyResults := true
  responseCode := -1
  matches := []string{}

  // get params
  query, ok := req.URL.Query()["q"]
  //optionalOrder, isOrderImportant := req.URL.Query()["orderImportant"]

  // need that query param, return 400 if can't find it
  if !ok || len(query[0]) < 1 {
    log.Println("query param is missing from search spelling request")
    http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
    return
  }

  // normalize query
  normalizeSearchTerm := strings.ToLower(query[0])
  // do the search
  var words = Search(normalizeSearchTerm,ConsonantsNotInWord(normalizeSearchTerm))
  // transform the array of struct to array of string
  for idx := 0; idx < len(words); idx++ {
    matches = append(matches, words[idx].Raw)
  }

  // check if there are any possible matches
  // no matches then 404 status
  // found matches then 200 status
  if (len(words)>0) {
    emptyResults = false
    responseCode = http.StatusOK
  } else {
    emptyResults = true
    responseCode = http.StatusNotFound
  }

  // populate our response
  response := ResponseBody {
    HTTPCode: responseCode,
    UserInput: query[0],
    Suggestions: matches,
    Repeating: false,
    MissingVowels: false,
    MixedCase: isMixedCase(query[0]),
    NotFound: emptyResults,
  }

  // empty then 404
  // otherwise 200 and print JSON
  if (emptyResults) {
    http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
  } else {
    var jsonData []byte
    jsonData, err := json.Marshal(response)
    if err != nil {
      log.Println(err)
    }
    fmt.Fprintf(w,string(jsonData))
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
