package parse

import (
  "errors"
  "fmt"
  "strconv"
  "strings"
)



type Options struct {
  Filename string
  Show ShowOptions
  Tags map[string]*tag
}

type ShowOptions struct {
  Set    bool
  Mode   ShowMode
  Format string
}

type tag struct {
  Set bool
  Value string
}

type showFunc func(ShowMode) bool

type ShowMode int
const (
  Default   ShowMode = iota
  Simple
  Technical
  Full
  Custom
)

type tagInfo struct {
  Short string
  Long  string
  Name  string
  Mutable bool
  Integer bool
  ShowCondition showFunc
  Description string
  DescriptionClear string
}

type parseAction int
const (
  actionParse   parseAction = iota
  actionReparse
  actionIgnore
)

type flag struct {
  flagArgs []flagArg
  finish func([]string, *flag, *Options, map[string]*parseAction) ([]string, error)
}

type numberCondition struct {
  description string
  restriction func(int) bool
}

type flagArg struct {
  pattern  string
  optional bool

  integer   bool
  condition numberCondition 

  restricted bool
  candidates []string
}

func (fa *flagArg) validate(arg string) (error) {
  if fa.integer {
    n, err := strconv.Atoi(arg)
    if err != nil {
      return errors.New(fmt.Sprintf("'%s' is not an integer", arg))
    }
    if !fa.condition.restriction(n) {
      return errors.New(fmt.Sprintf("number '%d' does not match condition '%s'",
        n, fa.condition.description))
    }
    return nil
  }

  if fa.restricted {
    found := false
    for _, v := range fa.candidates {
      if arg == v {
        found = true
        break
      }
    }
    if !found {
      return errors.New(fmt.Sprintf("argument '%s' does not match one of [%s]", arg,
        strings.Join(fa.candidates, ", ")))
    }
  }

  return nil
}

