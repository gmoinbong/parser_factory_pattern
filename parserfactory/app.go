package parserfactory

import (
	"fmt"
	"parser/attemptfactory"
	"parser/processor"
	"parser/reader"
	"parser/writer"
)

type AppFactory[T any] struct {
	reader         reader.Reader[T]
	processor      processor.Processor[T]
	writer         writer.Writer[T]
	attemptFactory attemptfactory.AttemptFactory[T]
}

func NewAppFactory[T any](reader reader.Reader[T], processor processor.Processor[T], writer writer.Writer[T], attemptFactory attemptfactory.AttemptFactory[T]) *AppFactory[T] {
	return &AppFactory[T]{
		reader:         reader,
		processor:      processor,
		writer:         writer,
		attemptFactory: attemptFactory,
	}
}

func convertToBytes(value interface{}) ([]byte, error) {
	switch v := value.(type) {
	case string:
		return []byte(v), nil
	default:
		return nil, fmt.Errorf("value not equal to string")
	}
}

func (af *AppFactory[T]) Run() error {
	it := af.reader.Run()
	if it == nil {
		return fmt.Errorf("bad getaway iterator")
	}

	validationFuncs := af.attemptFactory.GetValidationFuncs()

	for it.HasNext() {
		step := it.Next()
		if step == nil {
			return fmt.Errorf("error while get next step")
		}

		if err := step.Err(); err != nil {
			it.Close()
			return err
		}

		value := step.Val()

		valueBytes, err := convertToBytes(value)
		if err != nil {
			return err
		}

		var validationFuncsTyped []processor.ValidationFunc[T]
		for _, f := range validationFuncs {
			validationFuncsTyped = append(validationFuncsTyped, f)
		}

		if af.processor.Run(valueBytes, validationFuncsTyped) {
			af.writer.Write(value)
		}
	}

	it.Close()
	return nil
}
