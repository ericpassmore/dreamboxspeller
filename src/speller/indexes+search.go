package speller

import ( "os"
  "bufio"
  "log"
  "unicode"
  "errors"
)

// ############### CONFIG ##########################
// bucketSize is chunk size of data (8,32,64)
const bucketSize = 8
// size is derived from number of words in file
// plus one to cover remainder
// if there were more words this would need to be bigger
const bitmapSize = (109583/bucketSize)+1
const emptyByte byte = 0
// ############### END CONFIG #######################

// ############### INTERNAL DATA STRUCTURES ##################
// This is what we use to look things up
var dictionary = initIndexes()

type indexes struct {
  words []Word
  a [bitmapSize]byte
  b [bitmapSize]byte
  c [bitmapSize]byte
  d [bitmapSize]byte
  e [bitmapSize]byte
  f [bitmapSize]byte
  g [bitmapSize]byte
  h [bitmapSize]byte
  i [bitmapSize]byte
  j [bitmapSize]byte
  k [bitmapSize]byte
  l [bitmapSize]byte
  m [bitmapSize]byte
  n [bitmapSize]byte
  o [bitmapSize]byte
  p [bitmapSize]byte
  q [bitmapSize]byte
  r [bitmapSize]byte
  s [bitmapSize]byte
  t [bitmapSize]byte
  u [bitmapSize]byte
  v [bitmapSize]byte
  w [bitmapSize]byte
  x [bitmapSize]byte
  y [bitmapSize]byte
  z [bitmapSize]byte
  apos [bitmapSize]byte
  yconsonant [bitmapSize]byte
}

// ##################### PUBLIC INTERFACE #######################
// normalized structure for all indexed words
// supports comparison
type Letter struct {
  Count int
  Position int
  isVowel bool
}

// Data Structure for words
// Raw value is self explainitory
// WordMap is bag of Letters see struct above
// Length is the number of chars in the words
type Word struct {
  Raw string
  LetterMap map[rune]Letter
  Length int
}


// build indexes
func Build(file string) {
  dictionary = buildIndex(file)
}

