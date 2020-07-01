package speller

import ( "os"
  "bufio"
  "log"
)

// cardinality is a-z plus apostrophy
const bucketSize = 8
// size is derived from number of words in file
const bitmapSize = (109583/bucketSize)+1

type Word struct {
  Raw string
  Length int
}

type Indexes struct {
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
}

// ##################### PUBLIC INTERFACE: THE DATA STRUCTURE ######
// This is what we use to look things up
var Dictionary = initIndexes()

// build indexes
func Build(file string) {
  Dictionary = buildIndex(file)
}

func Search(char string) []Word {
  var matched []Word
  var residents = bytesToLocation(Dictionary.z)
  for at := 0; at < len(residents); at++ {
    matched = append(matched,Dictionary.words[residents[at]])
  }
  return matched
}
// ###################### END PUBLIC INTERFACE ####################

// make new array init with fill 0
func initByteArray() [bitmapSize]byte {
  newIndex := [bitmapSize]byte{}
  for i := 0; i < bitmapSize; i++ {
    newIndex[i] = 0
  }
  return newIndex
}

// make new indexes fill with 0
// lots of code now, easier to read later
func initIndexes() Indexes {
  newIndexes := Indexes{
    a: initByteArray(),
    b: initByteArray(),
    c: initByteArray(),
    d: initByteArray(),
    e: initByteArray(),
    f: initByteArray(),
    g: initByteArray(),
    h: initByteArray(),
    i: initByteArray(),
    j: initByteArray(),
    k: initByteArray(),
    l: initByteArray(),
    m: initByteArray(),
    n: initByteArray(),
    o: initByteArray(),
    p: initByteArray(),
    q: initByteArray(),
    r: initByteArray(),
    s: initByteArray(),
    t: initByteArray(),
    u: initByteArray(),
    v: initByteArray(),
    w: initByteArray(),
    x: initByteArray(),
    y: initByteArray(),
    z: initByteArray(),
    apos: initByteArray(),
   }
   return newIndexes
}

// need both index by byte and shift for bits
func byteIndex(index int) (int, int) {
  if index == 0 { return 0,0 }
  // less one because index starts at zero
  position := index / bucketSize
  // remainder reversed 7 -> 0 .. 0 -> 7
  shift := bucketSize - (index % bucketSize) -1
  return position, shift
}

// convert byte matches to index positions
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

// this tells us how many byte arrays we need
// one bytearray index per words
func buildIndex(file string) Indexes {
  var theseIndexes = initIndexes()

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
    // add our word to index
    theseIndexes.words = append(theseIndexes.words, Word{Raw:word, Length: len(word)})

    if (lineCount > 109571 || word == "year's") {
      // loop over each rune and add to indexes
      for i := 0; i < len(word); i++ {
        char := string([]rune(word)[i])
        index, shift := byteIndex(lineCount)
        // lots of code , and easy to read
        // better then compact, obtuse code #pickyourpoison
        switch char {
        case "a":
          theseIndexes.a[index] |= 1 << uint(shift)
        case "b":
          theseIndexes.b[index] |= 1 << uint(shift)
        case "c":
          theseIndexes.c[index] |= 1 << uint(shift)
        case "d":
          theseIndexes.d[index] |= 1 << uint(shift)
        case "e":
          theseIndexes.e[index] |= 1 << uint(shift)
        case "f":
          theseIndexes.f[index] |= 1 << uint(shift)
        case "g":
          theseIndexes.g[index] |= 1 << uint(shift)
        case "h":
          theseIndexes.h[index] |= 1 << uint(shift)
        case "i":
          theseIndexes.i[index] |= 1 << uint(shift)
        case "j":
          theseIndexes.j[index] |= 1 << uint(shift)
        case "k":
          theseIndexes.k[index] |= 1 << uint(shift)
        case "l":
          theseIndexes.l[index] |= 1 << uint(shift)
        case "m":
          theseIndexes.m[index] |= 1 << uint(shift)
        case "n":
          theseIndexes.n[index] |= 1 << uint(shift)
        case "o":
          theseIndexes.o[index] |= 1 << uint(shift)
        case "p":
          theseIndexes.p[index] |= 1 << uint(shift)
        case "q":
          theseIndexes.q[index] |= 1 << uint(shift)
        case "r":
          theseIndexes.r[index] |= 1 << uint(shift)
        case "s":
          theseIndexes.s[index] |= 1 << uint(shift)
        case "t":
          theseIndexes.t[index] |= 1 << uint(shift)
        case "u":
          theseIndexes.u[index] |= 1 << uint(shift)
        case "v":
          theseIndexes.v[index] |= 1 << uint(shift)
        case "w":
          theseIndexes.w[index] |= 1 << uint(shift)
        case "x":
          theseIndexes.x[index] |= 1 << uint(shift)
        case "y":
          theseIndexes.y[index] |= 1 << uint(shift)
        case "z":
          theseIndexes.z[index] |= 1 << uint(shift)
        case "'":
          theseIndexes.apos[index] |= 1 << uint(shift)
        default:
          log.Println("not expecting to index this rune: ", char)
        }
      }
    }
    lineCount += 1
  }

  return theseIndexes
}
