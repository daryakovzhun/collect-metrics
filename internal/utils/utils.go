package utils

func FromPointer[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}

func ToPointer[T any](p T) *T {
	return &p
}
