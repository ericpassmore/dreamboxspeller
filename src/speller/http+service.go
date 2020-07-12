package speller

import (
  "fmt"
  "strconv"
  "log"
  "net/http"
  "unicode"
  "encoding/json"
  "errors"
  "strings"
)

// ############### CONSTANTS ##########################
const doctype="<!DOCTYPE html>"
const defaultContentType="application/json; charset=utf-8"

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
  ExactMatch bool
  UserInput string
  Suggestions []string
  Repeating bool
  MissingVowels bool
  MixedCase bool
}
// ###################### END PUBLIC INTERFACE ####################


// spelling service
// takes user input from url query param and looks up spelling suggestions
// example http://host:port/spelling?q=myword
// returns 400 - bad user input or missing query
//         404 - no results no suggestions
//         200 - exact match or have suggestions
// 400 or 404 retunrs body with name human readable error name
// 200 returns JSON matcing ResponseBody struct
func spelling(w http.ResponseWriter, req *http.Request) {
  // set content type, applies to all functions
  w.Header().Set("Content-Type", defaultContentType)

  // exactMatch assume false
  // relaxY removed filter by 'Y' consonant , returns more suggestions
  // missingVowel tracks one reason matches fail
  // repeating tracks another reason matches fail, user input has too many letters
  // suggestions are possible spellings
  // words is the data struct returned from Search
  exactMatch := false
  relaxY := false
  missingVowel := false
  repeatingLetters := false
  suggestions := []string{}
  var words = []Word{}

  // get params pass in reference to http.Request
  // must have "q" param for the query example: params["q"]
  params, paramsOK := getParams(req)
  //optionalOrder, isOrderImportant := req.URL.Query()["orderImportant"]

  // need that query param, log error, return 400 if can't find it
  if paramsOK != nil {
    log.Println(paramsOK)
    http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
    return
  }

  // optional params, move from string to bool
  if (params["relaxy"] == "true") {
    relaxY = true
  }

  // convert query value to ordered collection of letters
  // note CreateLetterMap function found in index+search.go
  queryAsLetterMap, queryValueOK := CreateLetterMap(params["q"])

  // 400 response if nothing retured
  // return immediattly if bad user input
  if queryValueOK != nil {
    log.Println(queryValueOK)
    http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
    return
  }

  // pull out all the letters in the query
  mustHaveLetters := createString(queryAsLetterMap)
  // check if mixed case
  isMixedCase := isMixedCase(params["q"])
  // normalize to make matching easy
  lowerCaseQuery := strings.ToLower(params["q"])
  // do the search, ConsonantsNotInWord func in vowels+consonants.go
  // Search(mustHave, mustNotHave)
  // relaxY enables/disables filtering by 'y' consonants
  // set relaxY true to get more results
  // example of relaxY searching for 'ellow' suggests 'yellow'
  words = Search(mustHaveLetters,ConsonantsNotInWord(mustHaveLetters, relaxY))
  // loop through looking for exact match
  // when query is mixed case considered a misspelling, exact match not possible
  // if exact match end loop and clear out suggestions
  // otherwise build list of matching suggestions
  for idx := 0; idx < len(words); idx++ {
    if !isMixedCase && lowerCaseQuery == words[idx].Raw {
      exactMatch = true
      suggestions = []string{}
      // success, all done move to the exits
      break
    }

    // little debug
    if (words[idx].Raw == "balloon" || words[idx].Raw == "bicycle") {
      for l, details := range words[idx].LetterMap {
        log.Printf("Word %s letter %s position %d count %d",words[idx].Raw,string(l),details.position, details.count)
      }
    }


    // loop through our user input and dictionary word letter by letter
    // as we go track, numbers indicate returned paramater pos in func call
    // 1. count of missing vowels
    // 2. max distance between letter's position in word
    // 3. user input contiguous letters exceeding dictionary word
    missingVowelCount, orderOffset, repeating := compareLetterMaps(queryAsLetterMap, words[idx].LetterMap )
    // user input has a letter with higher contiguous count
    // across all dictionary words
    // clearly a case of repeating letting in user input
    // example consider "balllooon"
    repeatingLetters = repeating && repeatingLetters

    // too many missing vowels, cuttoff of two
    if (missingVowelCount < 2 && orderOffset < 3) {
      suggestions = append(suggestions, words[idx].Raw)
    }
  }

  // populate our response
  response := ResponseBody {
    ExactMatch: exactMatch,
    UserInput: params["q"],
    Suggestions: suggestions,
    Repeating: repeatingLetters,
    MissingVowels: missingVowel,
    MixedCase: isMixedCase,
  }

  // print JSON only when 200
  // 200 if exact match or we have suggestions
  // otherwise 404 not found
  if (exactMatch || len(suggestions)>0) {
    // convert structure to JSON
    var jsonData []byte
    jsonData, err := json.Marshal(response)
    if err != nil {
      log.Println(err)
    } else {
      fmt.Fprintf(w,string(jsonData))
    }
  } else {
    http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
  }

}

// simple handler can be used to see if service is up
func hello(w http.ResponseWriter, req *http.Request) {
  // set content type, applies to all functions
  w.Header().Set("Content-Type", "text/html; charset: utf-8")
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

  // optional params
  _ , relaxyOK := req.URL.Query()["relaxy"]
  if !relaxyOK {
    params["relaxy"] = "false"
  } else {
    params["relaxy"] = "true"
  }

  // optional params
  _ , optionalOrderOK := req.URL.Query()["ordernotimportant"]
  if !optionalOrderOK {
    params["ordernotimportant"] = "false"
  } else {
    params["ordernotimportant"] = "true"
  }

  // no errors
  return params, nil
}

// compare two letterMaps
// letterMap has position, count, and isVowel for each letter
// we use this information to compare two words and see if they
// are close enough to be considered a suggestion
// userInput with too many missing vowels and letters out of position
// example of out of position consider "spam" to "maps"
func compareLetterMaps(userInput map[rune]letter, dictionaryWord map[rune]letter) (missingVowelCount int, maxOrderOffset int, repeatingLetters bool) {
  missingVowelCount = 0
  maxOrderOffset = 0
  repeatingLetters = false

  for letter, dictionaryWordDetails := range dictionaryWord {
    if userInput[letter].count > 0 {
      thisOffset := dictionaryWordDetails.position - userInput[letter].position
      // change to abs offset
      if (thisOffset < 0) { thisOffset = thisOffset * -1 }
      if ( thisOffset > maxOrderOffset ) { maxOrderOffset = thisOffset }
      if (userInput[letter].count > dictionaryWordDetails.count) {
        repeatingLetters = true
      }
    } else {
      // check if non existing letter is a vowel
      if dictionaryWordDetails.isVowel {
        missingVowelCount++
      }
    }
  }
  return missingVowelCount, maxOrderOffset, repeatingLetters
}

// build a bag of letters from map kyes
// order of letters in map is not guarenteed
func createString(letterMap map[rune]letter) string {
  letterBag := ""
  // key is rune
  for key, _ := range letterMap {
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
