package parse

import (
  "errors"
  "fmt"
  "os"
  "strings"
)



// formatters

func red(s string) string {
  return "\033[31m" + s + "\033[0m"
}

func yellow(s string) string {
  return "\033[33m" + s + "\033[0m"
}

func fat(s string) string {
  return "\033[1m" + s + "\033[0m"
}

// console outputers

func LogErrorAndDie(showManualReference bool, format string, args ...interface{}) {
  fmt.Fprintln(os.Stderr, fmt.Sprintf(red("ERR") + " " + format, args...))
  if showManualReference {
    fmt.Println("\nuse `taggo --help` to view the manual")
  }
  os.Exit(1)
}

func LogWarning(message string) {
  fmt.Fprintln(os.Stderr, fmt.Sprintf(yellow("WARN") + " %s", message))
}

// error templates

func errNoFile() error {
  return errors.New("the tool can't be run without an audio file," +
    " specify it in the arguments (file needed)")
}

func errTooFewArguments(key string, flagArgs []flagArg) error {
  var flagNames []string
  for _, f := range flagArgs {
    flagNames = append(flagNames, f.pattern)
  }
  pattern := key
  if len(flagNames) > 0 {
    pattern += " " + strings.Join(flagNames, " ")
  }
  return errors.New(fmt.Sprintf("too few arguments in option '%s', pattern is: %s",
    key, pattern))
}

// manual

func printManualPageOptionsAndExit() {
  flags := getFlagDictionary()
  help := fat("taggo") + "\n" +
    " ...is a tool for reading and editing meta data" +
    " embedded into audio files\n" +
    "\n" +
    fat("Usage\n") +
    "        " + "taggo [options...] <file>\n" +
    "\n" +
    fat("Options\n")

  help += "      " + fat("set tag") + "\n"
  for _, t := range tags {
    if t.Mutable {
      help += "        " +
        fmt.Sprintf("%-28s", "-" + t.Short + ", --" + t.Long +
        " " + flags[t.Long].flagArgs[0].pattern) +
        t.Description + "\n"
    }
  }
  help += "\n"
  help += "      " + fat("clear tag(s)") + "\n"
  for _, t := range tags {
    if t.Mutable {
      help += "        " +
        fmt.Sprintf("%-28s", "--clear-" + t.Long) +
        t.DescriptionClear + "\n"
    }
  }
  help += "\n"
  help += "        " +
    fmt.Sprintf("%-28s", "--clear") + "clear all tags\n" +
    "\n"
  help += "      " + fat("display tags") + " (see Presentation)\n" +
    "        " +
    fmt.Sprintf("%-28s", "-s, --show " +
    flags["show"].flagArgs[0].pattern) +
    "show the tags defined by mode\n" +
    "        " +
    fmt.Sprintf("%-28s", "--show-format " +
    flags["show-format"].flagArgs[0].pattern) +
    "display tags and custom text defined by format\n" +
    "\n"
  hps := flags["help"].flagArgs[0].candidates[0]
  hpe := flags["help"].flagArgs[0].candidates[1]
  fpat := flags["file"].flagArgs[0].pattern
  help += "      " + fat("miscellaneous") + "\n" +
    "        " + fmt.Sprintf("%-28s", "-h, --help [" + hps + "|" + hpe + "]") +
    "show help page\n" +
    "        " + fmt.Sprintf("%-28s", "-f, --file " + fpat) +
    "explicitly take " + fpat + " as input file\n" +
    "\n"

  help += fat("Presentation") + "\n" +
    "        taggo --help " + hps + "\n" +
    "         ...to get further help on how to display the tags\n" +
    "            with the options --show and --show-format\n" +
    "\n" +
    fat("Examples") + "\n" +
    "        taggo --help " + hpe + "\n" +
    "         ...to show examples on how to use taggo\n" +
    "\n" +
    fat("Note") + "\n" +
    "        check out the source code on github: https://github.com/elias-boemeke/taggo"


  fmt.Println(help)
  os.Exit(0)
}

