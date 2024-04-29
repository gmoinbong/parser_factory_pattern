package converter

import "strconv"

type Converter[T any] interface {
	Convert(s string) (T, error)
}

type StringConverter struct{}

func (c StringConverter) Convert(s string) (string, error) {
	return s, nil
}

type IntConverter struct{}

func (c IntConverter) Convert(s string) (int, error) {
	return strconv.Atoi(s)
}
