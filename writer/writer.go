package writer

import (
	"fmt"
	"os"
	"sync"
)

type ConsoleAndFileWriter struct {
	consoleWriter *ConsoleWriter
	fileWriter    *FileWriter
	ResultChannel chan bool
}

func NewConsoleAndFileWriter(filename string, resultChannel chan bool) *ConsoleAndFileWriter {
	return &ConsoleAndFileWriter{
		consoleWriter: &ConsoleWriter{},
		fileWriter:    NewFileWriter(filename),
		ResultChannel: resultChannel,
	}
}

func (w *ConsoleAndFileWriter) Write(value string) {
	w.consoleWriter.Write(value)

	if err := w.fileWriter.Write(value); err != nil {
		fmt.Println("Error writing to file:", err)
		w.ResultChannel <- false
	}
}

type ConsoleWriter struct{}

func (cw *ConsoleWriter) Write(value string) {
	fmt.Println("Writing to console:", value)
}

type FileWriter struct {
	filename string
	mutex    sync.Mutex
}

func NewFileWriter(filename string) *FileWriter {
	return &FileWriter{filename: filename}
}

func (fw *FileWriter) Write(value string) error {
	fw.mutex.Lock()
	defer fw.mutex.Unlock()

	if _, err := os.Stat(fw.filename); os.IsNotExist(err) {
		if _, err := os.Create(fw.filename); err != nil {
			return fmt.Errorf("error creating file: %v", err)
		}
	}

	file, err := os.OpenFile(fw.filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(value + "\n"); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	return nil
}
