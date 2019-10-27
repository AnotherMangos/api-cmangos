package logger

import (
  "fmt"
  "os"
  "time"
)

const (
  _VERBOSITY_DEBUG = 0
  _VERBOSITY_INFO = 1
  _VERBOSITY_ERROR = 2
)

var Verbosity = _VERBOSITY_ERROR

func Debug(msg string) {
  log(msg, _VERBOSITY_DEBUG)
}

func Error(msg string) {
  log(msg, _VERBOSITY_ERROR)
}

func Info(msg string) {
  log(msg, _VERBOSITY_INFO)
}

func log(msg string, lvl int) {
  var l string
  var o *os.File

  if lvl == _VERBOSITY_DEBUG {
    l = "DEBUG"
    o = os.Stdout
  } else if lvl == _VERBOSITY_ERROR {
    l = "ERROR"
    o = os.Stderr
  } else {
    l = "INFO"
    o = os.Stdout
  }
  t := time.Now().UTC().Format("2006-01-02T15:04:05-0700")
  m := fmt.Sprintf("[%s] %s - %s", t, l, msg)

  if lvl >= Verbosity {
    fmt.Fprintf(o, fmt.Sprintf("%s\n", m))
  }
}
