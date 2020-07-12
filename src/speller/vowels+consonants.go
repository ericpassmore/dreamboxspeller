package speller

import (
  "unicode"
)

const a = 'a'
const e = 'e'
const i = 'i'
const o = 'o'
const u = 'u'
const y = 'y'
const emptyRune = '\000'

// find distinct vowels in a word
// results returned in order they appear in word
func Vowels(word string) string {
  var vowels = ""
  var previous = emptyRune
  var next = emptyRune
  // needed for distinct
  // yes we could map, but this is easy to read
  hasA := false; hasE := false; hasI := false; hasO := false; hasU := false;
  hasY := false

  for idx := 0; idx < len(word); idx++ {
    // get current char
    char := []rune(word)[idx]

    // defense 1: normalize to lower
    // 2: skip to next if not a letter
    if (unicode.IsLetter(char)) {
      // 1
      char = unicode.ToLower(char)
    } else {
      // 2
      continue
    }

    // lookahead one
    if (idx+1 < len(word)) {
      next = []rune(word)[idx+1]
    } else {
      next = emptyRune
    }

    // determin vowel
    // has check add once , make distinct list
    switch char {
    case a:
      if !hasA {
        vowels += "a"; hasA = true
      }
    case e:
      if !hasE {
        vowels += "e"; hasE = true
      }
    case i:
      if !hasI {
        vowels += "i"; hasI = true
      }
    case o:
      if !hasO {
        vowels += "o"; hasO = true
      }
    case u:
      if !hasU {
        vowels += "u"; hasU = true
      }
    case y:
      if !hasY {
        isLast := (idx == len(word)-1)
        if IsYAVowle(previous,next, isLast) {
          vowels += "y"; hasY = true
        }
      }
    default:
      // no actions needed for non-vowels
    }

    // lookback
    previous = char
  }

  return vowels
}

// find distinct list of vowels not in word
// note y not considered a vowel in this function
func VowelsNotInWord(word string) string {
  var vowels = Vowels(word)
  var missingVowels = ""
  // loop through once set these true as vowels found
  // at end use these bools to collect inverse
  // yes we could map, but this is easy to read
  hasA := false; hasE := false; hasI := false; hasO := false; hasU := false;

  for idx := 0; idx < len(vowels); idx++ {
    // get current char
    char := []rune(vowels)[idx]

    // defense 1: normalize to lower
    // 2: skip to next if not a letter
    if (unicode.IsLetter(char)) {
      // 1
      char = unicode.ToLower(char)
    } else {
      // 2
      continue
    }

    // mark found
    switch char {
    case a:
      hasA = true
    case e:
      hasE = true
    case i:
      hasI = true
    case o:
      hasO = true
    case u:
      hasU = true
    default:
      // nothing to do
    }
  }

  // add what was not found
  if (!hasA) { missingVowels += "a" }
  if (!hasE) { missingVowels += "e" }
  if (!hasI) { missingVowels += "i" }
  if (!hasO) { missingVowels += "o" }
  if (!hasU) { missingVowels += "u" }

  return missingVowels
}

// looks at a char determin if is a vowel
func IsVowel(char rune, prev rune, next rune, isLast bool) bool {
  isVowel := false
  // mark found
  switch char {
  case a:
    isVowel = true
  case e:
    isVowel = true
  case i:
    isVowel = true
  case o:
    isVowel = true
  case u:
    isVowel = true
  case y:
  default:
    // nothing to do
  }

  return isVowel
}

// find distinct list of consonants not in word
func ConsonantsNotInWord(word string, relaxY bool) string {
  var missingConsonants = ""

  consonants := map[rune]bool {
    'b' : false, 'c' : false, 'd' : false, 'f' : false, 'g' : false,
    'h' : false, 'j' : false, 'k' : false, 'l' : false, 'm' : false,
    'n' : false, 'p' : false, 'q' : false, 'r' : false, 's' : false,
    't' : false, 'v' : false, 'w' : false, 'x' : false, 'y' : false,
    'z' : false,
  }

  for idx := 0; idx < len(word); idx++ {
    // get current char
    char := []rune(word)[idx]

    // defense 1: normalize to lower
    // 2: skip to next if not a letter
    if (unicode.IsLetter(char)) {
      // 1
      char = unicode.ToLower(char)
      } else {
        // 2
        continue
      }

      consonants[char] = true
    }

    for letter, exists := range consonants {
      if !exists {
        // mostly 'y' is a vowel and sometimes is a consonant
        // setting relaxY ignores the different uses of 'y'
        // check if this is 'y' only add 'y' if relaxY is false
        if (letter != 'y') {
          missingConsonants += string(letter)
        } else {
          if !relaxY {
            missingConsonants += string(letter)
          }
        }
      }
    }
    return missingConsonants

  }

func IsYAVowle(prev rune, next rune, isLast bool) bool {
  // check previous is not a vowel
  if (prev != emptyRune && prev != a && prev != e && prev != i && prev != o && prev != u) {
    // case 2: y is at the end of a word && previous is not vowel
    // example candy
    if ( isLast ) {
      return true
    }
    // case 3: y is in the middle of a syllable or end of a syllable
    // example bicycle
    // case 1: y is the only vowel
    // example gym
    if (next != emptyRune && next != a && next != e && next != i && next != o && next != u) {
      return true
    }
  }

  return false
}