func Search(mustHave string, mustNotHave string) []Word {
  var matched []Word

  isEmpty := false
  var resultsByteArray [bitmapSize]byte = initByteArray(isEmpty)

  //lots of code, maybe there is a better way to do this
  // first go code, so forgivness please
  for midx:=0; midx < len(mustHave); midx++ {
    char := []rune(mustHave)[midx]
    switch char {
    case 'a':
      and(dictionary.a, resultsByteArray, &resultsByteArray)
    case 'b':
      and(dictionary.b, resultsByteArray, &resultsByteArray)
    case 'c':
      and(dictionary.c, resultsByteArray, &resultsByteArray)
    case 'd':
      and(dictionary.d, resultsByteArray, &resultsByteArray)
    case 'e':
      and(dictionary.e, resultsByteArray, &resultsByteArray)
    case 'f':
      and(dictionary.f, resultsByteArray, &resultsByteArray)
    case 'g':
      and(dictionary.g, resultsByteArray, &resultsByteArray)
    case 'h':
      and(dictionary.h, resultsByteArray, &resultsByteArray)
    case 'i':
      and(dictionary.i, resultsByteArray, &resultsByteArray)
    case 'j':
      and(dictionary.j, resultsByteArray, &resultsByteArray)
    case 'k':
      and(dictionary.k, resultsByteArray, &resultsByteArray)
    case 'l':
      and(dictionary.l, resultsByteArray, &resultsByteArray)
    case 'm':
      and(dictionary.m, resultsByteArray, &resultsByteArray)
    case 'n':
      and(dictionary.n, resultsByteArray, &resultsByteArray)
    case 'o':
      and(dictionary.o, resultsByteArray, &resultsByteArray)
    case 'p':
      and(dictionary.p, resultsByteArray, &resultsByteArray)
    case 'q':
      and(dictionary.q, resultsByteArray, &resultsByteArray)
    case 'r':
      and(dictionary.r, resultsByteArray, &resultsByteArray)
    case 's':
      and(dictionary.s, resultsByteArray, &resultsByteArray)
    case 't':
      and(dictionary.t, resultsByteArray, &resultsByteArray)
    case 'u':
      and(dictionary.u, resultsByteArray, &resultsByteArray)
    case 'v':
      and(dictionary.v, resultsByteArray, &resultsByteArray)
    case 'w':
      and(dictionary.w, resultsByteArray, &resultsByteArray)
    case 'x':
      and(dictionary.x, resultsByteArray, &resultsByteArray)
    case 'y':
      and(dictionary.y, resultsByteArray, &resultsByteArray)
    case 'z':
      and(dictionary.z, resultsByteArray, &resultsByteArray)
    case '\'':
      and(dictionary.apos, resultsByteArray, &resultsByteArray)
    default:
      log.Println("not expecting to search this rune: ", char)
    }
  }
  //lots of code, maybe there is a better way to do this
  // first go code, so forgivness please
  for midx:=0; midx < len(mustNotHave); midx++ {
    char := []rune(mustNotHave)[midx]
    switch char {
    case 'a':
      andNot(dictionary.a, resultsByteArray, &resultsByteArray)
    case 'b':
      andNot(dictionary.b, resultsByteArray, &resultsByteArray)
    case 'c':
      andNot(dictionary.c, resultsByteArray, &resultsByteArray)
    case 'd':
      andNot(dictionary.d, resultsByteArray, &resultsByteArray)
    case 'e':
      andNot(dictionary.e, resultsByteArray, &resultsByteArray)
    case 'f':
      andNot(dictionary.f, resultsByteArray, &resultsByteArray)
    case 'g':
      andNot(dictionary.g, resultsByteArray, &resultsByteArray)
    case 'h':
      andNot(dictionary.h, resultsByteArray, &resultsByteArray)
    case 'i':
      andNot(dictionary.i, resultsByteArray, &resultsByteArray)
    case 'j':
      andNot(dictionary.j, resultsByteArray, &resultsByteArray)
    case 'k':
      andNot(dictionary.k, resultsByteArray, &resultsByteArray)
    case 'l':
      andNot(dictionary.l, resultsByteArray, &resultsByteArray)
    case 'm':
      andNot(dictionary.m, resultsByteArray, &resultsByteArray)
    case 'n':
      andNot(dictionary.n, resultsByteArray, &resultsByteArray)
    case 'o':
      andNot(dictionary.o, resultsByteArray, &resultsByteArray)
    case 'p':
      andNot(dictionary.p, resultsByteArray, &resultsByteArray)
    case 'q':
      andNot(dictionary.q, resultsByteArray, &resultsByteArray)
    case 'r':
      andNot(dictionary.r, resultsByteArray, &resultsByteArray)
    case 's':
      andNot(dictionary.s, resultsByteArray, &resultsByteArray)
    case 't':
      andNot(dictionary.t, resultsByteArray, &resultsByteArray)
    case 'u':
      andNot(dictionary.u, resultsByteArray, &resultsByteArray)
    case 'v':
      andNot(dictionary.v, resultsByteArray, &resultsByteArray)
    case 'w':
      andNot(dictionary.w, resultsByteArray, &resultsByteArray)
    case 'x':
      andNot(dictionary.x, resultsByteArray, &resultsByteArray)
    case 'y':
      andNot(dictionary.yconsonant, resultsByteArray, &resultsByteArray)
    case 'z':
      andNot(dictionary.z, resultsByteArray, &resultsByteArray)
    case '\'':
      andNot(dictionary.apos, resultsByteArray, &resultsByteArray)
    default:
      log.Println("not expecting to search this rune: ", char)
    }
  }

  var residents = bytesToLocation(resultsByteArray)
  for at := 0; at < len(residents); at++ {
    matched = append(matched,dictionary.words[residents[at]])
  }
  return matched
}

