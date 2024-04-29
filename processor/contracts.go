package processor

type ValidationFunc[T any] func(T) bool

type Processor[T any] interface {
	Run(value []byte, attempts []ValidationFunc[T]) bool
}
