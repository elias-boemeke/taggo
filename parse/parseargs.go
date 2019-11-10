package parse

import (
  "errors"
  "fmt"
)



func ParseArgs(args []string) (*Options, error) {
  if len(args) == 0 {
    printManualPageOptionsAndExit()
  }

  // map (flag-)keys to flags; show -> flag{...}
  flagDict := getFlagDictionary()
  // map to store parse actions
  parseStatus := getParseStatus()
  // map (command line) flags to keys; -s -> show
  flagKeys := getFlagkeyMap()

  options := newOptions()
  var err  error
  var warn []string
  var warnings []string

  for len(args) > 0 {

    arg := args[0]
    if flagKey, ok := flagKeys[arg]; ok {
      f := flagDict[flagKey]
      args, err, warn = f.parse(arg, args[1:], options, parseStatus)
    } else {
      err, warn = parseFile(arg, options, parseStatus)
      args = args[1:]
    }

    if err != nil {
      return nil, err
    }
    warnings = append(warnings, warn...)
  }

  if *parseStatus["file"] == actionParse {
    return nil, errNoFile()
  }

  if !options.Show.Set {
    change := false
    for _, tag := range options.Tags {
      if tag.Set {
        change = true
        break
      }
    }
    if !change {
      options.Show.Set  = true
      options.Show.Mode = Default
    }
  }

  for _, warn := range warnings {
    LogWarning(warn)
  }

  return options, nil
}

func (f *flag) parse(key string, args []string, options *Options,
  parseStatus map[string]*parseAction) ([]string, error, []string) {

  var err error
  var warnings []string
  var consumed []string

  for _, flagarg := range f.flagArgs {
    if len(args) == 0 {
      if !flagarg.optional {
        return nil, errTooFewArguments(key, f.flagArgs), nil
      }

    } else {
      err = flagarg.validate(args[0])
      if err != nil {
        if flagarg.optional {
          // note that arg is not consumed
          continue
        } else {
          return nil, err, nil
        }
      }
      consumed = append(consumed, args[0])
      args = args[1:]
    }
  }
  warn, err := f.finish(consumed, f, options, parseStatus)
  if err != nil {
    return nil, err, nil
  }
  warnings = append(warnings, warn...)

  return args, nil, warnings
}

func parseFile(arg string, options *Options,
  parseStatus map[string]*parseAction) (error, []string) {

  warnings := make([]string, 0)

  if len(arg) > 0 && arg[0] == '-' {
    warnings = append(warnings, fmt.Sprintf("option '%s' begins with a dash" +
      " but is interpreted as a file; to hide this warning use the" +
      " --file option", arg))
  }

  status := parseStatus["file"]
  if *status == actionParse {
    options.Filename = arg
    *status = actionReparse

  } else {
    // this error is a duplicate, remove duplication by refactoring
    return errors.New(fmt.Sprintf("attempting to use file '%s'," +
      " but file is already set to '%s';" +
      " (make sure you didn't missspell an option)", arg, options.Filename)), nil
  }

  return nil, warnings
}