// create map of letters
// this tracks the position and checks if it is a vowel
func CreateLetterMap(userInput string) (map[rune]Letter, error) {
  // tracks if there are any valid characters
  // if none it is an error
  hasAnyValidChars := false
  prevChar := emptyRune
  nextChar := emptyRune

  letterMap := make(map[rune]Letter)
  letterCount := 1

  // loop through each char
  for idx := 0; idx < len(userInput); idx++ {
    char := []rune(userInput)[idx]
    // get next char, if no last char set to emptyRune
    if (idx+1 < len(userInput)) {
      nextChar = []rune(userInput)[idx+1]
    } else {
      nextChar = emptyRune
    }
    // if there continuous charaters, increment count, and skip to next
    // no need to set prevChar since the chars are the same
    if (char == nextChar) {
      letterCount++
      continue;
    }

    // check valid char
    if (unicode.IsLetter(char) || char == '\'') {
      hasAnyValidChars = true
      // handle apostrophy
      if ( char == '\'' ) {
        letterMap[char] = Letter {
          Count: letterCount,
          Position: idx,
          isVowel: false,
        }
      // handle a-z letters
      } else {
        char = unicode.ToLower(char)
        isLast := (idx == len(userInput)-1)
        // isVowel needs all those params to detect y as vowel
        // IsVowel function in vowels+consonants.go
        letterMap[char] = Letter {
          Count: letterCount,
          Position: idx,
          isVowel: IsVowel(char, prevChar, nextChar, isLast),
        }
      }
    }
    // we can only get here when char != nextChar
    // so set prevChar and letterCount appropriately
    prevChar = char
    letterCount = 1
  }

  // our only error
  if !hasAnyValidChars {
    return letterMap, errors.New("No valid Characters in user input: "+userInput)
  }

  return letterMap, nil
}
// ###################### END PUBLIC INTERFACE ####################

// and together, joined is an argument to support recursion
func and(first [bitmapSize]byte, second [bitmapSize]byte, joined *[bitmapSize]byte) {
  idx := 0

  // explicit goto loop enable go to inline
  // inline offers a speedup
  loop:
  // note byte arrays are all the same length
  if idx < len(first) {
    joined[idx] = first[idx] & second[idx]
    idx++
    goto loop
  }
}

// and not, joined is an argument to support recursion
// NOTE: first is the must not have vector
func andNot(first [bitmapSize]byte, second [bitmapSize]byte, joined *[bitmapSize]byte) {
  idx := 0

  // explicit goto loop enable go to inline
  // inline offers a speedup
  loop:
  // note byte arrays are all the same length
  if idx < len(first) {
    joined[idx] = ^first[idx] & second[idx]
    idx++
    goto loop
  }
}

