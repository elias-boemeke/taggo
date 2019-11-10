package parse

import (
  "errors"
  "fmt"
  "strings"
)



var tags = [...]tagInfo {
  {"l", "album",      "Album",      true,  false, showProd,   "set Album tag",   "clear Album tag"},
  {"r", "artist",     "Artist",     true,  false, showProd,   "set Artist tag",  "clear Artist tag"},
  {"b", "bitrate",    "Bitrate",    false, true,  showTech,   "", ""},
  {"h", "channels",   "Channels",   false, true,  showTech,   "", ""},
  {"c", "comment",    "Comment",    true,  false, showExtra,  "set Comment tag", "clear Comment tag"},
  {"g", "genre",      "Genre",      true,  false, showExtra,  "set Genre tag",   "clear Genre tag"},
  {"n", "length",     "Length",     false, false, showLength, "", ""},
  {"s", "samplerate", "Samplerate", false, true,  showTech,   "", ""},
  {"t", "title",      "Title",      true,  false, showProd,   "set Title tag",   "clear Title tag"},
  {"k", "track",      "Track",      true,  true,  showProd,   "set Track tag",   "clear Track tag"},
  {"y", "year",       "Year",       true,  true,  showExtra,  "set Year tag",    "clear Year tag"},
}

// used for LogErrorAndDie to indicate if an
// additional reference to the manual is shown
const RefManual = true
const NoRefManual = false

var showProd = func(mode ShowMode) bool {
  return mode == Default || mode == Simple || mode == Full
}

var showTech = func(mode ShowMode) bool {
  return mode == Technical || mode == Full
}

var showExtra = func(mode ShowMode) bool {
  return mode == Default || mode == Full
}

var showLength = func(mode ShowMode) bool {
  return mode == Simple || mode == Technical || mode == Full
}

var info []*tagInfo
var shortToLong map[string]string

func GetTagInfo() []*tagInfo {
  if info == nil {
    info = make([]*tagInfo, 0)
    for i := range tags {
      info = append(info, &tags[i])
    }
  }
  return info
}

func GetShortToLongMap() map[string]string {
  if len(shortToLong) == 0 {
    shortToLong = make(map[string]string)
    for _, t := range tags {
      shortToLong[t.Short] = t.Long
    }
  }
  return shortToLong
}

