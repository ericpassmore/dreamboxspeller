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
// tollerance used for internal score calulation, bigger values returns more results
const missingVowelTollerance = 1
const orderOffsetTollerance = 2

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

var debug = false

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
  // orderNotImportant ignore letter position during compare, returns more suggestions
  // first tracks the very first word returned from search
  // missingVowel tracks one reason matches fail
  // repeating tracks another reason matches fail, user input has too many letters
  // suggestions are possible spellings
  // words is the data struct returned from Search
  exactMatch := false
  relaxY := false
  orderNotImportant := false
  first := true
  allWordsHaveMissingVowels := false
  allWordsHaveRepeatingLetters := false
  suggestions := []string{}
  var words = []Word{}

  // get params pass in reference to http.Request
  // must have "q" param for the query example: params["q"]
  params, paramsOK := getParams(req)
  //optionalOrder, isOrderImportant := req.URL.Query()["orderImportant"]

  // need that query param, log error, return 400 if can't find it
  if paramsOK != nil {
    debug = false
    log.Println(paramsOK)
    http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
    return
  }

  // optional params, move from string to bool
  if (params["relaxy"] == "true") {
    relaxY = true
  }
  if (params["ordernotimportant"] == "true") {
    orderNotImportant = true
  }
  if (params["debug"] == "true" ) { debug = true }

  // convert query value to ordered collection of letters
  // note CreateLetterMap function found in index+search.go
  queryAsLetterMap, queryValueOK := CreateLetterMap(params["q"])

  // 400 response if nothing retured
  // return immediattly if bad user input
  if queryValueOK != nil {
    debug = false
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

    if (debug) { log.Printf("comparing userinput is %s dictionary word is %s ",lowerCaseQuery, words[idx].Raw) }

    // loop through our user input and dictionary word letter by letter
    // as we go track, numbers indicate returned paramater pos in func call
    // 1. count of missing vowels
    // 2. max distance between letter's position in word
    // 3. user input contiguous letters exceeding dictionary word
    isValidSuggestion, hasMissingVowel, hasRepeatingLetters := compareLetterMaps(queryAsLetterMap, words[idx].LetterMap, orderNotImportant )

    if (debug ) {
      if isValidSuggestion {
        log.Println("Is Valid Suggestion ")
      } else {
        log.Println("Is NOT Valid Suggestion")
      }
    }
    // too many missing vowels, cuttoff of two
    if (isValidSuggestion) {
      // user input has a letter with higher contiguous count
      // across all dictionary words
      // clearly a case of repeating letting in user input
      // example consider "balllooon"
      if (debug && hasMissingVowel) {log.Printf(" Has Missing Vowels ")}
      if (debug && hasRepeatingLetters) {log.Printf(" Has Repeating Letters ")}
      if debug { log.Println(" - end") }
      allWordsHaveRepeatingLetters = (allWordsHaveRepeatingLetters || first) && hasRepeatingLetters
      allWordsHaveMissingVowels = (allWordsHaveMissingVowels || first) && hasMissingVowel
      first = false

      suggestions = append(suggestions, words[idx].Raw)
    }

    if debug { log.Println("********************") }
  }

  // populate our response
  response := ResponseBody {
    ExactMatch: exactMatch,
    UserInput: params["q"],
    Suggestions: suggestions,
    Repeating: allWordsHaveRepeatingLetters,
    MissingVowels: allWordsHaveMissingVowels,
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

  debug = false

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

  // debug
  _ , debugOK := req.URL.Query()["debug"]
  if !debugOK {
    params["debug"] = "false"
  } else {
    params["debug"] = "true"
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
// when isOrderNotImportant is set letter maps are compared as a bag
func compareLetterMaps(userInput map[rune]letter, dictionaryWord map[rune]letter, orderNotImportant bool) (isValidSuggestion bool, hasMissingVowel bool, hasRepeatingLetters bool) {
  // default return values
  hasRepeatingLetters = false
  isValidSuggestion = false
  // next three are inputs to determin if it is a legit suggestion
  missingVowelCount := 0
  maxOrderOffset := 0
  orderOffsetFirstConsonant := 0
  // start with a big values and then track to smallest as we iterate
  minPositionFirstConsonant := len(userInput)

  // loop over letters in dictionary word
  for letter, dictionaryWordDetails := range dictionaryWord {
    // letter exist in user input?
    if userInput[letter].count > 0 {
      // calc positional difference in letters
      thisOffset := dictionaryWordDetails.position - userInput[letter].position
      // change to abs offset
      if (thisOffset < 0) { thisOffset = thisOffset * -1 }
      if ( thisOffset > maxOrderOffset ) { maxOrderOffset = thisOffset }
      // track positional difference in first consonant, this is weighted higher
      if ( !userInput[letter].isVowel && (minPositionFirstConsonant > dictionaryWordDetails.position) ) {
        minPositionFirstConsonant = dictionaryWordDetails.position
        orderOffsetFirstConsonant = thisOffset
      }
      // look for repeating letters, this is returned to user as a rational for no-match
      if (userInput[letter].count > dictionaryWordDetails.count && userInput[letter].count > 1 ) {
        hasRepeatingLetters = true
      }
    } else {
      if debug { log.Printf("missing letter %s \n",string(letter)) }
      // check if non existing letter is a vowel
      if dictionaryWordDetails.isVowel {
        if debug { log.Println(" .. And its a vowel") }
        missingVowelCount++
      }
    }
  }

  // ***** LAST THREE IF STMT CALC SCORE *******
  // missing vowels withing tollerance
  if (missingVowelCount <= missingVowelTollerance ) {
    isValidSuggestion = true
  } else {
    if debug { log.Printf("Is Missing Vowels") }
    isValidSuggestion = false
  }
  // positional compare within tollerance and we care about order
  if ( (maxOrderOffset <= orderOffsetTollerance) || orderNotImportant ) {
    // and with true to track across all conditions
    isValidSuggestion = isValidSuggestion && true
  } else {
    if debug { log.Printf("Max offset too high,  out of position") }
    isValidSuggestion = false
  }
  // first consonant within tollerance, tollerance is reasonable, and we care about order
  if ((orderOffsetFirstConsonant <= (orderOffsetTollerance -1) && (orderOffsetTollerance -1) >= 0) || orderNotImportant) {
    // and with true to track across all conditions
    isValidSuggestion = isValidSuggestion && true
  } else {
    if debug { log.Printf("First Consonant out of position") }
    isValidSuggestion = false
  }

  return isValidSuggestion, missingVowelCount>0, hasRepeatingLetters
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
