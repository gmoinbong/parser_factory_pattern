package main

import (
	"log"
	"parser/attemptfactory"
	"parser/converter"
	"parser/parserfactory"
	"parser/processor"
	"parser/reader"
	"parser/validator"
	"parser/writer"
	"time"
)

func main() {
	resultChannel := make(chan bool)
	errorChannel := make(chan error)

	jsonValidator := validator.NewJSONValidator()

	fileReader := reader.NewFileReader([]string{"file1.txt", "file2.txt", "file3.txt"}, converter.StringConverter{}, jsonValidator)

	parserJSON := processor.NewParserJSON(resultChannel, errorChannel)
	fileAndConsoleWriter := writer.NewConsoleAndFileWriter("output.txt", resultChannel)
	var attemptFactory attemptfactory.AttemptFactory[string]

	appFactory := parserfactory.NewAppFactory(fileReader, parserJSON, fileAndConsoleWriter, attemptFactory)

	go func() {
		log.Println("App started")
		err := appFactory.Run()
		if err != nil {
			log.Println("App error:", err)
		}
		log.Println("Execution completed")
	}()

	go func() {
		for range errorChannel { 
			log.Println("Received error signal from parser")
			resultChannel <- false 
		}
	}()

	go func() {
		for range resultChannel {
			log.Println("Received result signal from writer")
		}
	}()

	time.Sleep(20 * time.Second)
	log.Println("Time to close")
	log.Println("App closed")
}
