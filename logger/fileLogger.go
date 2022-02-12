package logger

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

type interactiveLogger struct {
	stream *os.File
	prefix string
}

func (il *interactiveLogger) PrintFormattedReader(skipBytes int, resultFormat string, reader io.Reader) (map[string]string, error) {
	bf := bufio.NewReader(reader)
	line, err := bf.ReadString('\n')

	reg := regexp.MustCompile("^::output::([a-zA-Z]+)::(.*)::$")
	output := make(map[string]string)

	for err == nil {
		formattedLine := string([]byte(line)[skipBytes:])
		trimmedLine := strings.Trim(formattedLine, "\n \t\r")

		// search for output
		submatchs := reg.FindStringSubmatch(trimmedLine)

		if submatchs != nil {
			varName := submatchs[1]
			value := submatchs[2]

			output[varName] = value
		}

		writeError := il.writeLine(fmt.Sprintf(resultFormat, trimmedLine))
		if writeError != nil {
			return nil, writeError
		}

		line, err = bf.ReadString('\n')
	}
	return output, nil
}

func (il *interactiveLogger) writeLine(line string) error {
	_, err := fmt.Fprintf(il.stream, "%s %s\n", il.prefix, line)
	return err
}

func (il *interactiveLogger) Copy(prefixExtension string) *interactiveLogger {
	return &interactiveLogger{
		stream: il.stream,
		prefix: fmt.Sprintf("%s%s", il.prefix, prefixExtension),
	}
}

func (il *interactiveLogger) Clean() {
	err := il.stream.Close()
	if err != nil {
		log.Printf("[WARNING] fail to close the log file : %v\n", err)
	}
}

func (il *interactiveLogger) Print(message string) error {
	return il.writeLine(message)
}
