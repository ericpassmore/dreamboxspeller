package speller

import (
  "testing"
  "fmt"
  "os"
)

func setup() {
  // build index, a one time event
  Build("/Users/eric/RandomRepos/dreamboxspeller/wordsEn.txt")
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
  var words = Search("year's",ConsonantsNotInWord("year's"))
  hasYears := false
  for idx := 0; idx < len(words); idx++ {
    if words[idx].Raw == "year's" { hasYears = true; break; }
  }

  if !hasYears {
    t.Errorf("SearchApos failed to find year's")
  }
}

func TestSearchIncomplete(t *testing.T) {
  var words = Search("years",ConsonantsNotInWord("years"))
  hasYears := false
  for idx := 0; idx < len(words); idx++ {
    if words[idx].Raw == "year's" { hasYears = true; break }
  }

  if !hasYears {
    t.Errorf("SearchApos failed to find year's")
  }
}

func TestSearchFailed(t *testing.T) {
  query := []string{"fsiuyfiusyifys","dua",}

  for at := 0 ; at < len(query); at++ {
    var words = Search(query[at],ConsonantsNotInWord(query[at]))
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
  query := []string{"balloon","aah","ab","a"}

  for at := 0 ; at < len(query); at++ {
    fmt.Printf("word is: %s, not vowels is %s\n", query[at],ConsonantsNotInWord(query[at]))
    var words = Search(query[at],VowelsNotInWord(query[at]))
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
  var words = []string{ "yugoslavia", "viability", }
  var results = ""
  for i := 0; i < len(words); i++ {
    results = ConsonantsNotInWord(words[i])

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
