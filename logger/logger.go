package logger

import "io"

type Context struct {
	RunName string
}

type Logger interface {
	Init(Context) error
	Log(string, io.Reader) error
	Debug(string)
}

var LOG = FileLogger{}
