package parse

import (
  "errors"
  "fmt"
  "os"
  "strings"
)



type parseAction int
const (
  actionParse   parseAction = iota
  actionReparse
  actionIgnore
)

type flag struct{
  keys    []string
  eat     int
  parseFn func(arg string, values []string, op *Options,
               parseActions map[string]parseAction) ([]string, error)
  description string
}

var flags = [...]flag{
  // control options
  flag{[]string{"-h", "--help"},  0, nil,
    "print this message and exit"},

  flag{[]string{"-f", "--file"},  1,
    func(arg string, values []string, op *Options,
         parseActions map[string]parseAction) ([]string, error) {
      return parseFile(values[0], true, op, parseActions)
    },
    "set the file for reading and editing (can be omitted)"},

  // show options
  flag{[]string{"-s", "--show"},  1,
    func(arg string, values []string, op *Options,
        parseActions map[string]parseAction) ([]string, error) {
      return parseShow(arg, values, "show", op, parseActions)
    },
    "mode of printing tags (for more information see README.md)"},

  flag{[]string{"--show-format"}, 1,
    func(arg string, values []string, op *Options,
        parseActions map[string]parseAction) ([]string, error) {
      return parseShow(arg, values, "show-format", op, parseActions)
    },
    "custom format for printing (for more information see README.md)"},

  // set options
  flag{[]string{"-l", "--album"},   1,
    func(arg string, values []string, op *Options,
         parseActions map[string]parseAction) ([]string, error) {
      return parseMutableOption(arg, values, "album", false, op, parseActions)
    },
    "set Album tag"},

  flag{[]string{"-r", "--artist"},  1,
    func(arg string, values []string, op *Options,
         parseActions map[string]parseAction) ([]string, error) {
      return parseMutableOption(arg, values, "artist", false, op, parseActions)
    },
    "set Artist tag"},

  flag{[]string{"-c", "--comment"}, 1,
    func(arg string, values []string, op *Options,
         parseActions map[string]parseAction) ([]string, error) {
      return parseMutableOption(arg, values, "comment", false, op, parseActions)
    },
    "set Comment tag"},

  flag{[]string{"-g", "--genre"},   1,
    func(arg string, values []string, op *Options,
         parseActions map[string]parseAction) ([]string, error) {
      return parseMutableOption(arg, values, "genre", false, op, parseActions)
    },
    "set Genre tag"},

  flag{[]string{"-t", "--title"},   1,
    func(arg string, values []string, op *Options,
         parseActions map[string]parseAction) ([]string, error) {
      return parseMutableOption(arg, values, "title", false, op, parseActions)
    },
    "set Title tag"},

  flag{[]string{"-k", "--track"},   1,
    func(arg string, values []string, op *Options,
         parseActions map[string]parseAction) ([]string, error) {
      return parseMutableOption(arg, values, "track", false, op, parseActions)
    },
    "set Track tag"},

  flag{[]string{"-y", "--year"},    1,
    func(arg string, values []string, op *Options,
         parseActions map[string]parseAction) ([]string, error) {
      return parseMutableOption(arg, values, "year", false, op, parseActions)
    },
    "set Year tag"},

  // unset options
  flag{[]string{"--clear-album"},   0,
    func(arg string, values []string, op *Options,
         parseActions map[string]parseAction) ([]string, error) {
      return parseMutableOption(arg, nil, "album", true, op, parseActions)
    },
    "clear Album tag"},

  flag{[]string{"--clear-artist"},  0,
    func(arg string, values []string, op *Options,
         parseActions map[string]parseAction) ([]string, error) {
      return parseMutableOption(arg, nil, "artist", true, op, parseActions)
    },
    "clear Artist tag"},

  flag{[]string{"--clear-comment"}, 0,
    func(arg string, values []string, op *Options,
         parseActions map[string]parseAction) ([]string, error) {
      return parseMutableOption(arg, nil, "comment", true, op, parseActions)
    },
    "clear Comment tag"},

  flag{[]string{"--clear-genre"},   0,
    func(arg string, values []string, op *Options,
         parseActions map[string]parseAction) ([]string, error) {
      return parseMutableOption(arg, nil, "genre", true, op, parseActions)
    },
    "clear Genre tag"},

  flag{[]string{"--clear-title"},   0,
    func(arg string, values []string, op *Options,
         parseActions map[string]parseAction) ([]string, error) {
      return parseMutableOption(arg, nil, "title", true, op, parseActions)
    },
    "clear Title tag"},

  flag{[]string{"--clear-track"},   0,
    func(arg string, values []string, op *Options,
         parseActions map[string]parseAction) ([]string, error) {
      return parseMutableOption(arg, nil, "track", true, op, parseActions)
    },
    "clear Track tag"},

  flag{[]string{"--clear-year"},    0,
    func(arg string, values []string, op *Options,
         parseActions map[string]parseAction) ([]string, error) {
      return parseMutableOption(arg, nil, "year", true, op, parseActions)
    },
    "clear Year tag"},

  // clear option
  flag{[]string{"--clear"}, 0,
    func(arg string, values []string, op *Options,
        parseActions map[string]parseAction) ([]string, error) {
      return parseClear(arg, op, parseActions)
    },
    "clear all tags"},
}

