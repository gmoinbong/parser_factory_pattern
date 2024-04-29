package reader

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"parser/converter"
	"parser/iterator"
	"parser/validator"
	"sync"
	"time"
)

type FileReader[T any] struct {
	files     []string
	converter converter.Converter[T]
	validator validator.Validator
}

func NewFileReader[T any](files []string, converter converter.Converter[T], validator validator.Validator) *FileReader[T] {
	return &FileReader[T]{files: files, converter: converter, validator: validator}
}

func (fr *FileReader[T]) Run() iterator.Iterator[T] {
	resultChannel := make(chan T)
	var wg sync.WaitGroup

	for _, file := range fr.files {
		wg.Add(1)
		go func(fileName string) {
			defer wg.Done()
			fr.processFile(fileName, resultChannel)
		}(file)
	}

	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	return &FileIterator[T]{resultChannel: resultChannel}
}

func (fr *FileReader[T]) processFile(fileName string, resultChannel chan T) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !fr.validator.Validate(line) {
			fmt.Printf("Invalid data in file %s: %s\n", fileName, line)
			return
		}

		val, err := fr.converter.Convert(line)
		if err != nil {
			fmt.Printf("Error converting data in file %s: %s\n", fileName, err)
			continue
		}

		randomDelay := rand.Intn(5000)
		time.Sleep(time.Duration(randomDelay) * time.Millisecond)

		resultChannel <- val
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
	}
}

type FileIterator[T any] struct {
	resultChannel <-chan T
}

func (fi *FileIterator[T]) HasNext() bool {
	return true
}

func (fi *FileIterator[T]) Next() iterator.Iteration[T] {
	val, ok := <-fi.resultChannel
	if !ok {
		return iterator.NewStepErr[T](fmt.Errorf("no more data"))
	}
	return iterator.NewStepVal[T](val)
}

func (fi *FileIterator[T]) Close() {
}