// make new array init with fill 0 if isEmpty 1 otherwise
func initByteArray(isEmpty bool) [bitmapSize]byte {
  newIndex := [bitmapSize]byte{}

  for i := 0; i < bitmapSize; i++ {
    if isEmpty {
      newIndex[i] = emptyByte
      } else {
        newIndex[i] = ^emptyByte
      }
    }
    return newIndex
  }

  // make new indexes fill with 0
  //lots of code, maybe there is a better way to do this
  // first go code, so forgivness please
  func initIndexes() indexes {
    empty := true
    newIndexes := indexes{
      a: initByteArray(empty),
      b: initByteArray(empty),
      c: initByteArray(empty),
      d: initByteArray(empty),
      e: initByteArray(empty),
      f: initByteArray(empty),
      g: initByteArray(empty),
      h: initByteArray(empty),
      i: initByteArray(empty),
      j: initByteArray(empty),
      k: initByteArray(empty),
      l: initByteArray(empty),
      m: initByteArray(empty),
      n: initByteArray(empty),
      o: initByteArray(empty),
      p: initByteArray(empty),
      q: initByteArray(empty),
      r: initByteArray(empty),
      s: initByteArray(empty),
      t: initByteArray(empty),
      u: initByteArray(empty),
      v: initByteArray(empty),
      w: initByteArray(empty),
      x: initByteArray(empty),
      y: initByteArray(empty),
      z: initByteArray(empty),
      apos: initByteArray(empty),
      yconsonant: initByteArray(empty),
    }
    return newIndexes
  }

  // turns index (int) location into position and shift for bit position
  // need both index by byte and shift for bits
  func byteIndex(index int) (int, int) {
    position := index / bucketSize
    // remainder reversed 7 -> 0 .. 0 -> 7
    shift := bucketSize - (index % bucketSize) -1
    return position, shift
  }

  // convert bits (0|1) into index (int) positions
  // each bit represents an index entry
  func bytesToLocation(index [bitmapSize]byte) []int {
    var residents []int
    for i := 0; i < bitmapSize; i++ {
      for l := bucketSize-1; l > 0; l-- {
        if index[i]&(1<<uint(l)) > 0 {
          // i is every chunk; a chunk is (8,32,64) depending on datatype used
          // l is the bit inside that chunk
          // matching positions inside the chunk we calc in facy if stmt above
          residents = append(residents,((bucketSize-1)-l)+(i*bucketSize))
        }
      }
    }
    return residents
  }

  // scans the dictionary building a-z bitmapped indexes
  //
  func buildIndex(file string) indexes {
    var theseIndexes = initIndexes()
    var length = 0
    // emptyRune defined in vowels+consonants.go
    var previous = emptyRune
    var next = emptyRune

    // open file
    fileHandle, err := os.Open(file)
    if err != nil {
      log.Println("Error opening dictionary: ",file, err)
      os.Exit(1)
    }
    fileScanner := bufio.NewScanner(fileHandle)
    lineCount := 0

    // loop over file
    for fileScanner.Scan() {
      word := fileScanner.Text()
      // calc length once
      length = len(word)
      // create map of words
      // if error log error, and skip to next word
      letterMap, error := CreateLetterMap(word)
      if error != nil {
        log.Println(error)
        // add one we start at zero
        log.Printf("error occured on line %d while indexing dict\n", lineCount+1)
        continue
      }
      // add our word to index
      theseIndexes.words = append(theseIndexes.words,
        Word{Raw:word, LetterMap: letterMap, Length: length})
      // init when indexing new word
      previous = emptyRune
      next = emptyRune

      // loop over each rune and add to indexes
      for i := 0; i < length; i++ {
        char := []rune(word)[i]

        // lookahead one, needed for is y a vowel
        if (i+1 < length) {
          next = []rune(word)[i+1]
          } else {
            next = emptyRune
          }

          // used to do bit manipulation
          index, shift := byteIndex(lineCount)
          // lots of code , and easy to read
          // better then compact, obtuse code #pickyourpoison
          // index for each a-z char + apostrophy
          switch char {
          case 'a':
            theseIndexes.a[index] |= 1 << uint(shift)
          case 'b':
            theseIndexes.b[index] |= 1 << uint(shift)
          case 'c':
            theseIndexes.c[index] |= 1 << uint(shift)
          case 'd':
            theseIndexes.d[index] |= 1 << uint(shift)
          case 'e':
            theseIndexes.e[index] |= 1 << uint(shift)
          case 'f':
            theseIndexes.f[index] |= 1 << uint(shift)
          case 'g':
            theseIndexes.g[index] |= 1 << uint(shift)
          case 'h':
            theseIndexes.h[index] |= 1 << uint(shift)
          case 'i':
            theseIndexes.i[index] |= 1 << uint(shift)
          case 'j':
            theseIndexes.j[index] |= 1 << uint(shift)
          case 'k':
            theseIndexes.k[index] |= 1 << uint(shift)
          case 'l':
            theseIndexes.l[index] |= 1 << uint(shift)
          case 'm':
            theseIndexes.m[index] |= 1 << uint(shift)
          case 'n':
            theseIndexes.n[index] |= 1 << uint(shift)
          case 'o':
            theseIndexes.o[index] |= 1 << uint(shift)
          case 'p':
            theseIndexes.p[index] |= 1 << uint(shift)
          case 'q':
            theseIndexes.q[index] |= 1 << uint(shift)
          case 'r':
            theseIndexes.r[index] |= 1 << uint(shift)
          case 's':
            theseIndexes.s[index] |= 1 << uint(shift)
          case 't':
            theseIndexes.t[index] |= 1 << uint(shift)
          case 'u':
            theseIndexes.u[index] |= 1 << uint(shift)
          case 'v':
            theseIndexes.v[index] |= 1 << uint(shift)
          case 'w':
            theseIndexes.w[index] |= 1 << uint(shift)
          case 'x':
            theseIndexes.x[index] |= 1 << uint(shift)
          case 'y':
            theseIndexes.y[index] |= 1 << uint(shift)
            isLast := (i == length-1)
            if !IsYAVowle(previous, next, isLast) {
              theseIndexes.yconsonant[index] |= 1 << uint(shift)
            }
          case 'z':
            theseIndexes.z[index] |= 1 << uint(shift)
          case '\'':
            theseIndexes.apos[index] |= 1 << uint(shift)
          default:
            log.Println("not expecting to index this rune: ", char)
          }

          // lookback
          previous = char
        }
        lineCount += 1
      }

      return theseIndexes
    }
