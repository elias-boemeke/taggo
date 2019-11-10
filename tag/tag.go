package tag

import (
  "strconv"
)

import (
  taglib "github.com/wtolson/go-taglib"
)

import (
  parse "github.com/elias-boemeke/taggo/parse"
)



func ReadFile(fileName string) *taglib.File {
  file, err := taglib.Read(fileName)

  if err != nil {
    parse.LogErrorAndDie(parse.NoRefManual, "unable to read file '%s': %s", fileName, err)
  }

  if file == nil {
    parse.LogErrorAndDie(parse.NoRefManual, "unable to read file '%s'", fileName)
  }

  return file
}

func WriteTags(file *taglib.File, op *parse.Options) error {
  forceInt := func(s string) int {
    n, _ := strconv.Atoi(s)
    return n
  }

  if op.Tags["album"].Set {
    file.SetAlbum(op.Tags["album"].Value)
  }
  if op.Tags["artist"].Set {
    file.SetArtist(op.Tags["artist"].Value)
  }
  if op.Tags["comment"].Set {
    file.SetComment(op.Tags["comment"].Value)
  }
  if op.Tags["genre"].Set {
    file.SetGenre(op.Tags["genre"].Value)
  }
  if op.Tags["title"].Set {
    file.SetTitle(op.Tags["title"].Value)
  }
  if op.Tags["track"].Set {
    file.SetTrack(forceInt(op.Tags["track"].Value))
  }
  if op.Tags["year"].Set {
    file.SetYear(forceInt(op.Tags["year"].Value))
  }

  err := file.Save()
  return err
}

func tagValuesFromFile(file *taglib.File) map[string]string {
  values := make(map[string]string)
  strHideZero := func(n int) string {
    if n == 0 {
      return ""
    }
    return strconv.Itoa(n)
  }
  values["album"]      = file.Album()
  values["artist"]     = file.Artist()
  values["bitrate"]    = strconv.Itoa(file.Bitrate())
  values["channels"]   = strconv.Itoa(file.Channels())
  values["comment"]    = file.Comment()
  values["genre"]      = file.Genre()
  values["length"]     = file.Length().String()
  values["samplerate"] = strconv.Itoa(file.Samplerate())
  values["title"]      = file.Title()
  values["track"]      = strHideZero(file.Track())
  values["year"]       = strHideZero(file.Year())
  return values
}