func getFlagDictionary() map[string]*flag {
  flags := make(map[string]*flag)

  // -h or --help
  flags["help"] = &flag{
    flagArgs: []flagArg{
      flagArg{
        pattern: "page",
        optional: true,
        restricted: true,
        candidates: []string{"show", "examples"},
      },
    },
    finish: func(args []string, f *flag, options *Options,
        parseStatus map[string]*parseAction) ([]string, error) {
      if len(args) == 0 {
        printManualPageOptionsAndExit()
      } else {
        if args[0] == "show" {
          printManualPageShowAndExit()
        } else if args[0] == "examples" {
          printManualPageExamplesAndExit()
        }
      }
      return nil, nil
    },
  }

  // -f or --file
  flags["file"] = &flag{
    flagArgs: []flagArg{
      flagArg{
        pattern: "FILE",
      },
    },
    finish: func(args []string, f *flag, options *Options,
        parseStatus map[string]*parseAction) ([]string, error) {
      status := parseStatus["file"]
      if *status == actionParse {
        options.Filename = args[0]
        *status = actionReparse
        return nil, nil
      } else {
        // this error is a duplicate, remove duplication by refactoring
        return nil, errors.New(fmt.Sprintf("attempting to use file '%s'," +
          " but file is already set to '%s'", args[0], options.Filename))
      }
    },
  }

  // -s or --show
  flags["show"] = &flag{
    flagArgs: []flagArg{
      flagArg{
        pattern: "[MODE]",
        optional: true,
        restricted: true,
        candidates: []string{"default", "simple", "technical", "full"},
      },
    },
    finish: func(args []string, f *flag, options *Options,
        parseStatus map[string]*parseAction) ([]string, error) {
      status := parseStatus["show"]
      switch *status {
        // actionParse
      case actionParse:
        var mode ShowMode
        if len(args) == 0 {
          mode = Default
        } else {
          switch args[0] {
          case "default":
            mode = Default
          case "simple":
            mode = Simple
          case "technical":
            mode = Technical
          case "full":
            mode = Full
          }
        }
        options.Show.Set = true
        options.Show.Mode = mode
        *status = actionReparse
        return nil, nil

        // actionReparse
      case actionReparse:
        *status = actionIgnore
        return []string{"show mode already given (ignoring)"}, nil

        // actionIgnore
      case actionIgnore:
        return nil, nil

      default:
        panic("unreachable :(")
      }
    },
  }

  // --show-format
  flags["show-format"] = &flag{
    flagArgs: []flagArg{
      flagArg{
        pattern: "FORMAT",
      },
    },
    finish: func(args []string, f *flag, options *Options,
        parseStatus map[string]*parseAction) ([]string, error) {
      status := parseStatus["show"]
      switch *status {
        // actionParse
      case actionParse:
        options.Show.Set = true
        options.Show.Mode = Custom
        options.Show.Format = args[0]
        *status = actionReparse
        return nil, nil

        // actionReparse
      case actionReparse:
        *status = actionIgnore
        return []string{"show mode already given (ignoring)"}, nil

        // actionIgnore
      case actionIgnore:
        return nil, nil

      default:
        panic("unreachable :(")
      }
    },
  }

  // tags
  for _, t := range(tags) {
    // for closure capturing
    key := t.Long

    if t.Mutable {
      // set flag
      var fa []flagArg
      var zeroval string

      if t.Integer {
        fa = []flagArg{
          flagArg{
            pattern: strings.ToUpper(t.Long),
            integer: true,
            condition: numberCondition{
              description: "x > 0",
              restriction: func(x int) bool { return x > 0 },
            },
          },
        }
        zeroval = "0"

      } else {
        fa = []flagArg{
          flagArg{
            pattern: strings.ToUpper(t.Long),
          },
        }
        zeroval = ""
      }

      flags[t.Long] = &flag{
        flagArgs: fa,
        finish: func(args []string, f *flag, options *Options,
            parseStatus map[string]*parseAction) ([]string, error) {
          return parseTag(key, args[0], f, options, parseStatus)
        },
      }
      // clear flag
      flags["clear-" + t.Long] = &flag{
        flagArgs: []flagArg{},
        finish: func(args []string, f *flag, options *Options,
            parseStatus map[string]*parseAction) ([]string, error) {
          return parseTag(key, zeroval, f, options, parseStatus)
        },
      }
    }
  }

  // clear
  flags["clear"] = &flag{
    flagArgs: []flagArg{},
    finish: func(args []string, f *flag, options *Options,
        parseStatus map[string]*parseAction) ([]string, error) {
      var warn []string
      for _, t := range(tags) {
        if t.Mutable {
          var zeroval string
          if t.Integer {
            zeroval = "0"
          } else {
            zeroval = ""
          }
          w, err := parseTag(t.Long, zeroval, f, options, parseStatus)
          if (err != nil) {
            return nil, err
          }
          warn = append(warn, w...)
        }
      }
      return warn, nil
    },
  }

  return flags
}

func getParseStatus() map[string]*parseAction {

  parseStatus := make(map[string]*parseAction)
  extraKeys := []string{"help", "file", "show"}

  for _, k := range(extraKeys) {
    parseStatus[k] = new(parseAction)
  }
  for _, t := range(tags) {
    if t.Mutable {
      k := t.Long
      parseStatus[k] = new(parseAction)
    }
  }

  return parseStatus
}

func getFlagkeyMap() map[string]string {
  keys := make(map[string]string)

  keys["-h"] = "help"
  keys["--help"] = keys["-h"]

  keys["-f"] = "file"
  keys["--file"] = keys["-f"]

  keys["-s"] = "show"
  keys["--show"] = keys["-s"]

  keys["--show-format"] = "show-format"

  // tags
  for _, t := range(tags) {
    if t.Mutable {
      keys["-" + t.Short] = t.Long
      keys["--" + t.Long] = t.Long
      keys["--clear-" + t.Long] = "clear-" + t.Long
    }
  }

  keys["--clear"] = "clear"

  return keys
}

func newOptions() (*Options) {
  op := &Options{}
  op.Tags = make(map[string]*tag)

  for _, t := range tags {
    if t.Mutable {
      op.Tags[t.Long] = &tag{}
    }
  }

  return op
}

func parseTag(key string, value string, f *flag, options *Options,
    parseStatus map[string]*parseAction) ([]string, error) {
  status := parseStatus[key]
  switch *status {
    // actionParse
  case actionParse:
    opt := options.Tags[key]
    opt.Set = true
    opt.Value = value
    *status = actionReparse
    return nil, nil

    // actionReparse
  case actionReparse:
    *status = actionIgnore
    return []string{fmt.Sprintf("tag '%s' already set, value remains" +
      " '%s' (ignoring)", key, options.Tags[key].Value)}, nil

    // actionIgnore
  case actionIgnore:
    return nil, nil

  default:
    panic("unreachable :(")
  }
}

