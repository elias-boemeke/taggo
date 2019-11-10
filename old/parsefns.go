package parse

import (
  "errors"
  "fmt"
  "strconv"
)



func parseFile(arg string, calledWithFlag bool, op *Options,
  parseActions map[string]parseAction) ([]string, error) {

  warnings := make([]string, 0)

  if !calledWithFlag && len(arg) > 0 && arg[0] == '-' {
    warnings = append(warnings, fmt.Sprintf("option '%s' begins with a dash" +
      " but is interpreted as a file; to hide this warning use the" +
      " --file option", arg))
  }

  if parseActions["file"] == actionParse {
    op.FileName = arg
    parseActions["file"] = actionReparse

  } else {
    return nil, errors.New(fmt.Sprintf("attempting to use file '%s'," +
      " but file is already set to '%s'", arg, op.FileName))
  }

  return warnings, nil
}

// it is guaranteed that exactly one the return values is nil and the other one is not 
func tagFromKey(key string, op *Options) (*tag, *itag) {
  var t  *tag  = nil
  var it *itag = nil

  switch key {
  case "album":
    t = &op.Album
  case "artist":
    t = &op.Artist
  case "comment":
    t = &op.Comment
  case "genre":
    t = &op.Genre
  case "title":
    t = &op.Title
  case "track":
    it = &op.Track
  case "year":
    it = &op.Year
  default:
    panic("unreachable code :(")
  }

  return t, it
}

func parseMutableOption(arg string, values []string, key string, clearTag bool, op *Options,
  parseActions map[string]parseAction) ([]string, error) {

  t, it := tagFromKey(key, op)
  intType := (t == nil)
  actionKey := key
  counterKey := "clear-" + key
  if clearTag {
    actionKey, counterKey = counterKey, actionKey
  }

  action := parseActions[actionKey]
  if counterAction := parseActions[counterKey];
    action != actionIgnore && counterAction != actionParse {
    
    parseActions[actionKey] = actionIgnore
    parseActions[counterKey] = actionIgnore
    if intType {
      it.Set = false
      it.Value = 0
    } else {
      t.Set = false
      t.Value = ""
    }
    return []string{fmt.Sprintf("conflicting options for tag '%s' at '%s'," +
      " old value is kept", key, arg)}, nil
  }

  warn := make([]string, 0)
  switch action {

  case actionParse:
    if intType {
      it.Set = true
      // setting logic for int type
      if clearTag {
        it.Value = 0
      } else {
        n, err := strconv.Atoi(values[0])
        if err != nil {
          return nil, err
        }
        if n < 1 {
          return nil, errors.New(fmt.Sprintf("value for tag '%s' has to be" +
            " greater than 0 (got %d)", key, n))
        }
        it.Value = n
      }

    } else {
      // setting logic for string type
      t.Set = true
      if clearTag {
        t.Value = ""
      } else {
        t.Value = values[0]
      }
    }
    // in any case update the parse action
    parseActions[actionKey] = actionReparse

  case actionReparse:
    var v string
    if intType {
      if it.Value == 0 {
        v = ""
      } else {
        v = strconv.Itoa(it.Value)
      }

    } else {
      v = t.Value
    }
    warn = append(warn, fmt.Sprintf("tag '%s' already set, value remains" +
      " '%s' (ignoring)", key, v))
    parseActions[actionKey] = actionIgnore

  case actionIgnore:
    // do nothing
  default:
    panic("unreachable code :(")
  }

  return warn, nil
}

func parseClear(arg string, op *Options,
  parseActions map[string]parseAction) ([]string, error) {

  key := "clear"
  var w []string
  warn := make([]string, 0)
  var err error
  switch parseActions[key] {
  case actionParse:
    for _, key := range mutableTags {
      w, err = parseMutableOption(arg, nil, key, true, op, parseActions)
      if err != nil {
        return nil, err
      }
      warn = append(warn, w...)
    }
    parseActions[key] = actionReparse

  case actionReparse:
    warn = append(warn, fmt.Sprintf("option '%s' already set (ignoring)", key))

  case actionIgnore:
    // do nothing
  default:
    panic("unreachable code :(")
  }

  return warn, nil
}

func showModeFromString(s string) (ShowMode, error) {
  var mode ShowMode
  switch s {
  case "default":
    mode = Default
  case "simple":
    mode = Simple
  case "tech":
    mode = Technical
  case "full":
    mode = Full
  default:
    return Default, errors.New(fmt.Sprintf("'%s' is not a valid show mode", s))
  }
  return mode, nil
}

func parseShow(arg string, values []string, key string, op *Options,
  parseActions map[string]parseAction) ([]string, error) {
    
  switch parseActions["show"] {
  case actionParse:
    switch key {
    case "show":
      mode, err := showModeFromString(values[0])
      if err != nil {
        return nil, err
      }
      op.Show.Set = true
      op.Show.Mode = mode

    case "show-format":
      op.Show.Set = true
      op.Show.Mode = Custom
      op.Show.Format = values[0]

    default:
      panic("unreachable code :(")
    }
    parseActions["show"] = actionReparse

  case actionReparse:
    parseActions["show"] = actionIgnore
    return []string{fmt.Sprintf("'%s' mode for showing already given (ignoring)", key)}, nil

  case actionIgnore:
    // do nothing

  default:
    panic("unreachable code :(")
  }

  return nil, nil
}