func printManualPageShowAndExit() {
  flags := getFlagDictionary()
  help := fat("taggo") + "\n" +
    " ...is a tool for reading and editing meta data" +
    " embedded into audio files\n" +
    "    showing '" + fat("show") + "' help page, for main page use taggo --help\n" +
    "\n"
  help += fat("Presentation") + "\n" +
    "      " + fmt.Sprintf("%-28s", "-s, --show " +
    flags["show"].flagArgs[0].pattern) + "\n" +
    "\n" +
    "        show the tags defined by MODE\n" +
    "        MODE is optional and can be ommited\n" +
    "        in this case mode defaults to 'default'\n" +
    "\n" +
    "        available modes:  " +
    fat("default") + ", " + fat("simple") + ", " +
    fat("technical") + ", " + fat("full") + "\n" +
    "\n" +
    "        there can only be one mode active at a time\n" +
    "        if you want a custom format use --show-format\n" +
    "\n" +
    "       tag        | shown by\n" +
    "       ------------------------------------\n" +
    "       Album      | default, simple, full\n" +
    "       Artist     | default, simple, full\n" +
    "       Bitrate    | technical, full\n" +
    "       Channels   | technical, full\n" +
    "       Comment    | default, full\n" +
    "       Genre      | default, full\n" +
    "       Length     | simple, technical, full\n" +
    "       Samplerate | technical, full\n" +
    "       Title      | default, simple, full\n" +
    "       Track      | default, simple, full\n" +
    "       Year       | default, full\n" +
    "\n" +
    "\n" +
    "      " + fmt.Sprintf("%-28s", "--show-format " +
    flags["show-format"].flagArgs[0].pattern) + "\n" +
    "\n" +
    "        display tags and custom text defined by format\n" +
    "        format is string that may contain the following escapes:\n" +
    "\n" +
    "       escape | expands to\n" +
    "       -----------------------\n" +
    "       %l     | Album tag\n" +
    "       %r     | Artist tag\n" +
    "       %b     | Bitrate tag\n" +
    "       %h     | Channels tag\n" +
    "       %c     | Comment tag\n" +
    "       %g     | Genre tag\n" +
    "       %n     | Length tag\n" +
    "       %s     | Samplerate tag\n" +
    "       %t     | Title tag\n" +
    "       %k     | Track tag\n" +
    "       %y     | Year\n" +
    "       %%     | literal %\n" +
    "\n" +
    "        after these escapes are resolved, the string is\n" +
    "        passed to strconv.Unquote such that you can use '\\n' and other escapes\n" +
    "\n" +
    fat("Note") + "\n" +
    "        for examples see taggo --help examples"

  fmt.Println(help)
  os.Exit(0)
}

func printManualPageExamplesAndExit() {
  help := fat("taggo") + "\n" +
    " ...is a tool for reading and editing meta data" +
    " embedded into audio files\n" +
    "    showing '" + fat("examples") + "' help page, for main page use taggo --help\n" +
    "\n"
  help += fat("Examples") + "\n" +
    "      " + "taggo test.mp3\n" +
    "        show tag information of file 'test.mp3'\n" +
    "\n" +
    "      " + "taggo test.mp3 --clear\n" +
    "        clear all tags of file 'test.mp3'\n" +
    "\n" +
    "      " + "taggo test.mp3 -r \"The Artist\"\n" +
    "        change artist tag of file 'test.mp3' to 'The Artist'\n" +
    "\n" +
    "      " + "taggo test.mp3 -c \"A Comment\" -s simple\n" +
    "        change comment tag of file 'test.mp3' to 'A Comment' and\n" +
    "        display the tags in simple mode afterwards\n" +
    "\n" +
    "      " + "taggo test.mp3 --clear-genre --show-format \"year;%y\\nalbum;%l\"\n" +
    "        clear the genre tag of file 'test.mp3' and\n" +
    "        display tags using a custom given format\n" +
    "\n" +
    "      " + "taggo -f -dashfile -k 5\n" +
    "        change track number tag of file '-dashfile' to 5"

  fmt.Println(help)
  os.Exit(0)
}

