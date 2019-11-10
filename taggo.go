package main

import (
  "os"
)

import (
  parse  "github.com/elias-boemeke/taggo/parse"
  tag  "github.com/elias-boemeke/taggo/tag"
)



func main() {
  options, err := parse.ParseArgs(os.Args[1:])
  if err != nil {
    parse.LogErrorAndDie(parse.RefManual, "parsing of arguments failed: %s", err)
  }

  fileName := options.Filename
  file := tag.ReadFile(fileName)
  defer file.Close()

  err = tag.WriteTags(file, options)
  if err != nil {
    parse.LogErrorAndDie(parse.NoRefManual, "failed to write tags: %s", err)
  }

  if options.Show.Set {
    tag.ShowTags(file, &options.Show)
  }
}

/*
-------------------------
  Available Tags
-------------------------
  Album       ~   string
  Artist      ~   string
  Bitrate         int
  Channels        int
  Comment     ~   string
  Genre       ~   string
  Length          time.Duration
  Samplerate      int
  Title       ~   string
  Track       ~   int
  Year        ~   int
-------------------------
   Escapes
-------------------------
  %l : Album
  %r : Artist
  %b : Bitrate
  %h : Channels
  %c : Comment
  %g : Genre
  %n : Length
  %s : Samplerate
  %t : Title
  %k : Track
  %y : Year
  %% : %
-------------------------
   Flags
-------------------------
  -h or --help        print help (this message) and exit, additional parameters: show, examples
  -f or --file        set the file for reading and editing (flag can be omitted)
  -s or --show        mode of printing tags, leaving out mode defaults to show mode default
  --show-format       custom format for printing tags
  -l or --album       set Album tag
  -r or --artist      set Artist tag
  -c or --comment     set Comment tag
  -g or --genre       set Genre tag
  -t or --title       set Title tag
  -k or --track       set Track tag
  -y or --year        set Year tag
  --clear-album       clear Album tag
  --clear-artist      clear Artist tag
  --clear-comment     clear Comment tag
  --clear-genre       clear Genre tag
  --clear-title       clear Title tag
  --clear-track       clear Track tag
  --clear-year        clear Year tag
  --clear             clear all tags
-------------------------
*/

