package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

type FileLogger struct {
	historyPath string
}

func (fl *FileLogger) Init(ctx Context) error {
	fl.historyPath = path.Join("./history", fmt.Sprintf("run-%s", ctx.RunName))

	err := os.MkdirAll(fl.historyPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("can not create run history directory : %v", err)
	}

	stepsPath := path.Join(fl.historyPath, "steps")

	err = os.Mkdir(stepsPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("can not create steps history directory : %v", err)
	}

	return nil
}

func (fl FileLogger) Log(historyRelativePath string, stream io.Reader) error {
	// Create file if not exists
	filePath := path.Join(fl.historyPath, historyRelativePath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create folder
		err := os.MkdirAll(path.Dir(filePath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("fail to create log directory %v", err)
		}
	}

	// Create file
	file, err := os.Create(filePath)
	defer file.Close()

	if err != nil {
		return fmt.Errorf("fail to open log file %v", err)
	}

	_, err = file.ReadFrom(stream)
	if err != nil {
		return fmt.Errorf("fail to write in log file %v", err)
	}
	return nil
}

func (fl FileLogger) Debug(message string) {
	f, err := os.OpenFile(path.Join(fl.historyPath, "debug.log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("[WARNING] : Fail to open debug file : %v", err)
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("%s\n", message)); err != nil {
		log.Printf("[WARNING] : Fail to write into debug file : %v", err)
	}
}
