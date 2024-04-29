package attemptfactory

type AttemptFactory[T any] struct {
	callbacks []func(T) bool
}
