package reader

import "parser/iterator"

type Reader[T any] interface {
	Run() iterator.Iterator[T]
}


