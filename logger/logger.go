package logger

import (
	"io"
	"os"
)

type Context struct {
	RunName string
}

type Logger interface {
	PrintFormattedReader(skipBytes int, resultFormat string, reader io.Reader) error
	Copy(prefixExtension string) *interactiveLogger
}

func New(basePrefix string, stream *os.File) Logger {
	return &interactiveLogger{
		prefix: basePrefix,
		stream: stream,
	}
}
