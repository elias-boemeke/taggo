# taggo

**taggo** is a command line tool for reading and editing meta data embedded
into audio files. It is written in go and depends on taglib and go-taglib.

**taglib** is a library providing reading and editing of meta data to other
programs. In order to run taggo, taglib has to be installed on your system.

**go-taglib** provides go language bindings from taglib (written in C) to go.
It is also required to run taggo.

Links: [taglib](https://taglib.org/) [go-taglib](https://github.com/wtolson/go-taglib)


## Installation

First set up your go environment if not already done.  
If you don't know how to do that, here is a tutorial:
[How to Write Go Code](https://golang.org/doc/code.html)

- recommended: add your go bin directory to PATH (re-login to apply changes)  
  this lets you run taggo without specifying an absolute path
```
# in file .profile
...
export PATH="${PATH}:${GOPATH}/bin"
```

- Install taglib with your favorite package manager i.e. `pacman -S taglib`
- Get the required source code:
```
go get github.com/wtolson/go-taglib
go get github.com/elias-boemeke/taggo
```

- Build and install the program
```
cd $GOPATH
go install github.com/elias-boemeke/taggo
```

## Usage

To get detailed information about taggo's functions, run
`taggo` without arguments or `taggo --help`


## Functions

taggo supports reading and editing of the tags
`Album, Artist, Comment, Genre, Title, Track, Year` and reading of the
tags `Bitrate, Channels, Length Samplerate`


## Examples

`taggo test.mp3` displays tags of file `test.mp3`

`taggo test.mp3 --clear` clear all tags of file `test.mp3`

`taggo test.mp3 -r "The Artist"` set the Artist tag of `test.mp3` to `The Artist`

`taggo test.mp3 -c "A Comment" -s simple` set the Comment tag of `test.mp3` to
`A Comment` and display a few basic tags (simple is a mode of display, for more
see `taggo --help show`)

`taggo test.mp3 --clear-genre --show-format "year;%y\nalbum;%l"` clear the Genre tag
of file `test.mp3` and display the tags Year and Album in the given format
(for more information on display modes and format see `taggo --help show`)

`taggo -f -dashfile -k 5` set the Track tag of file `-dashfile` to `5`

**Note:**

see `taggo --help` for the manual of the tool


