package cond

func Ter[T comparable](condition bool, A, B T) T {
	if condition {
		return A
	}

	return B
}
