package tag

import (
  "fmt"
  "strconv"
  "unicode/utf8"
)

import (
  parse "github.com/elias-boemeke/taggo/parse"
  taglib "github.com/wtolson/go-taglib"
)



func ShowTags(file *taglib.File, showOpt *parse.ShowOptions) {
  tagValues := tagValuesFromFile(file)
  if showOpt.Mode == parse.Custom {
    showTagsFromFormat(tagValues, showOpt.Format)
  } else {
    showTagsFromMode(tagValues, showOpt.Mode)
  }
}

func showTagsFromMode(tagValues map[string]string, mode parse.ShowMode) {
  var width string

  switch mode {
  case parse.Default:
    width = "7"
  case parse.Simple:
    width = "6"
  case parse.Technical:
    width = "10"
  case parse.Full:
    width = "10"
  }

  info := parse.GetTagInfo()
  for _, v := range info {
    if v.ShowCondition(mode) {
      fmt.Println(fmt.Sprintf("%" + width + "s: %s", v.Name, tagValues[v.Long]))
    }
  }
}

func showTagsFromFormat(tagValues map[string]string, format string) {
  show := ""
  stl := parse.GetShortToLongMap()

  for i := 0; i < len(format); {
    r, w := utf8.DecodeRuneInString(format[i:])

    if r == '%' {
      if i + w > len(format) - 1 {
        show += string(r)

      } else {
        next, wNext := utf8.DecodeRuneInString(format[i+w:])

        if next == '%' {
          show += string(next)
        } else if v, ok := stl[string(next)]; ok {
          show += tagValues[v]
        } else {
          show += string(r)
          show += string(next)
        }
        w += wNext
      }

    } else {
      show += string(r)
    }

    i += w
  }

  s, err := strconv.Unquote(fmt.Sprintf("\"%s\"", show))
  if err != nil {
    s = "Unable to resolve format string '" + format +
      "', make sure to escape quotes(\") and see" +
      " 'https://golang.org/pkg/strconv/#Unquote' for more information"
  }
  fmt.Println(s)
}

