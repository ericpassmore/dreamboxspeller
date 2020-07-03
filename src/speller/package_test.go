package speller

import (
  "testing"
  "fmt"
)

func TestSearchApos(t *testing.T) {
  // build index, a one time event
  Build("/Users/eric/RandomRepos/dreamboxspeller/wordsEn.txt")

  var words = Search("year's","iou")
  hasYears := false
  for idx := 0; idx < len(words); idx++ {
    fmt.Println("word is: ", words[idx].Raw )
    if words[idx].Raw == "year's" { hasYears = true }
  }

  if !hasYears {
    t.Errorf("SearchApos failed to find year's")
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
