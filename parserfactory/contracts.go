package parserfactory

type ParserFactory[Input, U any] interface {
	Create(Input) U
}
