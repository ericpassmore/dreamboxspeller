package speller

import (
  "testing"
  "fmt"
  "os"
  "net/http"
  "io/ioutil"
  "encoding/json"
)

func setup() {
  path, _ := os.Getwd()
  // build index, a one time event
  Build(path + "/../../wordsEn.txt")
}
func shutdown() {
  fmt.Println("Testing Shutdown***")
}

func TestMain(m *testing.M) {
    setup()
    code := m.Run()
    shutdown()
    os.Exit(code)
}

func TestSearchApos(t *testing.T) {
  relaxY := false
  var words = Search("year's",ConsonantsNotInWord("year's",relaxY))
  hasYears := false
  for idx := 0; idx < len(words); idx++ {
    if words[idx].Raw == "year's" { hasYears = true; break; }
  }

  if !hasYears {
    t.Errorf("SearchApos failed to find year's")
  }
}

func TestSearchIncomplete(t *testing.T) {
  relaxY := false
  var words = Search("years",ConsonantsNotInWord("years",relaxY))
  hasYears := false
  for idx := 0; idx < len(words); idx++ {
    if words[idx].Raw == "year's" { hasYears = true; break }
  }

  if !hasYears {
    t.Errorf("SearchApos failed to find year's")
  }
}

func TestSearchFailed(t *testing.T) {
  relaxY := false
  query := []string{"fsiuyfiusyifys","dua",}

  for at := 0 ; at < len(query); at++ {
    var words = Search(query[at],ConsonantsNotInWord(query[at],relaxY))
    hasQuery := false
    for idx := 0; idx < len(words); idx++ {
      if words[idx].Raw == query[at] { hasQuery = true; break }
    }

    if hasQuery {
      t.Errorf("SearchFailed found the unfindable %s",query[at])
    }
  }
}

func TestSearchSucceed(t *testing.T) {
  relaxY := false
  query := []string{"balloon","aah","ab","a"}

  for at := 0 ; at < len(query); at++ {
    fmt.Printf("word is: %s, not vowels is %s\n", query[at],ConsonantsNotInWord(query[at],relaxY))
    var words = Search(query[at],ConsonantsNotInWord(query[at],relaxY))
    hasQuery := false
    for idx := 0; idx < len(words); idx++ {
      //fmt.Println("word is: ", words[idx].Raw )
      if words[idx].Raw == query[at] { hasQuery = true; break }
    }

    if !hasQuery {
      for idx := 0; idx < len(words); idx++ {
        fmt.Println("word is: ", words[idx].Raw )
      }
      t.Errorf("SearchSucceed did not find %s",query[at])
    }
  }
}

func TestVowelsIgnorBadChars(t *testing.T) {
  var words = []string{ "no%way", "pound#", "Bug", }
  var results = ""
  for i := 0; i < len(words); i++ {
    results = Vowels(words[i])

    switch i {
    case 0:
      if (results != "oa") {
        t.Errorf("Vowels did not handle bad chars in %s: %s, want: %s.", words[i],results, "oa")
      }
    case 1:
      if (results != "ou") {
        t.Errorf("Vowels did not handle bad chars in %s: %s, want: %s.", words[i], results, "ou")
      }
    case 2:
      if (results != "u") {
        t.Errorf("Vowels did not handle bad chars in %s: %s, want: %s.", words[i], results, "u")
      }
    default:
    }

  }
}

func TestAlwaysVowels(t *testing.T) {
  var words = []string{ "winter", "summer", "spring", "fall", }
  var results = ""
  for i := 0; i < len(words); i++ {
    results = Vowels(words[i])

    switch i {
    case 0:
      if (results != "ie") {
        t.Errorf("AlwaysVowels did not work in %s: %s, want: %s.", words[i], results, "ie")
      }
    case 1:
      if (results != "ue") {
        t.Errorf("AlwaysVowels did not work in %s: %s, want: %s.", words[i], results, "ue")
      }
    case 2:
      if (results != "i") {
        t.Errorf("AlwaysVowels did not work in %s: %s, want: %s.", words[i], results, "i")
      }
    case 3:
      if (results != "a") {
        t.Errorf("AlwaysVowels did not work in %s: %s, want: %s.", words[i], results, "a")
      }
    default:
    }

  }
}

