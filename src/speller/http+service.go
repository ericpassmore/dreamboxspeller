package speller

import (
  "fmt"
  "strconv"
  "log"
  "net/http"
  "unicode"
  "encoding/json"
  "errors"
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
// returned as JSON in body on 200
// non 200 returns no body
type ResponseBody struct {
  HTTPCode int
  UserInput string
  Suggestions []string
  Repeating bool
  MissingVowels bool
  MixedCase bool
  NotFound bool
}
// ###################### END PUBLIC INTERFACE ####################


// spelling service
// gets the params from the URL
// uses the user input to do a search
// then looks for matches
// finally returns results
func spelling(w http.ResponseWriter, req *http.Request) {

  // tracks state first two booleans look at results
  // http status defaults to bad, updated later
  // matches are the raw suggestions
  // words is the data struct returned from Search
  emptyResults := true
  badUserInput := true
  responseCode := http.StatusBadRequest
  matches := []string{}
  var words = []Word{}

  // get params pass in reference to http.Request
  // must have "q" param for the query example: params["q"]
  params, paramsOK := getParams(req)
  //optionalOrder, isOrderImportant := req.URL.Query()["orderImportant"]

  // need that query param, log error, return 400 if can't find it
  if paramsOK != nil {
    log.Println(paramsOK)
    responseCode = http.StatusBadRequest
    http.Error(w, http.StatusText(responseCode), responseCode)
    return
  } else {
    badUserInput = false
  }


  // convert query value to ordered collection of letters
  // return error 404 if query value does not have any valid charaters
  // on error log message
  // note CreateLetterMap function found in index+search.go
  queryAsLetterMap, queryValueOK := CreateLetterMap(params["q"])

  // 400 response if nothing retured
  // return immediattly if bad user input
  if queryValueOK != nil {
    log.Println(queryValueOK)
    responseCode = http.StatusBadRequest
    http.Error(w, http.StatusText(responseCode), responseCode)
    return
  } else {
    badUserInput = false
  }


  // normalize query
  normalizeSearchTerm := stringFromMap(queryAsLetterMap)
  // do the search, ConsonantsNotInWord func in vowels+consonants.go
  words = Search(normalizeSearchTerm,ConsonantsNotInWord(normalizeSearchTerm))
  // transform the array of struct to array of string
  for idx := 0; idx < len(words); idx++ {
    matches = append(matches, words[idx].Raw)
  }


  // check if there are any possible matches
  // no matches then 404 status
  // found matches then 200 status
  // remember emptyResults defaults true
  if (len(words)>0) {
    emptyResults = false
    responseCode = http.StatusOK
  }

  // populate our response
  response := ResponseBody {
    HTTPCode: responseCode,
    UserInput: params["q"],
    Suggestions: matches,
    Repeating: false,
    MissingVowels: false,
    MixedCase: !badUserInput && isMixedCase(params["q"]),
    NotFound: emptyResults,
  }

  // convert structure to JSON
  var jsonData []byte
  jsonData, err := json.Marshal(response)
  if err != nil {
    log.Println(err)
  } else {
    fmt.Fprintf(w,string(jsonData))
  }

  // print JSON in all of the following cases
  // empty then 404
  // badUserInput then 400
  // otherwise 200
  if (emptyResults || badUserInput ) {
    http.Error(w, http.StatusText(responseCode), responseCode)
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

// get query return error if query param does not exists
func getParams(req *http.Request) (map[string]string, error) {
  params := make(map[string]string)

  // get params from request
  values, ok := req.URL.Query()["q"]
  //optionalOrder, isOrderImportant := req.URL.Query()["orderImportant"]

  // return error if query param is no good
  // otherwise add query param to map
  if !ok || len(values[0]) < 1 {
    return params, errors.New("query param is missing from search spelling request")
  } else {
    params["q"] = values[0]
  }

  // no errors
  return params, nil
}

// build a bag of letters from map kyes
// order of letters in map is not guarenteed
func stringFromMap(wordMap map[rune]Letter) string {
  letterBag := ""
  // key is rune
  for key, _ := range wordMap {
    letterBag += string(key)
  }
  return letterBag
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
