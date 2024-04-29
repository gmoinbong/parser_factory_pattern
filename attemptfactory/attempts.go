package attemptfactory

import (
	"reflect"
)

func NewAttemptFactory[T any]() *AttemptFactory[T] {
	return &AttemptFactory[T]{
		callbacks: make([]func(T) bool, 0),
	}
}

func (af *AttemptFactory[T]) Run() bool {
	for _, callback := range af.callbacks {
		zeroValue := reflect.Zero(reflect.TypeOf((*T)(nil)).Elem()).Interface().(T)
		if !callback(zeroValue) {
			return false
		}
	}
	return true
}

func (af *AttemptFactory[T]) Add(callback func(T) bool) {
	af.callbacks = append(af.callbacks, callback)
}

func (af *AttemptFactory[T]) GetValidationFuncs() []func(T) bool {
	return af.callbacks
}
