package logger

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type InteractiveLogger struct {
	stream *os.File
	prefix string
}

func (il *InteractiveLogger) PrintFormattedReader(skipBytes int, resultFormat string, reader io.Reader) error {
	bf := bufio.NewReader(reader)
	line, err := bf.ReadString('\n')

	for err == nil {
		formattedLine := string([]byte(line)[skipBytes:])
		trimmedLine := strings.Trim(formattedLine, "\n \t\r")

		writeError := il.writeLine(fmt.Sprintf(resultFormat, trimmedLine))
		if writeError != nil {
			return writeError
		}

		line, err = bf.ReadString('\n')
	}
	return nil
}

func (il *InteractiveLogger) writeLine(line string) error {
	_, err := fmt.Fprintf(il.stream, "%s %s\n", il.prefix, line)
	return err
}

func (il *InteractiveLogger) Copy(prefixExtension string) *InteractiveLogger {
	return &InteractiveLogger{
		stream: il.stream,
		prefix: fmt.Sprintf("%s%s", il.prefix, prefixExtension),
	}
}

func New(basePrefix string, stream *os.File) *InteractiveLogger {
	return &InteractiveLogger{
		prefix: basePrefix,
		stream: stream,
	}
}