func TestYVowels(t *testing.T) {
  var words = []string{ "candy", "bicycle", "gym", "year", }
  var results = ""
  for i := 0; i < len(words); i++ {
    results = Vowels(words[i])

    switch i {
    case 0:
      if (results != "ay") {
        t.Errorf("YVowels did not work in %s: %s, want: %s.", words[i], results, "ay")
      }
    case 1:
      if (results != "iye") {
        t.Errorf("YVowels did not work in %s: %s, want: %s.", words[i], results, "iye")
      }
    case 2:
      if (results != "y") {
        t.Errorf("YVowels did not work in %s: %s, want: %s.", words[i], results, "y")
      }
    case 3:
      if (results != "ea") {
        t.Errorf("YVowels did not work is %s: %s, want: %s.", words[i], results, "ea")
      }
    default:
    }
  }
}

func TestDistinctVowels(t *testing.T) {
  var words = []string{ "yugoslavia", "viability", }
  var results = ""
  for i := 0; i < len(words); i++ {
    results = Vowels(words[i])

    switch i {
    case 0:
      if (results != "uoai") {
        t.Errorf("DistinctVowels did not work in %s: %s, want: %s.", words[i], results, "uoai")
      }
    case 1:
      if (results != "iay") {
        t.Errorf("DistinctVowels did not work in %s: %s, want: %s.", words[i], results, "iay")
      }
    default:
    }
  }
}

func TestInverseVowels(t *testing.T) {
  var words = []string{ "yugoslavia", "viability", }
  var results = ""
  for i := 0; i < len(words); i++ {
    results = VowelsNotInWord(words[i])

    switch i {
    case 0:
      if (results != "e") {
        t.Errorf("InverseVowels did not work in %s: %s, want: %s.", words[i], results, "e")
      }
    case 1:
      if (results != "eou") {
        t.Errorf("InverseVowels did not work in %s: %s, want: %s.", words[i], results, "eou")
      }
    default:
    }
  }
}

func TestInverseConsonants(t *testing.T) {
  relaxY := false
  var words = []string{ "yugoslavia", "viability", }
  var results = ""
  for i := 0; i < len(words); i++ {
    results = ConsonantsNotInWord(words[i],relaxY)

    switch i {
    case 0:
      if (len(results) != len("bcdfhjkmnpqrtwxz")) {
        t.Errorf("InverseConsonants did not work in %s: %s, want: %s.", words[i], results, "bcdfhjkmnpqrtwxz")
      }
    case 1:
      if (len(results) != len("cdfghjkmnpqrswxz")) {
        t.Errorf("InverseConsonants did not work in %s: %s, want: %s.", words[i], results, "cdfghjkmnpqrswxz")
      }
    default:
    }
  }
}

func TestServiceUp(t *testing.T) {
  url := "http://localhost:8080/hello"
  resp, err := http.Get(url)
  if err != nil {
    t.Errorf("is service up? %s failed with %s", url, err)
  }

  defer resp.Body.Close()

  if resp.StatusCode == http.StatusOK {
    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      t.Errorf("error reading http response %s",err)
    }
    bodyString := string(bodyBytes)
    if len(bodyString) <= 10 {
      t.Errorf("is service up? response too short from %s", url)
    }
  } else {
    t.Errorf("non 2xx code from %s",url)
  }
}

func TestExactMatch(t *testing.T) {
  url := "http://localhost:8080/spelling?q=year%27s"
  resp, err := http.Get(url)
  if err != nil {
    t.Errorf("is service up? %s failed", url)
  }

  defer resp.Body.Close()

  var response ResponseBody

  if resp.StatusCode == http.StatusOK {
    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      t.Errorf("error reading http response %s",err)
    }
    json.Unmarshal(bodyBytes, &response)

    if (!response.ExactMatch) {
      t.Errorf("Failed Exact Match for url %s", url)
    }

  } else {
    t.Errorf("non 2xx code from %s",url)
  }
}

