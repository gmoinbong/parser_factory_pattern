package writer

type Writer[T any] interface {
	Write(value T)
}
