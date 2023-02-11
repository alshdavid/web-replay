package extras

func Must[T any](value T, err error) T {
	return value
}