func TestNotExactMatch(t *testing.T) {
  url := "http://localhost:8080/spelling?q=yer"
  resp, err := http.Get(url)
  if err != nil {
    t.Errorf("is service up? %s failed", url)
  }

  defer resp.Body.Close()

  var response ResponseBody

  if resp.StatusCode == http.StatusOK {
    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      t.Errorf("error reading http response %s",err)
    }
    json.Unmarshal(bodyBytes, &response)

    if (response.ExactMatch) {
      t.Errorf("Failed Did Not Expect Exact Match for url %s", url)
    }

  } else {
    t.Errorf("non 2xx code from %s",url)
  }
}

func TestRepeatingLetter(t *testing.T) {
  url := "http://localhost:8080/spelling?q=balllooon"
  resp, err := http.Get(url)
  if err != nil {
    t.Errorf("is service up? %s failed", url)
  }

  defer resp.Body.Close()

  var response ResponseBody

  if resp.StatusCode == http.StatusOK {
    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      t.Errorf("error reading http response %s",err)
    }
    json.Unmarshal(bodyBytes, &response)

    if (!response.Repeating) {
      t.Errorf("Expected Repeating Found Non-Repeating for url %s", url)
    }

  } else {
    t.Errorf("non 2xx code from %s",url)
  }
}

func TestNonRepeatingLetter(t *testing.T) {
  url := "http://localhost:8080/spelling?q=balon"
  resp, err := http.Get(url)
  if err != nil {
    t.Errorf("is service up? %s failed", url)
  }

  defer resp.Body.Close()

  var response ResponseBody

  if resp.StatusCode == http.StatusOK {
    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      t.Errorf("error reading http response %s",err)
    }
    json.Unmarshal(bodyBytes, &response)

    if (response.Repeating) {
      t.Errorf("Expected Non-Repeating Found Repeating for url %s", url)
    }

  } else {
    t.Errorf("non 2xx code from %s",url)
  }
}

func TestNonMixedCase(t *testing.T) {
  url := "http://localhost:8080/spelling?q=balloon"
  resp, err := http.Get(url)
  if err != nil {
    t.Errorf("is service up? %s failed", url)
  }

  defer resp.Body.Close()

  var response ResponseBody

  if resp.StatusCode == http.StatusOK {
    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      t.Errorf("error reading http response %s",err)
    }
    json.Unmarshal(bodyBytes, &response)

    if (response.MixedCase) {
      t.Errorf("Expected Non-MixedCase and found otherwise for url %s", url)
    }

  } else {
    t.Errorf("non 2xx code from %s",url)
  }
}

func TestMissingVowels(t *testing.T) {
  url := "http://localhost:8080/spelling?q=pp"
  resp, err := http.Get(url)
  if err != nil {
    t.Errorf("is service up? %s failed", url)
  }

  defer resp.Body.Close()

  var response ResponseBody

  if resp.StatusCode == http.StatusOK {
    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      t.Errorf("error reading http response %s",err)
    }
    json.Unmarshal(bodyBytes, &response)

    if (!response.MissingVowels) {
      t.Errorf("Expected MissingVowels and found otherwise for url %s", url)
    }

  } else {
    t.Errorf("non 2xx code from %s",url)
  }
}

func TestRelaxY(t *testing.T) {
  relaxYCount := 0
  hardYCount := 0

  url := "http://localhost:8080/spelling?q=ast&relaxy"

  resp, err := http.Get(url)
  if err != nil {
    t.Errorf("is service up? %s failed", url)
  }

  defer resp.Body.Close()

  var response ResponseBody

  if resp.StatusCode == http.StatusOK {
    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      t.Errorf("error reading http response %s",err)
    }
    json.Unmarshal(bodyBytes, &response)

    relaxYCount = len(response.Suggestions)

  } else {
    t.Errorf("non 2xx code from %s",url)
  }

  url = "http://localhost:8080/spelling?q=ast"

  resp, err = http.Get(url)
  if err != nil {
    t.Errorf("is service up? %s failed", url)
  }

  defer resp.Body.Close()

  if resp.StatusCode == http.StatusOK {
    bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      t.Errorf("error reading http response %s",err)
    }
    json.Unmarshal(bodyBytes, &response)

    hardYCount = len(response.Suggestions)

  } else {
    t.Errorf("non 2xx code from %s",url)
  }

  if relaxYCount <= hardYCount {
    t.Errorf("Expected relaxY to have more suggestions for url http://localhost:8080/spelling?q=ast&relaxy")
  }

}