var mutableTags = [...]string{"album", "artist", "comment",
  "genre", "title", "track", "year"}



func printManualAndExit() {
  help := display.Fat("taggo") + " is a tool for reading and editing meta data" +
    " embedded into audio files\n" +
    "For more information see README.md\n" +
    "\n" +
    display.Fat("Usage\n") +
    "        " + "taggo [options] <file> ..." +
    "\n" +
    "\n" +
    display.Fat("Options\n")

  for _, f := range flags {
    help += "        " +
      fmt.Sprintf("%-20s", strings.Join(f.keys, " or ")) +
      f.description + "\n"
  }

  help += "\n" +
    "\n" +
    display.Fat("Show Modes") + "\n" +
    "        " + "default | simple | tech | full\n" +
    "\n" +
    "\n" +
    display.Fat("Examples") + "\n" +
    "        " + "taggo test.mp3\n" +
    "        " + "taggo test.mp3 --clear\n" +
    "        " + "taggo test.mp3 -r \"The Artist\"\n" +
    "        " + "taggo test.mp3 -c \"A Comment\" -s simple\n" +
    "        " + "taggo test.mp3 --clear-genre --show-format \"year;%y\\nalbum;%l\"\n" +
    "        " + "taggo -f -dashfile -k 5\n" +
    "\n" +
    "\n" +
    display.Fat("Note") + "\n" +
    "        " + "Further reading on the Options, formatting with --show-format" +
    " and general information\n" +
    "        " + "can be found in README.md or on github:" +
    " 'https://github.com/elias-boemeke/taggo'\n"


  fmt.Println(help)
  os.Exit(0)
}

func lookupFlag(key string) *flag {
  for _, f := range flags {
    for _, k := range f.keys {
      if key == k {
        return &f
      }
    }
  }
  return nil
}

func ParseArgs(args []string) (Options, error) {

  if len(args) == 0 {
    printManualAndExit()
  }
  flags[0].parseFn = func(arg string, values []string, op *Options,
      parseActions map[string]parseAction) ([]string, error) {
    printManualAndExit()
    return nil, nil
  }

  var op Options

  warnings := make([]string, 0)

  parseActions := make(map[string]parseAction)
  parseKeys := []string{"file", "show", "clear"}
  for _, key := range mutableTags {
    parseKeys = append(parseKeys, key)
    parseKeys = append(parseKeys, "clear-" + key)
  }
  for _, k := range parseKeys {
    parseActions[k] = actionParse
  }

  elvis := func(cond bool, ts string, fs string) string {
    if cond {
      return ts
    } else {
      return fs
    }
  }

  var warn []string
  var err  error
  // iterate through arguments
  for i := 0; i < len(args); i++ {
    // search for a matching key
    f := lookupFlag(args[i])

    if f == nil {
      warn, err = parseFile(args[i], false, &op, parseActions)

    } else {
      if i + f.eat >= len(args) {
        return Options{}, errors.New(fmt.Sprintf("option '%s' expects" + 
          " %d " + elvis(f.eat == 1, "value", "values") +
          ", got %d (unexpected end of arguments)",
          args[i], f.eat, len(args)-i-1))
      }
      warn, err = f.parseFn(args[i], args[i+1:i+1+f.eat], &op, parseActions)
      i += f.eat
    }

    if err != nil {
      return Options{}, err
    }
    warnings = append(warnings, warn...)
  }

  if parseActions["file"] == actionParse {
    return Options{}, errors.New("the tool can't be run without a file;" +
      " specify it in the arguments")
  }

  // if no tags were modified (and ShowMode was not set), set ShowMode to default
  if !op.Show.Set && !(op.Album.Set || op.Artist.Set || op.Comment.Set ||
      op.Genre.Set || op.Title.Set || op.Track.Set || op.Year.Set) {
    op.Show.Set = true
    op.Show.Mode = Default
  }

  yellow := func(s string) string { return "\033[33m" + s + "\033[0m" }
  // print warnings
  for _, warn := range warnings {
    fmt.Fprintln(os.Stderr, fmt.Sprintf(yellow("WARN") + " %s", warn))
  }

  return op, nil;
}

